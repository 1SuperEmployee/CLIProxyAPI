package auth

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestQuotaCacheFreshPassiveAndResetBoundary(t *testing.T) {
	now := time.Unix(1800000000, 0)
	cfg := DefaultQuotaSelectionConfig()
	c := newSelectionQuotaCache(&cfg)
	c.now = func() time.Time { return now }
	a := quotaTestAuth("codex", "a")
	a.Quota = quotaTestSignals("codex", now, time.Minute, 20)
	var calls int
	fetch := func(context.Context, *Auth) ([]byte, error) {
		calls++
		return quotaTestBody("codex", now, time.Hour, 10), nil
	}
	if c.refresh(t.Context(), []*Auth{a}, fetch) || calls != 0 {
		t.Fatal("fresh passive observation triggered network")
	}
	now = now.Add(time.Minute) // exact reset invalidates despite young observation
	if !c.refresh(t.Context(), []*Auth{a}, fetch) || calls != 1 {
		t.Fatal("reset boundary did not refresh")
	}
	observed := c.snapshots([]*Auth{a})[a.ID].weekly.observedAt
	now = now.Add(time.Minute)
	if !c.snapshots([]*Auth{a})[a.ID].weekly.observedAt.Equal(observed) {
		t.Fatal("read extended freshness")
	}
}

func TestQuotaCacheCoalescesAndWaiterCancellationDoesNotCancelSharedLookup(t *testing.T) {
	now := time.Unix(1800000000, 0)
	cfg := DefaultQuotaSelectionConfig()
	c := newSelectionQuotaCache(&cfg)
	c.now = func() time.Time { return now }
	a := quotaTestAuth("codex", "a")
	started, release := make(chan struct{}), make(chan struct{})
	var calls atomic.Int32
	fetch := func(ctx context.Context, _ *Auth) ([]byte, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		select {
		case <-release:
			return quotaTestBody("codex", now, time.Hour, 20), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	ctx, cancel := context.WithCancel(t.Context())
	first := make(chan struct{})
	go func() { c.refresh(ctx, []*Auth{a}, fetch); close(first) }()
	<-started
	cancel()
	<-first // departing waiter does not destroy the pending shared lookup
	second := make(chan struct{})
	go func() { c.refresh(t.Context(), []*Auth{a}, fetch); close(second) }()
	close(release)
	<-second
	if calls.Load() != 1 || len(c.snapshots([]*Auth{a})) != 1 {
		t.Fatalf("calls=%d; shared result missing or duplicated", calls.Load())
	}
}

func TestQuotaCacheFailureBackoffAndCredentialReplacement(t *testing.T) {
	now := time.Unix(1800000000, 0)
	cfg := DefaultQuotaSelectionConfig()
	c := newSelectionQuotaCache(&cfg)
	c.now = func() time.Time { return now }
	a := quotaTestAuth("codex", "a")
	var calls int
	fetch := func(context.Context, *Auth) ([]byte, error) { calls++; return nil, errors.New("unavailable") }
	c.refresh(t.Context(), []*Auth{a}, fetch)
	c.refresh(t.Context(), []*Auth{a}, fetch)
	if calls != 1 {
		t.Fatal("failure did not back off")
	}
	now = now.Add(cfg.FailureBackoff)
	c.refresh(t.Context(), []*Auth{a}, fetch)
	if calls != 2 {
		t.Fatal("failure backoff did not expire")
	}
	b := a.Clone()
	b.CredentialVersion++
	c.refresh(t.Context(), []*Auth{b}, fetch)
	if calls != 3 {
		t.Fatal("replacement reused prior credential's failure backoff")
	}
	if len(c.snapshots([]*Auth{a})) != 0 {
		t.Fatal("stale caller rolled back credential identity")
	}
}

func TestQuotaCacheOlderLookupCannotOverwritePassiveOrReplacement(t *testing.T) {
	for _, replacement := range []bool{false, true} {
		t.Run(map[bool]string{false: "newer passive", true: "replacement"}[replacement], func(t *testing.T) {
			now := time.Unix(1800000000, 0)
			cfg := DefaultQuotaSelectionConfig()
			c := newSelectionQuotaCache(&cfg)
			var clock atomic.Int64
			clock.Store(now.UnixNano())
			c.now = func() time.Time { return time.Unix(0, clock.Load()) }
			a := quotaTestAuth("codex", "a")
			started, release, complete := make(chan struct{}), make(chan struct{}), make(chan struct{})
			go func() {
				c.refresh(t.Context(), []*Auth{a}, func(context.Context, *Auth) ([]byte, error) {
					close(started)
					<-release
					return quotaTestBody("codex", now, time.Hour, 10), nil
				})
				close(complete)
			}()
			<-started
			clock.Store(now.Add(time.Second).UnixNano())
			b := a.Clone()
			if replacement {
				b.RegistrationEpoch++
			}
			b.Quota = quotaTestSignals("codex", now.Add(time.Second), 2*time.Hour, 80)
			c.snapshots([]*Auth{b})
			close(release)
			<-complete
			snapshot := c.snapshots([]*Auth{b})[b.ID]
			if snapshot.weekly.remainingFraction > 0.21 || !snapshot.weekly.observedAt.Equal(b.Quota.ObservedAt) {
				t.Fatal("older lookup overwrote authoritative newer observation")
			}
		})
	}
}

func TestQuotaCacheBoundedParallelDeadline(t *testing.T) {
	cfg := DefaultQuotaSelectionConfig()
	// A real deadline test uses context completion, not sleep-based TTL tests.
	cfg.LookupTimeout = 50 * time.Millisecond
	c := newSelectionQuotaCache(&cfg)
	auths := []*Auth{quotaTestAuth("codex", "a"), quotaTestAuth("codex", "b")}
	started, completed := make(chan string, 2), make(chan struct{}, 2)
	fetch := func(ctx context.Context, a *Auth) ([]byte, error) {
		started <- a.ID
		<-ctx.Done()
		completed <- struct{}{}
		return nil, ctx.Err()
	}
	done := make(chan struct{})
	go func() { c.refresh(t.Context(), auths, fetch); close(done) }()
	// Both start before either finishes: one total deadline, not sequential 2s.
	<-started
	<-started
	<-done
	<-completed
	<-completed
	if len(c.snapshots(auths)) != 0 {
		t.Fatal("deadline failure produced a fabricated measurement")
	}
}

func TestQuotaCacheReplacementCannotReuseCarriedPassiveReading(t *testing.T) {
	now := time.Unix(1800000000, 0)
	cfg := DefaultQuotaSelectionConfig()
	c := newSelectionQuotaCache(&cfg)
	c.now = func() time.Time { return now }
	a := quotaTestAuth("codex", "a")
	a.Quota = quotaTestSignals("codex", now, time.Hour, 20)
	if len(c.snapshots([]*Auth{a})) != 1 {
		t.Fatal("initial passive reading missing")
	}
	b := a.Clone()
	b.CredentialVersion++
	if len(c.snapshots([]*Auth{b})) != 0 {
		t.Fatal("replacement reused carried passive reading")
	}
	var calls int
	c.refresh(t.Context(), []*Auth{b}, func(context.Context, *Auth) ([]byte, error) {
		calls++
		return quotaTestBody("codex", now, 2*time.Hour, 50), nil
	})
	if calls != 1 || c.snapshots([]*Auth{b})[b.ID].weekly.remainingFraction != 0.5 {
		t.Fatal("replacement did not acquire its own reading")
	}
}

func TestQuotaCacheInvalidFutureObservationDoesNotPoisonRefresh(t *testing.T) {
	now := time.Unix(1800000000, 0)
	cfg := DefaultQuotaSelectionConfig()
	c := newSelectionQuotaCache(&cfg)
	c.now = func() time.Time { return now }
	a := quotaTestAuth("codex", "a")
	a.Quota = quotaTestSignals("codex", now.Add(time.Hour), time.Hour, 20)
	c.refresh(t.Context(), []*Auth{a}, func(context.Context, *Auth) ([]byte, error) { return quotaTestBody("codex", now, time.Hour, 30), nil })
	s := c.snapshots([]*Auth{a})[a.ID]
	if s.weekly.remainingFraction != 0.7 || !s.weekly.observedAt.Equal(now) {
		t.Fatal("invalid future passive observation blocked successful refresh")
	}
}
