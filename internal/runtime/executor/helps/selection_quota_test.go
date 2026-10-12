package helps

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

type quotaReadCloser struct {
	io.Reader
	closed bool
}

func (b *quotaReadCloser) Close() error { b.closed = true; return nil }

func TestSelectionQuotaGETIsBoundedReadOnlyAndRejectsRedirects(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		for _, status := range []int{200, 302, 429} {
			t.Run(provider+http.StatusText(status), func(t *testing.T) {
				var calls int
				body := &quotaReadCloser{Reader: strings.NewReader(`{"sanitized":true}`)}
				rt := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					calls++
					if req.Method != http.MethodGet || req.Body != nil || req.URL.Scheme != "https" {
						t.Fatal("not a read-only usage GET")
					}
					want := "https://api.anthropic.com/api/oauth/usage"
					if provider == "codex" {
						want = "https://chatgpt.com/backend-api/wham/usage"
					}
					if req.URL.String() != want {
						t.Fatal("unexpected target")
					}
					if _, ok := req.Context().Deadline(); !ok {
						t.Fatal("deadline lost")
					}
					if req.Header.Get("Authorization") != "Bearer sanitized-test-credential" {
						t.Fatal("credential preparation not used")
					}
					if provider == "claude" && req.Header.Get("Anthropic-Beta") != "oauth-2025-04-20" {
						t.Fatal("OAuth beta missing")
					}
					if provider == "codex" && req.Header.Get("Chatgpt-Account-Id") != "sanitized-account" {
						t.Fatal("account header missing")
					}
					return &http.Response{StatusCode: status, Body: body, Header: http.Header{"Location": {"https://invalid.example/credential-leak"}}, Request: req}, nil
				})
				base := context.WithValue(t.Context(), "cliproxy.roundtripper", http.RoundTripper(rt))
				ctx, cancel := context.WithTimeout(base, time.Second)
				defer cancel()
				a := &cliproxyauth.Auth{Provider: provider, Attributes: map[string]string{cliproxyauth.AttributeAuthKind: cliproxyauth.AuthKindOAuth}, Metadata: map[string]any{"account_id": "sanitized-account"}}
				data, err := FetchSharedSelectionQuota(ctx, nil, a, func(req *http.Request, _ *cliproxyauth.Auth) error {
					req.Header.Set("Authorization", "Bearer sanitized-test-credential")
					return nil
				})
				if calls != 1 || !body.closed {
					t.Fatal("retried/redirected or leaked response body")
				}
				if status == 200 && (err != nil || len(data) == 0) {
					t.Fatalf("successful body unavailable: %v", err)
				}
				if status != 200 && (err == nil || len(data) != 0) {
					t.Fatal("upstream error treated as successful quota")
				}
			})
		}
	}
}

func TestSelectionQuotaRejectsOversizedBodyAndUnboundedContext(t *testing.T) {
	a := &cliproxyauth.Auth{Provider: "codex", Attributes: map[string]string{cliproxyauth.AttributeAuthKind: cliproxyauth.AuthKindOAuth}}
	prepare := func(*http.Request, *cliproxyauth.Auth) error { return nil }
	if _, err := FetchSharedSelectionQuota(t.Context(), nil, a, prepare); err == nil {
		t.Fatal("accepted unbounded lookup")
	}
	body := &quotaReadCloser{Reader: strings.NewReader(strings.Repeat("x", maxSelectionQuotaBody+1))}
	rt := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: body, Header: make(http.Header), Request: req}, nil
	})
	ctx, cancel := context.WithTimeout(context.WithValue(t.Context(), "cliproxy.roundtripper", http.RoundTripper(rt)), time.Second)
	defer cancel()
	if data, err := FetchSharedSelectionQuota(ctx, nil, a, prepare); err == nil || len(data) != 0 || !body.closed {
		t.Fatal("oversized quota response accepted or body leaked")
	}
}
