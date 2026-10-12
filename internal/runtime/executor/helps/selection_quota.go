package helps

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

const maxSelectionQuotaBody = 1 << 20

// FetchSharedSelectionQuota performs exactly one read-only usage GET through
// the existing provider transport. It has no refresh, logging or credit-spend
// path. Its selection-stage deadline never applies to inference streams.
func FetchSharedSelectionQuota(ctx context.Context, cfg *config.Config, auth *cliproxyauth.Auth, prepare func(*http.Request, *cliproxyauth.Auth) error) ([]byte, error) {
	if auth == nil || auth.AuthKind() != cliproxyauth.AuthKindOAuth {
		return nil, errors.New("shared quota lookup requires OAuth")
	}
	if _, bounded := ctx.Deadline(); !bounded {
		return nil, errors.New("shared quota lookup requires a deadline")
	}
	var target string
	switch auth.Provider {
	case "claude":
		target = "https://api.anthropic.com/api/oauth/usage"
	case "codex":
		target = "https://chatgpt.com/backend-api/wham/usage"
	default:
		return nil, errors.New("shared quota lookup unsupported provider")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, errors.New("shared quota request unavailable")
	}
	if err = prepare(req, auth); err != nil {
		return nil, errors.New("shared quota credential preparation failed")
	}
	req.Header.Set("Accept", "application/json")
	if auth.Provider == "claude" {
		req.Header.Set("Anthropic-Beta", "oauth-2025-04-20")
		req.Header.Set("User-Agent", defaultClaudeDeviceProfile(cfg).UserAgent)
	} else if account, _ := auth.Metadata["account_id"].(string); strings.TrimSpace(account) != "" {
		req.Header.Set("Chatgpt-Account-Id", account)
	}
	client := NewUtlsHTTPClient(ctx, cfg, auth, 0)
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := client.Do(req)
	if err != nil {
		// Do not return URL errors or transport errors that may embed headers.
		return nil, errors.New("shared quota transport failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shared quota upstream status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxSelectionQuotaBody+1))
	if err != nil || len(body) > maxSelectionQuotaBody {
		return nil, errors.New("shared quota body unavailable or oversized")
	}
	return body, nil
}
