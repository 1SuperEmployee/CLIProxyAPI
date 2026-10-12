package auth

import (
	"math"
	"time"
)

// These types and the decision helper are a qualification prototype, not a
// registered Selector. No live routing, quota reads or binding writes use them.
// Provider normalization must establish every applicable scope before setting
// scopesComplete. Primary/secondary names alone do not identify window duration.
type selectionQuotaWindow struct {
	id                string
	duration          time.Duration
	remainingFraction float64
	resetAt           time.Time
	observedAt        time.Time
}

type selectionQuotaSnapshot struct {
	scopesComplete bool
	// rankingWindow identifies the relevant weekly allowance explicitly. Other
	// applicable windows also constrain headroom but do not redefine the ranking.
	rankingWindow string
	windows       []selectionQuotaWindow
}

type quotaSelectionPolicy struct {
	maxAge               time.Duration
	minRemainingFraction float64
}

type quotaSelectionDecision struct {
	candidates []*Auth
	reason     string
}

// nearestResetCandidates narrows an already eligible, ID-sorted candidate pool.
// The caller owns provider/model eligibility, explicit pins, cooldowns, transport
// preference and weight exclusions. A healthy binding must be checked across
// priority tiers, before the cold-selection pool is narrowed to the highest tier.
// The returned candidates retain input order for the existing tie-break selector.
//
// It makes no I/O, uses the supplied clock, does not mutate inputs, and must not
// be called in place of the atomic binding lookup/recheck. The healthyBoundID
// input models that lookup for policy tests; integration remains separate.
func nearestResetCandidates(eligible []*Auth, healthyBoundID string, snapshots map[string]selectionQuotaSnapshot, now time.Time, policy quotaSelectionPolicy) quotaSelectionDecision {
	if healthyBoundID != "" {
		for _, candidate := range eligible {
			if candidate.ID == healthyBoundID {
				return quotaSelectionDecision{[]*Auth{candidate}, "healthy-binding"}
			}
		}
	}
	pool := highestPriorityAuths(eligible)
	fallback := func(reason string) quotaSelectionDecision {
		return quotaSelectionDecision{pool, reason}
	}
	if len(pool) == 0 {
		return fallback("no-candidates")
	}
	if len(pool) == 1 {
		return fallback("single-candidate")
	}
	if now.IsZero() || policy.maxAge <= 0 || !validFraction(policy.minRemainingFraction) || policy.minRemainingFraction == 0 {
		return fallback("invalid-policy")
	}

	var bestReset time.Time
	var best []*Auth
	for _, candidate := range pool {
		snapshot, exists := snapshots[candidate.ID]
		if !exists {
			return fallback("missing-quota")
		}
		reset, hasHeadroom, reason := assessSelectionQuota(snapshot, now, policy)
		if reason != "" {
			// An unknown candidate is not an exhausted candidate. Refresh/resolve
			// its data before narrowing the pool; otherwise preserve baseline.
			return fallback(reason)
		}
		if !hasHeadroom {
			continue
		}
		if bestReset.IsZero() || reset.Before(bestReset) {
			bestReset = reset
			best = []*Auth{candidate}
		} else if reset.Equal(bestReset) {
			best = append(best, candidate)
		}
	}
	if len(best) == 0 {
		// A reserve is a preference, not a provider denial. Leave true quota
		// rejection/cooldown enforcement to the authoritative eligibility path.
		return fallback("insufficient-headroom")
	}
	return quotaSelectionDecision{best, "nearest-reset"}
}

func assessSelectionQuota(snapshot selectionQuotaSnapshot, now time.Time, policy quotaSelectionPolicy) (time.Time, bool, string) {
	if !snapshot.scopesComplete {
		return time.Time{}, false, "unresolved-quota-scope"
	}
	if snapshot.rankingWindow == "" || len(snapshot.windows) == 0 {
		return time.Time{}, false, "incomplete-quota"
	}
	seen := make(map[string]struct{}, len(snapshot.windows))
	var rankReset time.Time
	hasHeadroom := true
	for _, window := range snapshot.windows {
		if _, duplicate := seen[window.id]; duplicate || window.id == "" || window.duration <= 0 || !validFraction(window.remainingFraction) {
			return time.Time{}, false, "invalid-quota"
		}
		seen[window.id] = struct{}{}
		if window.observedAt.IsZero() || window.observedAt.After(now) {
			return time.Time{}, false, "invalid-quota-time"
		}
		// Exact boundaries expire. Reading a cached value never extends its age,
		// and a passed reset never implies a fresh, fully available allowance.
		if now.Sub(window.observedAt) >= policy.maxAge || !window.resetAt.After(now) {
			return time.Time{}, false, "stale-quota"
		}
		if window.resetAt.After(window.observedAt.Add(window.duration)) {
			return time.Time{}, false, "invalid-quota-time"
		}
		if window.remainingFraction < policy.minRemainingFraction {
			hasHeadroom = false
		}
		if window.id == snapshot.rankingWindow {
			if window.duration != 7*24*time.Hour {
				return time.Time{}, false, "invalid-ranking-window"
			}
			rankReset = window.resetAt
		}
	}
	if rankReset.IsZero() {
		return time.Time{}, false, "incomplete-quota"
	}
	return rankReset, hasHeadroom, ""
}

func validFraction(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 1
}
