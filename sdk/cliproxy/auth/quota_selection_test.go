package auth

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func selectionTestSnapshot(now time.Time, resetAfter time.Duration, remaining float64) selectionQuotaSnapshot {
	return selectionQuotaSnapshot{
		weekly: selectionQuotaWindow{
			duration:          7 * 24 * time.Hour,
			remainingFraction: remaining, resetAt: now.Add(resetAfter), observedAt: now,
		},
	}
}

func selectionTestIDs(decision quotaSelectionDecision) []string {
	ids := make([]string, 0, len(decision.candidates))
	for _, candidate := range decision.candidates {
		ids = append(ids, candidate.ID)
	}
	return ids
}

func TestQuotaSelectionPreservesBindingAndPriority(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	pool := []*Auth{
		{ID: "a", Attributes: map[string]string{"priority": "1"}},
		{ID: "b", Attributes: map[string]string{"priority": "10"}},
		{ID: "c", Attributes: map[string]string{"priority": "10"}},
	}
	// Existing bindings outrank cold-selection priority and do not need quota
	// or even a usable policy. This models a healthy long-lived conversation.
	bound := nearestResetCandidates(pool, "a", nil, now, quotaSelectionPolicy{})
	if got := selectionTestIDs(bound); !reflect.DeepEqual(got, []string{"a"}) || bound.reason != "healthy-binding" {
		t.Fatalf("healthy binding: ids=%v reason=%s", got, bound.reason)
	}
	policy := quotaSelectionPolicy{15 * time.Minute, 0.05}
	snapshots := map[string]selectionQuotaSnapshot{
		"b": selectionTestSnapshot(now, 2*time.Hour, 0.5),
		"c": selectionTestSnapshot(now, time.Hour, 0.5),
	}
	// An absent binding is not restored, and missing lower-priority data must
	// not trigger a fallback or expand the cold-selection tier.
	decision := nearestResetCandidates(pool, "no-longer-eligible", snapshots, now, policy)
	if got := selectionTestIDs(decision); !reflect.DeepEqual(got, []string{"c"}) || decision.reason != "quota-pace" {
		t.Fatalf("cold priority: ids=%v reason=%s", got, decision.reason)
	}
}

func TestQuotaSelectionSharedWindows(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	policy := quotaSelectionPolicy{15 * time.Minute, 0.05}
	pool := []*Auth{{ID: "a"}, {ID: "b"}}
	t.Run("weekly-only Codex does not need an invented secondary window", func(t *testing.T) {
		snapshots := map[string]selectionQuotaSnapshot{
			"a": selectionTestSnapshot(now, time.Hour, 0.5),
			"b": selectionTestSnapshot(now, 2*time.Hour, 0.95),
		}
		decision := nearestResetCandidates(pool, "", snapshots, now, policy)
		if got := selectionTestIDs(decision); !reflect.DeepEqual(got, []string{"a"}) {
			t.Fatalf("remaining allowance per time should win: %v (%s)", got, decision.reason)
		}
	})
	t.Run("short window constrains headroom but not ranking", func(t *testing.T) {
		a := selectionTestSnapshot(now, time.Hour, 0.5)
		a.shortWindow = &selectionQuotaWindow{
			duration:          5 * time.Hour,
			remainingFraction: 0.01, resetAt: now.Add(time.Minute), observedAt: now,
		}
		b := selectionTestSnapshot(now, 2*time.Hour, 0.5)
		snapshots := map[string]selectionQuotaSnapshot{"a": a, "b": b}
		decision := nearestResetCandidates(pool, "", snapshots, now, policy)
		if got := selectionTestIDs(decision); !reflect.DeepEqual(got, []string{"b"}) {
			t.Fatalf("nearer reset cannot bypass short-window reserve: %v (%s)", got, decision.reason)
		}
	})
	t.Run("short-window reset cannot displace shared-weekly ranking", func(t *testing.T) {
		a := selectionTestSnapshot(now, time.Hour, 0.5)
		b := selectionTestSnapshot(now, 2*time.Hour, 0.5)
		a.shortWindow = &selectionQuotaWindow{
			duration:          5 * time.Hour,
			remainingFraction: 0.5, resetAt: now.Add(3 * time.Hour), observedAt: now,
		}
		b.shortWindow = &selectionQuotaWindow{
			duration:          5 * time.Hour,
			remainingFraction: 0.5, resetAt: now.Add(time.Minute), observedAt: now,
		}
		decision := nearestResetCandidates(pool, "", map[string]selectionQuotaSnapshot{"a": a, "b": b}, now, policy)
		if got := selectionTestIDs(decision); !reflect.DeepEqual(got, []string{"a"}) || decision.reason != "quota-pace" {
			t.Fatalf("short reset cannot replace weekly allowance pacing: %v (%s)", got, decision.reason)
		}
	})
}

func TestQuotaSelectionFallsBackForUnknownCandidate(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	policy := quotaSelectionPolicy{15 * time.Minute, 0.05}
	pool := []*Auth{{ID: "a"}, {ID: "b"}}
	tests := []struct {
		name   string
		edit   func(*selectionQuotaSnapshot)
		reason string
	}{
		{"missing weekly window", func(s *selectionQuotaSnapshot) { s.weekly = selectionQuotaWindow{} }, "incomplete-quota"},
		{"age boundary", func(s *selectionQuotaSnapshot) { s.weekly.observedAt = now.Add(-policy.maxAge) }, "stale-quota"},
		{"passed reset", func(s *selectionQuotaSnapshot) { s.weekly.resetAt = now.Add(-time.Second) }, "stale-quota"},
		{"reset boundary", func(s *selectionQuotaSnapshot) { s.weekly.resetAt = now }, "stale-quota"},
		{"absent observation", func(s *selectionQuotaSnapshot) { s.weekly.observedAt = time.Time{} }, "invalid-quota-time"},
		{"future observation", func(s *selectionQuotaSnapshot) { s.weekly.observedAt = now.Add(time.Second) }, "invalid-quota-time"},
		{"epoch milliseconds mistaken for seconds", func(s *selectionQuotaSnapshot) { s.weekly.resetAt = time.Unix(now.Unix()*1000, 0) }, "invalid-quota-time"},
		{"percentage points mistaken for fractions", func(s *selectionQuotaSnapshot) { s.weekly.remainingFraction = 80 }, "invalid-quota"},
		{"negative fraction", func(s *selectionQuotaSnapshot) { s.weekly.remainingFraction = -0.1 }, "invalid-quota"},
		{"nonfinite fraction", func(s *selectionQuotaSnapshot) { s.weekly.remainingFraction = math.NaN() }, "invalid-quota"},
		{"infinite fraction", func(s *selectionQuotaSnapshot) { s.weekly.remainingFraction = math.Inf(1) }, "invalid-quota"},
		{"absent weekly duration", func(s *selectionQuotaSnapshot) { s.weekly.duration = 0 }, "incomplete-quota"},
		{"short primary assumed weekly", func(s *selectionQuotaSnapshot) { s.weekly.duration = 5 * time.Hour }, "invalid-ranking-window"},
		{"invalid short duration", func(s *selectionQuotaSnapshot) { s.shortWindow = &selectionQuotaWindow{} }, "invalid-quota"},
		{"weekly window supplied as short", func(s *selectionQuotaSnapshot) { s.shortWindow = &s.weekly }, "invalid-short-window"},
		{"one stale window in otherwise fresh snapshot", func(s *selectionQuotaSnapshot) {
			s.shortWindow = &selectionQuotaWindow{
				duration:          5 * time.Hour,
				remainingFraction: 0.5, resetAt: now.Add(time.Hour), observedAt: now.Add(-policy.maxAge),
			}
		}, "stale-quota"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := selectionTestSnapshot(now, time.Hour, 0.5)
			tc.edit(&a)
			snapshots := map[string]selectionQuotaSnapshot{
				"a": a, "b": selectionTestSnapshot(now, 2*time.Hour, 0.5),
			}
			decision := nearestResetCandidates(pool, "", snapshots, now, policy)
			if got := selectionTestIDs(decision); !reflect.DeepEqual(got, []string{"a", "b"}) || decision.reason != tc.reason {
				t.Fatalf("unknown candidate must preserve baseline pool: ids=%v reason=%s; want %s", got, decision.reason, tc.reason)
			}
		})
	}
	missing := nearestResetCandidates(pool, "", map[string]selectionQuotaSnapshot{
		"b": selectionTestSnapshot(now, time.Hour, 0.5),
	}, now, policy)
	if got := selectionTestIDs(missing); !reflect.DeepEqual(got, []string{"a", "b"}) || missing.reason != "missing-quota" {
		t.Fatalf("missing candidate cannot be silently dropped: %v (%s)", got, missing.reason)
	}
	for _, remaining := range []float64{0.01, 0.8} {
		unknownLast := nearestResetCandidates(pool, "", map[string]selectionQuotaSnapshot{
			"a": selectionTestSnapshot(now, time.Hour, remaining),
		}, now, policy)
		if unknownLast.reason != "missing-quota" || len(unknownLast.candidates) != 2 {
			t.Fatal("a known low or attractive account cannot hide an unknown later candidate")
		}
	}
}

func TestQuotaSelectionTiesReservesAndInputs(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	policy := quotaSelectionPolicy{15 * time.Minute, 0.05}
	pool := []*Auth{{ID: "a"}, {ID: "b"}, {ID: "c"}}
	snapshots := map[string]selectionQuotaSnapshot{
		"a": selectionTestSnapshot(now, time.Hour, 0.05),
		"b": selectionTestSnapshot(now, time.Hour, 0.05),
		"c": selectionTestSnapshot(now, 2*time.Hour, 0.05),
	}
	decision := nearestResetCandidates(pool, "", snapshots, now, policy)
	if got := selectionTestIDs(decision); !reflect.DeepEqual(got, []string{"a", "b"}) || decision.reason != "quota-pace" {
		t.Fatalf("reserve boundary and tied resets should preserve baseline tie order: %v (%s)", got, decision.reason)
	}
	// A later decision uses the original observation time; an earlier read
	// cannot refresh the cache implicitly, nor can selection reorder inputs.
	expired := nearestResetCandidates(pool, "", snapshots, now.Add(policy.maxAge), policy)
	if expired.reason != "stale-quota" || !snapshots["a"].weekly.observedAt.Equal(now) || pool[0].ID != "a" {
		t.Fatal("reading quota changed its age or candidate ordering")
	}
	for id, snapshot := range snapshots {
		snapshot.weekly.remainingFraction = 0.01
		snapshots[id] = snapshot
	}
	low := nearestResetCandidates(pool, "", snapshots, now, policy)
	if got := selectionTestIDs(low); !reflect.DeepEqual(got, []string{"a", "b", "c"}) || low.reason != "insufficient-headroom" {
		t.Fatalf("reserve should not invent a provider denial: %v (%s)", got, low.reason)
	}
}

func TestQuotaSelectionPacesRemainingWeeklyAllowance(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	pool := []*Auth{{ID: "a"}, {ID: "b"}}
	policy := quotaSelectionPolicy{15 * time.Minute, 0.05}
	for _, tc := range []struct {
		name   string
		a, b   float64
		ta, tb time.Duration
		want   []string
	}{
		{"observed imbalance", .09, .97, 48 * time.Hour, 73 * time.Hour, []string{"b"}},
		{"equal balances favour sooner reset", .8, .8, 24 * time.Hour, 72 * time.Hour, []string{"a"}},
		{"expiring allowance can beat fuller account", .25, .95, 6 * time.Hour, 120 * time.Hour, []string{"a"}},
		{"imminent reset remains useful", .06, .95, time.Nanosecond, 72 * time.Hour, []string{"a"}},
		{"equal rates preserve fallback tie", .25, .5, 24 * time.Hour, 48 * time.Hour, []string{"a", "b"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			snapshots := map[string]selectionQuotaSnapshot{
				"a": selectionTestSnapshot(now, tc.ta, tc.a),
				"b": selectionTestSnapshot(now, tc.tb, tc.b),
			}
			got := nearestResetCandidates(pool, "", snapshots, now, policy)
			if !reflect.DeepEqual(selectionTestIDs(got), tc.want) || got.reason != "quota-pace" {
				t.Fatalf("got %v (%s), want %v", selectionTestIDs(got), got.reason, tc.want)
			}
		})
	}
}

func TestQuotaSelectionDegeneratePoolsAndPolicy(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	policy := quotaSelectionPolicy{15 * time.Minute, 0.05}
	empty := nearestResetCandidates(nil, "", nil, now, policy)
	if empty.reason != "no-candidates" || len(empty.candidates) != 0 {
		t.Fatalf("empty pool: %+v", empty)
	}
	single := nearestResetCandidates([]*Auth{{ID: "a"}}, "", nil, now, policy)
	if single.reason != "single-candidate" || len(single.candidates) != 1 {
		t.Fatalf("single eligible account needs no ranking lookup: %+v", single)
	}
	for _, invalid := range []quotaSelectionPolicy{
		{}, {time.Minute, 0}, {-time.Minute, 0.05}, {time.Minute, 1.1},
		{time.Minute, math.NaN()}, {time.Minute, math.Inf(1)},
	} {
		decision := nearestResetCandidates([]*Auth{{ID: "a"}, {ID: "b"}}, "", nil, now, invalid)
		if decision.reason != "invalid-policy" || len(decision.candidates) != 2 {
			t.Fatalf("invalid policy must preserve eligible candidates: %+v", decision)
		}
	}
}

func BenchmarkQuotaSelection(b *testing.B) {
	now := time.Unix(1_800_000_000, 0)
	policy := quotaSelectionPolicy{15 * time.Minute, 0.05}
	pool := []*Auth{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}}
	snapshots := make(map[string]selectionQuotaSnapshot, len(pool))
	for i, candidate := range pool {
		snapshots[candidate.ID] = selectionTestSnapshot(now, time.Duration(i+1)*time.Hour, 0.5)
	}
	b.Run("healthy-binding", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			nearestResetCandidates(pool, "b", nil, now, policy)
		}
	})
	b.Run("four-fresh-accounts", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			nearestResetCandidates(pool, "", snapshots, now, policy)
		}
	})
}
