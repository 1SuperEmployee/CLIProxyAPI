package logging

import (
	"net/http"
	"testing"
	"time"
)

func TestResponseHeaderCaptureTimesRemainPerFieldAndIsolated(t *testing.T) {
	ctx := WithResponseHeadersHolder(t.Context())
	before := time.Unix(1800000000, 0)
	later := before.Add(20 * time.Minute)
	setResponseHeadersAt(ctx, http.Header{"X-Codex-Primary-Used-Percent": {"20"}, "X-Codex-Secondary-Used-Percent": {"30"}}, before)
	mergeResponseHeadersAt(ctx, http.Header{"X-Codex-Primary-Used-Percent": {"40"}}, later)
	headers, times := GetResponseHeadersWithCaptureTimes(ctx)
	if headers.Get("X-Codex-Primary-Used-Percent") != "40" || !times["X-Codex-Primary-Used-Percent"].Equal(later) || !times["X-Codex-Secondary-Used-Percent"].Equal(before) {
		t.Fatal("partial event freshened an untouched quota field")
	}
	delete(times, "X-Codex-Secondary-Used-Percent")
	headers.Set("X-Codex-Primary-Used-Percent", "99")
	headers, times = GetResponseHeadersWithCaptureTimes(ctx)
	if headers.Get("X-Codex-Primary-Used-Percent") != "40" || !times["X-Codex-Secondary-Used-Percent"].Equal(before) {
		t.Fatal("returned metadata aliases mutable holder state")
	}
	fresh := WithFreshResponseHeadersHolder(ctx)
	if h, times := GetResponseHeadersWithCaptureTimes(fresh); len(h) != 0 || len(times) != 0 {
		t.Fatal("retry inherited previous attempt's capture times")
	}
	setResponseHeadersAt(ctx, http.Header{"X-Codex-Primary-Used-Percent": {"50"}}, later)
	_, times = GetResponseHeadersWithCaptureTimes(ctx)
	if len(times) != 1 {
		t.Fatal("replacing response headers retained old capture times")
	}
}
