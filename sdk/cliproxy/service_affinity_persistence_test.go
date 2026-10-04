package cliproxy

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	internalconfig "github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
)

func TestRoutingSelectorPersistentAffinityEnvironment(t *testing.T) {
	t.Setenv("CLIPROXY_AFFINITY_STATE_FILE", filepath.Join(t.TempDir(), "routing.json"))
	state := normalizedRoutingRuntimeState(&internalconfig.Config{Routing: internalconfig.RoutingConfig{SessionAffinity: true, SessionAffinityTTL: "24h"}})
	makeSelector := func() *coreauth.SessionAffinitySelector {
		s := newRoutingSelector(state).(*coreauth.SessionAffinitySelector)
		t.Cleanup(s.Stop)
		return s
	}
	s := makeSelector()
	a := &coreauth.Auth{ID: "a", Provider: "codex", Status: coreauth.StatusActive}
	b := &coreauth.Auth{ID: "b", Provider: "codex", Status: coreauth.StatusActive}
	opts := cliproxyexecutor.Options{Headers: http.Header{"Session-Id": {"integration-session"}}}
	if got, err := s.Pick(context.Background(), "codex", "model", opts, []*coreauth.Auth{b}); err != nil || got.ID != "b" {
		t.Fatalf("initial selection: %v %v", got, err)
	}
	s.Stop()
	s = makeSelector()
	if got, err := s.Pick(context.Background(), "codex", "model", opts, []*coreauth.Auth{a, b}); err != nil || got.ID != "b" {
		t.Fatalf("restored selection: %v %v", got, err)
	}
}
