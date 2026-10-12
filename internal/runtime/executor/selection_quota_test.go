package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
)

type selectionQuotaTransport func(*http.Request) (*http.Response, error)

func (f selectionQuotaTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestSelectionQuotaRealExecutorsAndFollowup(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			cfg := &config.Config{}
			quota := cliproxyauth.DefaultQuotaSelectionConfig()
			s := cliproxyauth.NewSessionAffinitySelectorWithConfig(cliproxyauth.SessionAffinityConfig{QuotaSelection: &quota})
			defer s.Stop()
			m := cliproxyauth.NewManager(nil, s, nil)
			if provider == "claude" {
				m.RegisterExecutor(NewClaudeExecutor(cfg))
			} else {
				m.RegisterExecutor(NewCodexAutoExecutor(cfg))
			}
			for i := range 2 {
				_, err := m.Register(t.Context(), &cliproxyauth.Auth{ID: fmt.Sprintf("quota-executor-%d", i), Provider: provider, Attributes: map[string]string{cliproxyauth.AttributeAuthKind: cliproxyauth.AuthKindOAuth}, Metadata: map[string]any{"access_token": fmt.Sprintf("sanitized-token-%d", i), "account_id": "sanitized-account"}})
				if err != nil {
					t.Fatal(err)
				}
			}
			var calls atomic.Int32
			now := time.Now().Truncate(time.Second)
			rt := selectionQuotaTransport(func(r *http.Request) (*http.Response, error) {
				calls.Add(1)
				reset := 6 * time.Hour
				if r.Header.Get("Authorization") == "Bearer sanitized-token-1" {
					reset = 3 * time.Hour
				}
				var body []byte
				if provider == "claude" {
					body = fmt.Appendf(nil, `{"five_hour":{"utilization":5,"resets_at":%q},"seven_day":{"utilization":10,"resets_at":%q}}`, now.Add(time.Hour).Format(time.RFC3339), now.Add(reset).Format(time.RFC3339))
				} else {
					body = fmt.Appendf(nil, `{"rate_limit":{"primary_window":{"used_percent":10,"limit_window_seconds":604800,"reset_at":%d},"secondary_window":null}}`, now.Add(reset).Unix())
				}
				if r.Method != http.MethodGet {
					return nil, fmt.Errorf("unexpected test method")
				}
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(string(body))), Header: make(http.Header), Request: r}, nil
			})
			ctx := context.WithValue(t.Context(), "cliproxy.roundtripper", http.RoundTripper(rt))
			for range 2 {
				a, err := m.SelectAuth(ctx, provider, "", cliproxyexecutor.Options{Headers: http.Header{"X-Session-Id": {"real-executor-test-session"}}})
				if err != nil || a == nil || a.ID != "quota-executor-1" {
					t.Fatalf("real executor selection failed: %v", err)
				}
			}
			if calls.Load() != 2 {
				t.Fatalf("quota GETs=%d, want one per account then none on follow-up", calls.Load())
			}
			if _, ok := ctx.Deadline(); ok {
				t.Fatal("quota deadline leaked into inference context")
			}
		})
	}
}
