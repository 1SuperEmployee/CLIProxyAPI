package cliproxy

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/runtime/executor"
	coreauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
)

func TestRoutingQuotaOptInAndDefaultPreservation(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		cfg := &config.Config{Routing: config.RoutingConfig{SessionAffinity: true, NearestReset: enabled}}
		state := normalizedRoutingRuntimeState(cfg)
		sel := newRoutingSelector(state)
		affinity := sel.(*coreauth.SessionAffinitySelector)
		defer affinity.Stop()
		m := coreauth.NewManager(nil, sel, nil)
		// Real CodexAutoExecutor plus passive observations proves the normal
		// service config path, without making a provider request or using OAuth.
		m.RegisterExecutor(executor.NewCodexAutoExecutor(cfg))
		now := time.Now().Truncate(time.Second)
		for i, reset := range []time.Duration{6 * time.Hour, 3 * time.Hour} {
			_, err := m.Register(context.Background(), &coreauth.Auth{ID: fmt.Sprintf("config-quota-%d", i), Provider: "codex", Attributes: map[string]string{coreauth.AttributeAuthKind: coreauth.AuthKindOAuth}, Quota: coreauth.QuotaState{ObservedAt: now, Signals: map[string]string{"X-Codex-Primary-Used-Percent": "20", "X-Codex-Primary-Window-Minutes": "10080", "X-Codex-Primary-Reset-At": fmt.Sprint(now.Add(reset).Unix())}}})
			if err != nil {
				t.Fatal(err)
			}
		}
		a, err := m.SelectAuth(t.Context(), "codex", "", cliproxyexecutor.Options{Headers: http.Header{"X-Session-Id": {"config-quota-session"}}})
		want := "config-quota-0"
		if enabled {
			want = "config-quota-1"
		}
		if err != nil || a == nil || a.ID != want {
			t.Fatalf("enabled=%v: auth=%v error=%v", enabled, a, err)
		}
		withoutAffinity := normalizedRoutingRuntimeState(&config.Config{Routing: config.RoutingConfig{NearestReset: true}})
		if withoutAffinity.nearestReset {
			t.Fatal("quota feature enabled without affinity")
		}
	}
}
