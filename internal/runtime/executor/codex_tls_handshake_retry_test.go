package executor

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/registry"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
	sdktranslator "github.com/router-for-me/CLIProxyAPI/v8/sdk/translator"
)

// Exercise the real Codex HTTP streaming executor, not the WebSocket dialer.
// The synthetic CONNECT proxy never forwards traffic to an external provider.
func TestCodexHTTPStreamManagerRetriesSilentTLSHandshakeThroughConnectProxy(t *testing.T) {
	proxy, proxyURL := startSilentConnectProxy(t)

	cfg := &config.Config{}
	cfg.ProxyURL = proxyURL
	cfg.TLSHandshakeTimeoutSeconds = 1

	manager := cliproxyauth.NewManager(nil, nil, nil)
	manager.SetRetryConfig(3, 0, 0)
	manager.RegisterExecutor(NewCodexExecutor(cfg))

	authID := uuid.NewString() + "-silent-tls-handshake"
	model := "gpt-5.4"
	auth := &cliproxyauth.Auth{
		ID:       authID,
		Provider: "codex",
		Attributes: map[string]string{
			"api_key":  "synthetic-silent-handshake-test",
			"base_url": "https://chatgpt.com/backend-api/codex",
		},
	}
	reg := registry.GetGlobalRegistry()
	reg.RegisterClient(authID, "codex", []*registry.ModelInfo{{ID: model}})
	t.Cleanup(func() { reg.UnregisterClient(authID) })
	if _, errRegister := manager.Register(context.Background(), auth); errRegister != nil {
		t.Fatalf("register auth: %v", errRegister)
	}

	payload := []byte(`{"model":"gpt-5.4","input":[{"role":"user","content":"hello"}],"stream":true}`)
	executeDone := make(chan error, 1)
	go func() {
		_, errExecute := manager.ExecuteStream(context.Background(), []string{"codex"}, cliproxyexecutor.Request{
			Model:   model,
			Payload: payload,
		}, cliproxyexecutor.Options{SourceFormat: sdktranslator.FromString("openai-response")})
		executeDone <- errExecute
	}()

	var errExecute error
	select {
	case errExecute = <-executeDone:
	case <-time.After(15 * time.Second):
		t.Fatal("Manager.Execute did not finish after four bounded TLS handshakes")
	}
	if errExecute == nil {
		t.Fatal("Execute() error = nil, want final configured TLS handshake timeout")
	}
	if !cliproxyexecutor.IsConfiguredTLSHandshakeTimeoutError(errExecute) {
		t.Fatalf("Execute() error = %v, want typed configured TLS handshake timeout", errExecute)
	}
	if attempts := proxy.attempts.Load(); attempts != 4 {
		t.Fatalf("silent proxy handshake attempts = %d, want 4", attempts)
	}
	for attempt := 1; attempt <= 4; attempt++ {
		select {
		case <-proxy.closed:
		case <-time.After(5 * time.Second):
			t.Fatalf("silent proxy did not observe failed socket %d being closed", attempt)
		}
	}

	updatedAuth, ok := manager.GetByID(authID)
	if !ok || updatedAuth == nil {
		t.Fatal("registered auth not found after handshake retries")
	}
	if updatedAuth.Unavailable || !updatedAuth.NextRetryAfter.IsZero() {
		t.Fatalf("configured handshake timeout cooled the credential: unavailable=%t next_retry_after=%v", updatedAuth.Unavailable, updatedAuth.NextRetryAfter)
	}
	if state := updatedAuth.ModelStates[model]; state != nil && (state.Unavailable || !state.NextRetryAfter.IsZero()) {
		t.Fatalf("configured handshake timeout cooled the model: %#v", state)
	}
}
