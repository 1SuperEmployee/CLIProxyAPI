package auth

import (
	"math"
	"time"
)

// These types support the opt-in shared-quota preference inside session affinity.
// Only shared account windows participate in ranking. Supplemental model limits
// remain enforced by existing provider rejection/model cooldown handling.
type selectionQuotaWindow struct {
	duration          time.Duration
	remainingFraction float64
	resetAt           time.Time
	observedAt        time.Time
}

type selectionQuotaSnapshot struct {
	weekly selectionQuotaWindow
	// nil means the provider reports no active shared short window, not that a
	// required measurement is missing. Normalization must distinguish the two.
	// Primary/secondary names alone do not identify window duration.
	shortWindow *selectionQuotaWindow
}

type quotaSelectionPolicy struct {
	maxAge               time.Duration
	minRemainingFraction float64
}

type quotaSelectionDecision struct {
	candidates []*Auth
	reason     string
}

// nearestResetCandidates ranks unused weekly allowance against time until reset
// inside the legacy nearest-reset mode. It narrows an eligible, ID-sorted pool.
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

	var bestScore float64
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
		// Validation above guarantees a positive time and a finite fraction.
		// Equal scores retain the existing selector's tie-breaking order.
		score := snapshot.weekly.remainingFraction / reset.Sub(now).Seconds()
		if len(best) == 0 || score > bestScore {
			bestScore = score
			best = []*Auth{candidate}
		} else if score == bestScore {
			best = append(best, candidate)
		}
	}
	if len(best) == 0 {
		// A reserve is a preference, not a provider denial. Leave true quota
		// rejection/cooldown enforcement to the authoritative eligibility path.
		return fallback("insufficient-headroom")
	}
	return quotaSelectionDecision{best, "quota-pace"}
}

func assessSelectionQuota(snapshot selectionQuotaSnapshot, now time.Time, policy quotaSelectionPolicy) (time.Time, bool, string) {
	if snapshot.weekly.duration == 0 {
		return time.Time{}, false, "incomplete-quota"
	}
	if snapshot.weekly.duration != 7*24*time.Hour {
		return time.Time{}, false, "invalid-ranking-window"
	}
	if reason := validateSelectionQuotaWindow(snapshot.weekly, now, policy); reason != "" {
		return time.Time{}, false, reason
	}
	hasHeadroom := snapshot.weekly.remainingFraction >= policy.minRemainingFraction
	if window := snapshot.shortWindow; window != nil {
		if window.duration >= snapshot.weekly.duration {
			return time.Time{}, false, "invalid-short-window"
		}
		if reason := validateSelectionQuotaWindow(*window, now, policy); reason != "" {
			return time.Time{}, false, reason
		}
		hasHeadroom = hasHeadroom && window.remainingFraction >= policy.minRemainingFraction
	}
	return snapshot.weekly.resetAt, hasHeadroom, ""
}

func validateSelectionQuotaWindow(window selectionQuotaWindow, now time.Time, policy quotaSelectionPolicy) string {
	if window.duration <= 0 || !validFraction(window.remainingFraction) {
		return "invalid-quota"
	}
	if window.observedAt.IsZero() || window.observedAt.After(now) {
		return "invalid-quota-time"
	}
	// Exact boundaries expire. Reading a cached value never extends its age,
	// and a passed reset never implies a fresh, fully available allowance.
	if now.Sub(window.observedAt) >= policy.maxAge || !window.resetAt.After(now) {
		return "stale-quota"
	}
	if window.resetAt.After(window.observedAt.Add(window.duration)) {
		return "invalid-quota-time"
	}
	return ""
}

func validFraction(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 1
}
