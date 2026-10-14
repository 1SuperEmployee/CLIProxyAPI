package auth

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	internallogging "github.com/router-for-me/CLIProxyAPI/v8/internal/logging"
)

func TestQuotaCaptureTimeLongStreamAndPartialEvents(t *testing.T) {
	old := time.Unix(1800000000, 0)
	now := old.Add(20 * time.Minute)
	headers := http.Header{"X-Codex-Primary-Used-Percent": {"20"}, "X-Codex-Secondary-Used-Percent": {"30"}, "Content-Type": {"text/event-stream"}}
	times := map[string]time.Time{"X-Codex-Primary-Used-Percent": now, "X-Codex-Secondary-Used-Percent": old, "Content-Type": old.Add(-time.Hour)}
	if !quotaObservationTime("codex", headers, times, now).Equal(old) {
		t.Fatal("partial event or completion freshened old quota fields")
	}
	delete(headers, "X-Codex-Secondary-Used-Percent")
	if !quotaObservationTime("codex", headers, times, now).Equal(now) {
		t.Fatal("non-quota field incorrectly aged a fresh observation")
	}
	q := QuotaState{}
	q.ObserveResponseHeadersForProvider("codex", headers, now)
	headers.Set("X-Codex-Primary-Used-Percent", "1")
	if q.ObserveResponseHeadersForProvider("codex", headers, old) || q.Signals["X-Codex-Primary-Used-Percent"] != "20" {
		t.Fatal("older completed stream overwrote newer account usage")
	}
}

func TestQuotaManagerRecordsArrivalBeforeCompletion(t *testing.T) {
	m := NewManager(nil, &RoundRobinSelector{}, nil)
	a, err := m.Register(t.Context(), quotaTestAuth("codex", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	ctx := internallogging.WithResponseHeadersHolder(t.Context())
	internallogging.SetResponseHeaders(ctx, http.Header{"X-Codex-Primary-Used-Percent": {"20"}})
	_, times := internallogging.GetResponseHeadersWithCaptureTimes(ctx)
	captured := times["X-Codex-Primary-Used-Percent"]
	m.MarkResult(ctx, Result{AuthID: a.ID, Provider: "codex", Model: "quota-test-model", Success: true})
	current, _ := m.GetByID(a.ID)
	if !current.Quota.ObservedAt.Equal(captured) {
		t.Fatalf("observation timestamp was completion rather than arrival: got=%v want=%v", current.Quota.ObservedAt, captured)
	}
}

func TestQuotaRelativeResetUsesOwnArrivalWithoutFresheningSnapshot(t *testing.T) {
	old := time.Unix(1800000000, 0)
	arrival := old.Add(5 * time.Minute)
	for _, oldAbsolute := range []bool{false, true} {
		t.Run(strconv.FormatBool(oldAbsolute), func(t *testing.T) {
			headers := http.Header{
				"X-Codex-Plan-Type": {"synthetic-plan"}, "X-Codex-Primary-Window-Minutes": {"10080"},
				"X-Codex-Primary-Used-Percent": {"20"}, "X-Codex-Primary-Reset-After-Seconds": {"3600"},
			}
			times := map[string]time.Time{}
			for name := range headers {
				times[name] = arrival
			}
			times["X-Codex-Plan-Type"] = old
			if oldAbsolute {
				headers.Set("X-Codex-Primary-Reset-At", strconv.FormatInt(old.Add(time.Hour).Unix(), 10))
				times["X-Codex-Primary-Reset-At"] = old
			}
			anchorSharedCodexRelativeResets(headers, times)
			observedAt := quotaObservationTime("codex", headers, times, arrival)
			if !observedAt.Equal(old) {
				t.Fatal("anchoring relative reset freshened an older retained field")
			}
			q := QuotaState{}
			q.ObserveResponseHeadersForProvider("codex", headers, observedAt)
			snapshot, valid := selectionQuotaFromSignals("codex", q)
			if !valid || !snapshot.weekly.resetAt.Equal(arrival.Add(time.Hour)) {
				t.Fatalf("relative reset moved to older snapshot time: valid=%v reset=%v", valid, snapshot.weekly.resetAt)
			}
			// An absolute reset arriving with the same event wins even if the
			// relative field disagrees; do not rewrite provider data.
			absolute := arrival.Add(2 * time.Hour).Unix()
			headers.Set("X-Codex-Primary-Reset-At", strconv.FormatInt(absolute, 10))
			times["X-Codex-Primary-Reset-At"] = arrival
			anchorSharedCodexRelativeResets(headers, times)
			if headers.Get("X-Codex-Primary-Reset-At") != strconv.FormatInt(absolute, 10) {
				t.Fatal("relative reset overrode an equally recent absolute reset")
			}
		})
	}
}
