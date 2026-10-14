package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQuotaFleetTwentyAccountsConcurrentConversations(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			const accounts, conversations = 20, 6
			now := time.Now().Truncate(time.Second)
			started, release := make(chan string, accounts), make(chan struct{})
			var releaseOnce sync.Once
			t.Cleanup(func() { releaseOnce.Do(func() { close(release) }) })
			var mu sync.Mutex
			calls := map[string]int{}
			resets := map[string]time.Duration{}
			m, s, auths := quotaTestManagerAccounts(t, provider, accounts, func(ctx context.Context, a *Auth) ([]byte, error) {
				mu.Lock()
				calls[a.ID]++
				mu.Unlock()
				started <- a.ID
				select {
				case <-release:
					return quotaTestBody(provider, now, resets[a.ID], 20), nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			})
			for i, a := range auths {
				resets[a.ID] = time.Duration(48+i) * time.Hour
			}
			preferred := auths[7].ID
			resets[preferred] = time.Hour
			type result struct {
				a   *Auth
				err error
			}
			results := make(chan result, conversations)
			for i := range conversations {
				go func() {
					a, err := m.SelectAuth(t.Context(), provider, "quota-test-model", quotaTestOptions(fmt.Sprintf("fleet-conversation-%d", i)))
					results <- result{a, err}
				}()
			}
			// Every required lookup starts before any is allowed to complete.
			for range accounts {
				select {
				case <-started:
				case <-t.Context().Done():
					t.Fatal("parallel lookups did not start")
				}
			}
			releaseOnce.Do(func() { close(release) })
			for range conversations {
				r := <-results
				if r.err != nil || r.a == nil || r.a.ID != preferred {
					t.Fatalf("new conversation did not use nearest reset: auth=%v err=%v", r.a, r.err)
				}
			}
			mu.Lock()
			for _, a := range auths {
				if calls[a.ID] != 1 {
					t.Errorf("account %s had %d lookups, want one shared lookup", a.ID, calls[a.ID])
				}
			}
			mu.Unlock()
			// Reverse reset preference. Existing conversations must retain the
			// account; a new one sees the updated passive reading without a GET.
			s.quota.now = func() time.Time { return now.Add(time.Second) }
			setQuotaTestPassive(m, auths[0], quotaTestSignals(provider, now.Add(time.Second), 30*time.Minute, 20))
			for i := range conversations {
				a, err := m.SelectAuth(t.Context(), provider, "quota-test-model", quotaTestOptions(fmt.Sprintf("fleet-conversation-%d", i)))
				if err != nil || a.ID != preferred {
					t.Fatal("quota update moved a healthy conversation")
				}
			}
			a, err := m.SelectAuth(t.Context(), provider, "quota-test-model", quotaTestOptions("fresh-seventh-conversation"))
			if err != nil || a.ID != auths[0].ID {
				t.Fatalf("new conversation missed updated passive reading: auth=%v err=%v", a, err)
			}
			mu.Lock()
			defer mu.Unlock()
			for _, a := range auths {
				if calls[a.ID] != 1 {
					t.Fatal("follow-up or fresh-cache selection fetched again")
				}
			}
		})
	}
}

func TestQuotaFleetOneFailedAccountFallsBackAndRecovers(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			now := time.Now().Truncate(time.Second)
			var clock atomic.Int64
			clock.Store(now.UnixNano())
			var broken string
			var calls atomic.Int32
			m, s, auths := quotaTestManagerAccounts(t, provider, 20, func(context.Context, *Auth) ([]byte, error) {
				calls.Add(1)
				return nil, errors.New("usage endpoint unavailable")
			})
			s.quota.now = func() time.Time { return time.Unix(0, clock.Load()) }
			broken = auths[19].ID
			for i, a := range auths[:19] {
				setQuotaTestPassive(m, a, quotaTestSignals(provider, now, time.Duration(i+1)*time.Hour, 20))
			}
			var selected []string
			for i := range 20 {
				a, err := m.SelectAuth(t.Context(), provider, "quota-test-model", quotaTestOptions(fmt.Sprintf("fallback-%d", i)))
				if err != nil || a == nil {
					t.Fatalf("one broken quota endpoint stopped work: %v", err)
				}
				selected = append(selected, a.ID)
			}
			seen := map[string]bool{}
			for _, id := range selected {
				seen[id] = true
			}
			if len(seen) != 20 || !seen[broken] || calls.Load() != 1 {
				t.Fatalf("baseline pool/backoff changed: accounts=%d lookups=%d", len(seen), calls.Load())
			}
			// Fresh passive data can restore quota ranking even during the
			// failed usage endpoint's backoff, without an additional network call.
			clock.Store(now.Add(time.Second).UnixNano())
			setQuotaTestPassive(m, auths[19], quotaTestSignals(provider, now.Add(time.Second), 30*time.Minute, 20))
			a, err := m.SelectAuth(t.Context(), provider, "quota-test-model", quotaTestOptions("recovered"))
			if err != nil || a.ID != broken || calls.Load() != 1 {
				t.Fatal("passive recovery did not restore nearest-reset ranking")
			}
		})
	}
}

func TestQuotaFleetResetDuringLookupDoesNotFabricateCapacity(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	var clock atomic.Int64
	clock.Store(now.UnixNano())
	started, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	t.Cleanup(func() { releaseOnce.Do(func() { close(release) }) })
	m, s, auths := quotaTestManagerAccounts(t, "claude", 20, func(ctx context.Context, _ *Auth) ([]byte, error) {
		close(started)
		select {
		case <-release:
			return quotaTestBody("claude", now, time.Minute, 20), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	})
	s.quota.now = func() time.Time { return time.Unix(0, clock.Load()) }
	for _, a := range auths[1:] {
		setQuotaTestPassive(m, a, quotaTestSignals("claude", now, 24*time.Hour, 20))
	}
	done := make(chan error, 1)
	go func() {
		_, err := m.SelectAuth(t.Context(), "claude", "quota-test-model", quotaTestOptions("reset-during-lookup"))
		done <- err
	}()
	<-started
	clock.Store(now.Add(time.Minute).UnixNano())
	releaseOnce.Do(func() { close(release) })
	if err := <-done; err != nil {
		t.Fatalf("expired measurement stopped baseline selection: %v", err)
	}
	if _, ok := s.quota.snapshots(auths)[auths[0].ID]; ok {
		t.Fatal("crossing reset fabricated or accepted expired quota")
	}
}
