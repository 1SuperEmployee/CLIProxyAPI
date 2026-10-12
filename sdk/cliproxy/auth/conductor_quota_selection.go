package auth

import (
	"context"
	"errors"

	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
	cliproxysession "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/session"
)

type quotaSelectionPreparedKey struct{}

// Preparation happens with neither Manager.mu nor the affinity pick lock held.
// After any network wait the manager gathers eligibility again, and Pick still
// atomically checks the actual binding before consulting quota preferences.
func (m *Manager) prepareSelectionQuotas(ctx context.Context, selector Selector, provider, model string, opts cliproxyexecutor.Options, auths []*Auth) bool {
	s, ok := selector.(*SessionAffinitySelector)
	if !ok || s.quota == nil || ctx.Value(quotaSelectionPreparedKey{}) != nil {
		return false
	}
	if _, weighted := s.fallback.(*WeightedRoundRobinSelector); weighted {
		auths = positiveWeightAuths(auths)
	}
	pool := highestPriorityAuths(auths)
	if s.hasHealthyQuotaBinding(provider, model, opts, auths) {
		return false
	}
	pool = preferCodexWebsocketAuths(ctx, provider, pool)
	if !comparableQuotaAccounts(pool) {
		return false
	}
	return s.quota.refresh(ctx, pool, func(fetchCtx context.Context, a *Auth) ([]byte, error) {
		exec, exists := m.Executor(executorKeyFromAuth(a))
		fetcher, supported := exec.(quotaUsageFetcher)
		if !exists || !supported {
			return nil, errors.New("shared quota lookup unavailable")
		}
		if rt := m.roundTripperFor(a); rt != nil {
			fetchCtx = context.WithValue(fetchCtx, "cliproxy.roundtripper", rt)
		}
		// No independent OAuth refresh or persistence path here. If current
		// credentials cannot read usage, preserve normal selection/refresh.
		return fetcher.FetchSelectionQuota(fetchCtx, a)
	})
}

func comparableQuotaAccounts(pool []*Auth) bool {
	if len(pool) < 2 {
		return false
	}
	provider := selectionQuotaProvider(pool[0])
	if provider == "" {
		return false
	}
	for _, a := range pool[1:] {
		if selectionQuotaProvider(a) != provider {
			return false
		}
	}
	return true
}

func (s *SessionAffinitySelector) hasHealthyQuotaBinding(provider, model string, opts cliproxyexecutor.Options, available []*Auth) bool {
	primary, parent := extractExplicitSessionIDs(opts.Headers, opts.OriginalRequest, opts.Metadata)
	if primary == "" && s.matcher != nil {
		namespace := lcpAffinityNamespace(provider, model, opts.Metadata)
		turns := cliproxysession.ExtractCanonicalTurns(opts.SourceFormat, opts.OriginalRequest)
		if namespace != "" && len(turns) > 0 {
			fingerprints, minPrefix, tails, env := s.matcher.PrepareExt(turns)
			if minPrefix > 0 && minPrefix <= len(fingerprints) {
				if match, ok := s.matcher.MatchFingerprintsWithContext(namespace, fingerprints, tails, env, minPrefix); ok {
					for _, a := range available {
						if a.ID == match.AuthID {
							return true
						}
					}
				}
			}
		}
	}
	if primary == "" {
		primary, parent = extractSessionIDs(opts.Headers, opts.OriginalRequest, opts.Metadata)
	}
	if primary == "" {
		return false
	}
	isFork, _ := opts.Metadata[cliproxyexecutor.IsForkMetadataKey].(bool)
	inheritParent := isFork || s.subagentAffinity || !isSubagentSession(cliproxysession.BoundSessionIdentity(primary), cliproxysession.BoundSessionIdentity(parent))
	for _, id := range []string{primary, parent} {
		if id == "" || (id == parent && id != primary && !inheritParent) {
			continue
		}
		key := provider + "::" + cliproxysession.BoundSessionIdentity(id) + "::" + canonicalModelKey(model)
		if bound, ok := s.cache.Get(key); ok {
			for _, a := range available {
				if a.ID == bound {
					return true
				}
			}
			// Pick does not inherit the parent when its own binding exists but
			// is unavailable; it uses genuine failover in that branch.
			if id == primary {
				return false
			}
		}
	}
	return false
}

func (s *SessionAffinitySelector) pickFallback(ctx context.Context, provider, model string, opts cliproxyexecutor.Options, pool []*Auth) (*Auth, error) {
	if s.quota != nil {
		// Preserve the baseline transport preference before quota narrows the
		// pool; applying it afterward could force a nearer HTTP-only account.
		pool = preferCodexWebsocketAuths(ctx, provider, pool)
	}
	if s.quota != nil && comparableQuotaAccounts(pool) {
		decision := nearestResetCandidates(pool, "", s.quota.snapshots(pool), s.quota.now(), s.quota.policy())
		// Reason and candidate counts only: no quota bodies, tokens, account
		// emails or prompts. Existing affinity logs record the chosen auth ID.
		selectorLogEntry(ctx).WithField("quota_selection", decision.reason).WithField("eligible", len(pool)).WithField("preferred", len(decision.candidates)).Debug("shared quota selection")
		pool = decision.candidates
	}
	return s.fallback.Pick(ctx, provider, model, opts, pool)
}
