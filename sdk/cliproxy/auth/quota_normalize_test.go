package auth

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func quotaTestBody(provider string, now time.Time, resetAfter time.Duration, used float64) []byte {
	reset := now.Add(resetAfter)
	if provider == "claude" {
		return fmt.Appendf(nil, `{"five_hour":{"utilization":2,"resets_at":%q},"seven_day":{"utilization":%g,"resets_at":%q},"iguana_necktie":{"utilization":102},"extra_usage":{"is_enabled":false}}`, now.Add(time.Hour).Format(time.RFC3339), used, reset.Format(time.RFC3339))
	}
	return fmt.Appendf(nil, `{"rate_limit":{"allowed":true,"limit_reached":false,"primary_window":{"used_percent":%g,"limit_window_seconds":604800,"reset_at":%d},"secondary_window":null},"additional_rate_limits":[{"limit_name":"model-cap","rate_limit":{"limit_reached":true}}]}`, used, reset.Unix())
}

func quotaTestAuth(provider, id string) *Auth {
	return &Auth{ID: id, Provider: provider, RegistrationEpoch: 1, CredentialVersion: 1, Attributes: map[string]string{AttributeAuthKind: AuthKindOAuth}}
}

func quotaTestSignals(provider string, now time.Time, resetAfter time.Duration, used float64) QuotaState {
	q := QuotaState{ObservedAt: now, Signals: map[string]string{}}
	if provider == "claude" {
		for name, duration := range map[string]time.Duration{"5h": time.Hour, "7d": resetAfter} {
			prefix := "anthropic-ratelimit-unified-" + name + "-"
			q.Signals[prefix+"utilization"] = fmt.Sprint(used / 100)
			q.Signals[prefix+"reset"] = fmt.Sprint(now.Add(duration).Unix())
			q.Signals[prefix+"status"] = "allowed"
		}
	} else {
		q.Signals["x-codex-primary-used-percent"] = fmt.Sprint(used)
		q.Signals["x-codex-primary-window-minutes"] = "10080"
		q.Signals["x-codex-primary-reset-after-seconds"] = fmt.Sprint(int64(resetAfter / time.Second))
	}
	return q
}

func TestQuotaNormalizeSharedBodiesAndSignals(t *testing.T) {
	now := time.Unix(1800000000, 0)
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			for _, source := range []string{"body", "signals"} {
				var snapshot selectionQuotaSnapshot
				var ok bool
				if source == "body" {
					snapshot, ok = selectionQuotaFromBody(provider, quotaTestBody(provider, now, 3*time.Hour, 20), now)
				} else {
					snapshot, ok = selectionQuotaFromSignals(provider, quotaTestSignals(provider, now, 3*time.Hour, 20))
				}
				if !ok || snapshot.weekly.remainingFraction != 0.8 || !snapshot.weekly.resetAt.Equal(now.Add(3*time.Hour)) {
					t.Fatalf("%s: shared units/reset incorrect: %+v, ok=%v", source, snapshot, ok)
				}
				if (snapshot.shortWindow == nil) != (provider == "codex") {
					t.Fatalf("%s: short-window applicability incorrect", source)
				}
			}
		})
	}
}

func TestQuotaNormalizeRejectsIncompleteOrMalformed(t *testing.T) {
	now := time.Unix(1800000000, 0)
	for _, body := range []string{
		`{}`, `{"five_hour":null}`, `{"five_hour":null,"seven_day":{"utilization":-2,"resets_at":"2027-01-16T09:00:00Z"}}`,
		`{"five_hour":null,"seven_day":{"resets_at":"2027-01-16T09:00:00Z"}}`,
		`{"five_hour":{"utilization":1},"seven_day":null}`,
		`{"rate_limit":{"primary_window":{"used_percent":1,"limit_window_seconds":604800,"reset_at":1800000010}}}`,
		`{"rate_limit":{"primary_window":{"used_percent":101,"limit_window_seconds":604800,"reset_at":1800000010},"secondary_window":null}}`,
		`{"rate_limit":{"primary_window":{"used_percent":1,"limit_window_seconds":9223372036854775807,"reset_at":1800000010},"secondary_window":null}}`,
	} {
		for _, provider := range []string{"claude", "codex"} {
			if _, ok := selectionQuotaFromBody(provider, []byte(body), now); ok {
				t.Fatalf("accepted incomplete/malformed %s quota", provider)
			}
		}
	}
	for _, mutation := range []func(map[string]string){
		func(h map[string]string) { h["x-codex-primary-used-percent"] = "NaN" },
		func(h map[string]string) { h["x-codex-primary-used-percent"] = "Inf" },
		func(h map[string]string) { h["x-codex-primary-window-minutes"] = "604800" },
		func(h map[string]string) { h["x-codex-secondary-used-percent"] = "5" },
		func(h map[string]string) { h["x-codex-primary-reset-at"] = "broken" },
	} {
		q := quotaTestSignals("codex", now, time.Hour, 20)
		mutation(q.Signals)
		if _, ok := selectionQuotaFromSignals("codex", q); ok {
			t.Fatal("accepted malformed signal")
		}
	}
}

func TestQuotaNormalizeDurationNotPrimaryNameAndResetAnchor(t *testing.T) {
	now := time.Unix(1800000000, 0)
	q := quotaTestSignals("codex", now, time.Hour, 20)
	q.Signals["x-codex-primary-window-minutes"] = "300"
	q.Signals["x-codex-secondary-used-percent"] = "40"
	q.Signals["x-codex-secondary-window-minutes"] = "10080"
	q.Signals["x-codex-secondary-reset-at"] = fmt.Sprint(now.Add(3 * time.Hour).Unix())
	q.Signals["x-codex-secondary-reset-after-seconds"] = "1" // absolute wins
	snapshot, ok := selectionQuotaFromSignals("codex", q)
	if !ok || snapshot.weekly.remainingFraction != 0.6 || snapshot.shortWindow.duration != 5*time.Hour || !snapshot.shortWindow.resetAt.Equal(now.Add(time.Hour)) {
		t.Fatalf("duration/absolute reset interpretation incorrect: %+v", snapshot)
	}
	// Re-reading at a later time must not re-anchor relative reset or freshness.
	again, _ := selectionQuotaFromSignals("codex", q)
	if !again.shortWindow.resetAt.Equal(snapshot.shortWindow.resetAt) || !again.weekly.observedAt.Equal(now) {
		t.Fatal("cache read changed observation/reset")
	}
}

func TestQuotaNormalizeSharedDenialDoesNotUseModelOrCreditLimits(t *testing.T) {
	now := time.Unix(1800000000, 0)
	body := quotaTestBody("codex", now, time.Hour, 20)
	s, ok := selectionQuotaFromBody("codex", body, now)
	if !ok || s.weekly.remainingFraction != 0.8 {
		t.Fatal("additional model limit contaminated shared allowance")
	}
	body = []byte(strings.Replace(string(body), `"allowed":true`, `"allowed":false`, 1))
	s, ok = selectionQuotaFromBody("codex", body, now)
	if !ok || s.weekly.remainingFraction != 0 {
		t.Fatal("shared denial ignored")
	}
	q := quotaTestSignals("claude", now, time.Hour, 20)
	q.Signals["anthropic-ratelimit-unified-status"] = "rejected"
	q.Signals["anthropic-ratelimit-unified-7d_oi-status"] = "rejected"
	q.Signals["anthropic-ratelimit-unified-overage-status"] = "rejected"
	s, ok = selectionQuotaFromSignals("claude", q)
	if !ok || s.weekly.remainingFraction != 0.8 {
		t.Fatal("supplemental rejection contaminated shared allowance")
	}
	q.Signals["anthropic-ratelimit-unified-7d-status"] = "rejected"
	s, ok = selectionQuotaFromSignals("claude", q)
	if !ok || s.weekly.remainingFraction != 0 {
		t.Fatal("shared rejection ignored")
	}
}
