package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestUncertainCodexRefreshKeepsValidAccessWithoutRepeatingGrant(t *testing.T) {
	now := time.Now()
	expiry := now.Add(24 * time.Hour)
	executor := &unauthorizedRefreshExecutor{id: "codex", refreshErr: errors.New("refresh_outcome_unknown: token refresh response was lost; sign in again before another renewal")}
	m := NewManager(nil, &FillFirstSelector{}, nil)
	m.RegisterExecutor(executor)
	auth := &Auth{ID: "uncertain-refresh-test", Provider: "codex", Status: StatusActive, Metadata: map[string]any{
		"access_token": "synthetic-access", "refresh_token": "synthetic-refresh", "expired": expiry.Format(time.RFC3339),
	}}
	if _, err := m.Register(context.Background(), auth); err != nil {
		t.Fatal(err)
	}
	m.refreshAuth(context.Background(), auth.ID)
	current, _ := m.GetByID(auth.ID)
	if current == nil || current.Unavailable || current.Status != StatusActive || !current.HasValidAccessToken(now) {
		t.Fatal("uncertain renewal disabled usable access")
	}
	if !current.RequiresSignIn || current.RefreshError == nil || current.RefreshError.Code != "refresh_outcome_unknown" {
		t.Fatalf("uncertain renewal not recorded: %+v", current.RefreshError)
	}
	if current.Metadata["access_token"] != "synthetic-access" || current.Metadata["refresh_token"] != "synthetic-refresh" {
		t.Fatal("existing credentials changed")
	}
	m.refreshAuth(context.Background(), auth.ID)
	if got := executor.RefreshCalls(); got != 1 {
		t.Fatalf("unchanged uncertain grant refreshed %d times", got)
	}
	if m.shouldRefresh(current, now.Add(time.Hour)) {
		t.Fatal("automatic refresh repeats uncertain grant before expiry")
	}
}
