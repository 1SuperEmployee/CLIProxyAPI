package auth

import (
	"context"
	"sync"
	"time"
)

// QuotaSelectionConfig is opt-in; nil leaves routing unchanged. These defaults
// are candidate settings, not a promise of provider freshness or equal capacity.
type QuotaSelectionConfig struct {
	MaxAge               time.Duration
	MinRemainingFraction float64
	LookupTimeout        time.Duration
	FailureBackoff       time.Duration
}

func DefaultQuotaSelectionConfig() QuotaSelectionConfig {
	return QuotaSelectionConfig{15 * time.Minute, 0.05, 2 * time.Second, 30 * time.Second}
}

type quotaUsageFetcher interface {
	// FetchSelectionQuota reads only shared subscription usage. The caller owns
	// its bounded context. Implementations must not log bodies or credentials.
	FetchSelectionQuota(context.Context, *Auth) ([]byte, error)
}

type quotaCacheEntry struct {
	epoch, version uint64
	snapshot       selectionQuotaSnapshot
	hasSnapshot    bool
	retryAt        time.Time
	pending        chan struct{}
	passiveThrough time.Time
}

// One bounded in-memory entry per credential, with no poller or durable store.
// Pending lookups use their own bounded context: one departing waiter does not
// cancel a lookup another conversation is already waiting for.
type selectionQuotaCache struct {
	mu      sync.Mutex
	entries map[string]*quotaCacheEntry
	config  QuotaSelectionConfig
	now     func() time.Time
}

const maxSelectionQuotaAccounts = 256

func newSelectionQuotaCache(cfg *QuotaSelectionConfig) *selectionQuotaCache {
	if cfg == nil || cfg.MaxAge <= 0 || cfg.LookupTimeout <= 0 || cfg.FailureBackoff <= 0 || !validFraction(cfg.MinRemainingFraction) || cfg.MinRemainingFraction == 0 {
		return nil
	}
	return &selectionQuotaCache{entries: make(map[string]*quotaCacheEntry), config: *cfg, now: time.Now}
}

func (c *selectionQuotaCache) policy() quotaSelectionPolicy {
	return quotaSelectionPolicy{c.config.MaxAge, c.config.MinRemainingFraction}
}

// entryLocked invalidates measurements and pending results after a credential
// replacement/re-registration. An older caller cannot replace the newer entry.
func (c *selectionQuotaCache) entryLocked(a *Auth) *quotaCacheEntry {
	e := c.entries[a.ID]
	if e != nil && (a.RegistrationEpoch < e.epoch || (a.RegistrationEpoch == e.epoch && a.CredentialVersion < e.version)) {
		return nil
	}
	if e == nil || e.epoch != a.RegistrationEpoch || e.version != a.CredentialVersion {
		if e == nil && len(c.entries) >= maxSelectionQuotaAccounts {
			// Capacity never changes eligibility: unknown data preserves baseline.
			return nil
		}
		var passiveThrough time.Time
		if e != nil {
			passiveThrough = e.passiveThrough
			if e.hasSnapshot && e.snapshot.weekly.observedAt.After(passiveThrough) {
				passiveThrough = e.snapshot.weekly.observedAt
			}
		}
		e = &quotaCacheEntry{epoch: a.RegistrationEpoch, version: a.CredentialVersion, passiveThrough: passiveThrough}
		c.entries[a.ID] = e
	}
	if a.Quota.ObservedAt.After(e.passiveThrough) && !a.Quota.ObservedAt.After(c.now()) {
		if snapshot, ok := selectionQuotaFromSignals(selectionQuotaProvider(a), a.Quota); ok && (!e.hasSnapshot || snapshot.weekly.observedAt.After(e.snapshot.weekly.observedAt)) {
			if _, _, reason := assessSelectionQuota(snapshot, c.now(), c.policy()); reason == "" {
				e.snapshot, e.hasSnapshot = snapshot, true
			}
		}
		e.passiveThrough = a.Quota.ObservedAt
	}
	return e
}

func (c *selectionQuotaCache) snapshots(auths []*Auth) map[string]selectionQuotaSnapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]selectionQuotaSnapshot, len(auths))
	for _, a := range auths {
		if selectionQuotaProvider(a) == "" {
			continue
		}
		if e := c.entryLocked(a); e != nil && e.hasSnapshot {
			out[a.ID] = e.snapshot
		}
	}
	return out
}

func (c *selectionQuotaCache) refresh(ctx context.Context, auths []*Auth, fetch func(context.Context, *Auth) ([]byte, error)) bool {
	if ctx.Err() != nil {
		return false
	}
	// All missing accounts start before waiting for any one of them. The total
	// wait is bounded, including joining another conversation's pending lookup.
	waitCtx, cancel := context.WithTimeout(ctx, c.config.LookupTimeout)
	defer cancel()
	var pending []<-chan struct{}
	c.mu.Lock()
	now := c.now()
	for _, a := range auths {
		if selectionQuotaProvider(a) == "" {
			continue
		}
		e := c.entryLocked(a)
		if e == nil {
			continue
		}
		if e.hasSnapshot {
			if _, _, reason := assessSelectionQuota(e.snapshot, now, c.policy()); reason == "" {
				continue
			}
		}
		if e.pending != nil {
			pending = append(pending, e.pending)
			continue
		}
		if now.Before(e.retryAt) {
			continue
		}
		done := make(chan struct{})
		e.pending = done
		pending = append(pending, done)
		copy := a.Clone()
		go c.lookup(context.WithoutCancel(ctx), copy, e, done, now, fetch)
	}
	c.mu.Unlock()
	for _, done := range pending {
		select {
		case <-done:
		case <-waitCtx.Done():
			return len(pending) > 0
		}
	}
	return len(pending) > 0
}

func (c *selectionQuotaCache) lookup(ctx context.Context, a *Auth, entry *quotaCacheEntry, done chan struct{}, started time.Time, fetch func(context.Context, *Auth) ([]byte, error)) {
	ctx, cancel := context.WithTimeout(ctx, c.config.LookupTimeout)
	defer cancel()
	body, err := fetch(ctx, a)
	// Anchor relative reset times and freshness to the beginning of the read,
	// conservatively including network time. Never extend them at cache reads.
	snapshot, valid := selectionQuotaFromBody(selectionQuotaProvider(a), body, started)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries[a.ID] == entry {
		if err == nil && ctx.Err() == nil && valid {
			_, _, reason := assessSelectionQuota(snapshot, c.now(), c.policy())
			if reason == "" {
				if !entry.hasSnapshot || started.After(entry.snapshot.weekly.observedAt) {
					entry.snapshot, entry.hasSnapshot = snapshot, true
				}
				entry.retryAt = time.Time{}
			} else {
				entry.retryAt = c.now().Add(c.config.FailureBackoff)
			}
		} else {
			entry.retryAt = c.now().Add(c.config.FailureBackoff)
		}
		entry.pending = nil
	}
	close(done)
}
