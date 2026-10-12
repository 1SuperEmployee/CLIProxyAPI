package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/registry"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
)

type quotaManagerTestExecutor struct {
	schedulerProviderTestExecutor
	fetch func(context.Context, *Auth) ([]byte, error)
}

func (e quotaManagerTestExecutor) FetchSelectionQuota(ctx context.Context, a *Auth) ([]byte, error) {
	return e.fetch(ctx, a)
}

func quotaTestOptions(id string) cliproxyexecutor.Options {
	return cliproxyexecutor.Options{Headers: http.Header{"X-Session-Id": {id}}}
}

func quotaTestManager(t *testing.T, provider string, fetch func(context.Context, *Auth) ([]byte, error)) (*Manager, *SessionAffinitySelector, []*Auth) {
	t.Helper()
	cfg := DefaultQuotaSelectionConfig()
	s := NewSessionAffinitySelectorWithConfig(SessionAffinityConfig{QuotaSelection: &cfg})
	t.Cleanup(s.Stop)
	m := NewManager(nil, s, nil)
	m.RegisterExecutor(quotaManagerTestExecutor{schedulerProviderTestExecutor{provider}, fetch})
	var auths []*Auth
	for _, label := range []string{"a", "b"} {
		a, err := m.Register(t.Context(), quotaTestAuth(provider, t.Name()+"-"+label))
		if err != nil {
			t.Fatal(err)
		}
		auths = append(auths, a)
		registry.GetGlobalRegistry().RegisterClient(a.ID, provider, []*registry.ModelInfo{{ID: "quota-test-model"}, {ID: "other-quota-test-model"}})
		t.Cleanup(func() { registry.GetGlobalRegistry().UnregisterClient(a.ID) })
	}
	return m, s, auths
}

func setQuotaTestPassive(m *Manager, a *Auth, q QuotaState) {
	m.mu.Lock()
	m.auths[a.ID].Quota = q
	m.mu.Unlock()
}

func TestQuotaManagerFreshPassiveAcrossModelsAndHealthyFollowup(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			var calls atomic.Int32
			m, _, auths := quotaTestManager(t, provider, func(context.Context, *Auth) ([]byte, error) {
				calls.Add(1)
				return nil, errors.New("must not fetch")
			})
			now := time.Now().Truncate(time.Second)
			setQuotaTestPassive(m, auths[0], quotaTestSignals(provider, now, 6*time.Hour, 20))
			setQuotaTestPassive(m, auths[1], quotaTestSignals(provider, now, 3*time.Hour, 50))
			for _, model := range []string{"quota-test-model", "other-quota-test-model"} {
				a, err := m.SelectAuth(t.Context(), provider, model, quotaTestOptions(model))
				if err != nil || a.ID != auths[1].ID {
					t.Fatalf("model %s did not use shared nearest reset: auth=%v err=%v", model, a, err)
				}
			}
			// Stale quotas on BOTH accounts cannot disturb an established thread.
			setQuotaTestPassive(m, auths[0], quotaTestSignals(provider, now.Add(-time.Hour), 6*time.Hour, 20))
			setQuotaTestPassive(m, auths[1], quotaTestSignals(provider, now.Add(-time.Hour), 3*time.Hour, 99))
			a, err := m.SelectAuth(t.Context(), provider, "quota-test-model", quotaTestOptions("quota-test-model"))
			if err != nil || a.ID != auths[1].ID || calls.Load() != 0 {
				t.Fatalf("healthy follow-up fetched or switched: auth=%v err=%v calls=%d", a, err, calls.Load())
			}
		})
	}
}

func TestQuotaManagerConcurrentColdStartsOneLookupAndOneBinding(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	started, release := make(chan string, 2), make(chan struct{})
	var mu sync.Mutex
	calls := map[string]int{}
	var bID string
	m, _, auths := quotaTestManager(t, "codex", func(ctx context.Context, a *Auth) ([]byte, error) {
		mu.Lock()
		calls[a.ID]++
		mu.Unlock()
		started <- a.ID
		select {
		case <-release:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		reset := 6 * time.Hour
		if a.ID == bID {
			reset = 3 * time.Hour
		}
		return quotaTestBody("codex", now, reset, 20), nil
	})
	bID = auths[1].ID
	type result struct {
		a   *Auth
		err error
	}
	results := make(chan result, 8)
	for range 8 {
		go func() {
			a, err := m.SelectAuth(t.Context(), "codex", "quota-test-model", quotaTestOptions("one-concurrent-conversation"))
			results <- result{a, err}
		}()
	}
	<-started
	<-started
	close(release)
	for range 8 {
		r := <-results
		if r.err != nil || r.a.ID != bID {
			t.Fatalf("cold conversation split or failed: %v %v", r.a, r.err)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	for _, a := range auths {
		if calls[a.ID] != 1 {
			t.Fatalf("%s: quota calls=%d, want 1", a.ID, calls[a.ID])
		}
	}
}

func TestQuotaManagerRechecksBindingAndEligibilityAfterLookup(t *testing.T) {
	for _, change := range []string{"binding", "disable", "replace"} {
		t.Run(change, func(t *testing.T) {
			now := time.Now().Truncate(time.Second)
			started, release := make(chan struct{}, 2), make(chan struct{})
			m, s, auths := quotaTestManager(t, "codex", func(ctx context.Context, _ *Auth) ([]byte, error) {
				started <- struct{}{}
				select {
				case <-release:
					return quotaTestBody("codex", now, 3*time.Hour, 20), nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			})
			opts := quotaTestOptions("waiting-conversation")
			result := make(chan *Auth, 1)
			errs := make(chan error, 1)
			go func() {
				a, err := m.SelectAuth(t.Context(), "codex", "quota-test-model", opts)
				result <- a
				errs <- err
			}()
			<-started
			<-started
			switch change {
			case "binding":
				primary, _ := extractExplicitSessionIDs(opts.Headers, nil, nil)
				s.cache.Set("codex::"+primary+"::quota-test-model", auths[1].ID)
			case "disable":
				a, _ := m.GetByID(auths[0].ID)
				a.Disabled = true
				if _, err := m.Update(t.Context(), a); err != nil {
					t.Fatal(err)
				}
			case "replace":
				a := auths[0].Clone()
				if _, err := m.Register(t.Context(), a); err != nil {
					t.Fatal(err)
				}
			}
			close(release)
			a, err := <-result, <-errs
			if err != nil || a == nil {
				t.Fatalf("selection failed: %v", err)
			}
			if change != "replace" && a.ID != auths[1].ID {
				t.Fatalf("did not recheck %s: got %s", change, a.ID)
			}
			if change == "replace" {
				current, _ := m.GetByID(a.ID)
				if a.RegistrationEpoch != current.RegistrationEpoch || a.CredentialVersion != current.CredentialVersion {
					t.Fatal("returned replaced credentials")
				}
				if len(s.quota.snapshots([]*Auth{current})) != 0 && a.ID == auths[0].ID {
					t.Fatal("replacement reused old lookup measurement")
				}
			}
		})
	}
}

func TestQuotaManagerLookupFailureKeepsBaselineAndSingleAccountDoesNotFetch(t *testing.T) {
	var calls atomic.Int32
	m, _, auths := quotaTestManager(t, "codex", func(context.Context, *Auth) ([]byte, error) {
		calls.Add(1)
		return nil, errors.New("usage unavailable")
	})
	first, err := m.SelectAuth(t.Context(), "codex", "quota-test-model", quotaTestOptions("first"))
	if err != nil || first == nil || calls.Load() != 2 {
		t.Fatalf("fallback failed: %v calls=%d", err, calls.Load())
	}
	second, err := m.SelectAuth(t.Context(), "codex", "quota-test-model", quotaTestOptions("second"))
	if err != nil || second.ID == first.ID || calls.Load() != 2 {
		t.Fatal("lookup failure broke rotation or ignored backoff")
	}
	a, _ := m.GetByID(auths[1].ID)
	a.Disabled = true
	if _, err := m.Update(t.Context(), a); err != nil {
		t.Fatal(err)
	}
	single, err := m.SelectAuth(t.Context(), "codex", "quota-test-model", quotaTestOptions("third"))
	if err != nil || single.ID != auths[0].ID || calls.Load() != 2 {
		t.Fatal("single-account selection fetched or failed")
	}
}

func TestQuotaManagerPreservesTransportPreferenceBeforeQuota(t *testing.T) {
	for _, fallback := range []Selector{&RoundRobinSelector{}, &WeightedRoundRobinSelector{}, &FillFirstSelector{}} {
		t.Run(fmt.Sprintf("%T", fallback), func(t *testing.T) {
			var calls atomic.Int32
			m, s, auths := quotaTestManager(t, "codex", func(context.Context, *Auth) ([]byte, error) { calls.Add(1); return nil, errors.New("must not fetch") })
			s.fallback = fallback
			now := time.Now().Truncate(time.Second)
			setQuotaTestPassive(m, auths[0], quotaTestSignals("codex", now, 3*time.Hour, 20))
			setQuotaTestPassive(m, auths[1], quotaTestSignals("codex", now, 6*time.Hour, 20))
			m.mu.Lock()
			m.auths[auths[1].ID].Attributes["websockets"] = "true"
			m.mu.Unlock()
			ctx := cliproxyexecutor.WithDownstreamWebsocket(t.Context())
			a, err := m.SelectAuth(ctx, "codex", "quota-test-model", quotaTestOptions("ws-cold"))
			if err != nil || a.ID != auths[1].ID || calls.Load() != 0 {
				t.Fatal("quota overrode transport preference or fetched a single eligible account")
			}
			// A healthy HTTP binding still outranks the cold WS preference.
			a, err = m.SelectAuth(t.Context(), "codex", "quota-test-model", quotaTestOptions("existing-http"))
			if err != nil || a.ID != auths[0].ID {
				t.Fatal("could not establish HTTP binding")
			}
			a, err = m.SelectAuth(ctx, "codex", "quota-test-model", quotaTestOptions("existing-http"))
			if err != nil || a.ID != auths[0].ID || calls.Load() != 0 {
				t.Fatal("transport preference disturbed healthy binding")
			}
		})
	}
}
