package executor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	codexresponses "github.com/router-for-me/CLIProxyAPI/v8/internal/translator/codex/openai/responses"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
	sdktranslator "github.com/router-for-me/CLIProxyAPI/v8/sdk/translator"
	"github.com/tidwall/gjson"
)

// Local fixtures only. A clean transport EOF is not a complete model response.
func TestClaudeStreamCompletionRegression(t *testing.T) {
	for _, complete := range []bool{true, false} {
		name := "truncated"
		if complete {
			name = "complete"
		}
		t.Run(name, func(t *testing.T) {
			stream := "event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_audit\",\"type\":\"message\",\"role\":\"assistant\",\"model\":\"claude-sonnet-5-5\",\"content\":[],\"usage\":{\"input_tokens\":5,\"output_tokens\":0}}}\n\nevent: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"partial\"}}\n\n"
			if complete {
				stream += "event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				_, _ = w.Write([]byte(stream))
			}))
			defer server.Close()
			e := NewClaudeExecutor(&config.Config{})
			auth := &cliproxyauth.Auth{Attributes: map[string]string{"api_key": "synthetic-audit-key", "base_url": server.URL}}
			r, err := e.ExecuteStream(context.Background(), auth, cliproxyexecutor.Request{Model: "claude-sonnet-5-5", Payload: []byte(`{"messages":[{"role":"user","content":"test"}]}`)}, cliproxyexecutor.Options{SourceFormat: sdktranslator.FormatClaude})
			if err != nil {
				t.Fatalf("fixture startup: %v", err)
			}
			gotErr := false
			for chunk := range r.Chunks {
				if chunk.Err != nil {
					gotErr = true
				}
			}
			if complete && gotErr {
				t.Fatal("complete fixture failed")
			}
			if !complete && !gotErr {
				t.Fatal("truncated stream closed with no terminal error")
			}
		})
	}
}

func TestClaudeOnDemandCompactionRegression(t *testing.T) {
	payload := []byte(`{"model":"claude-opus-5","max_tokens":4096,"thinking":{"type":"adaptive"},"messages":[{"role":"user","content":"test"}],"compaction":{"type":"summarize"}}`)
	for _, stream := range []bool{false, true} {
		body := executeClaudeContextManagementRequest(t, &config.Config{}, payload, stream)
		if gjson.GetBytes(body, "context_management").Exists() {
			t.Errorf("stream=%v: proxy injected context_management alongside on-demand compaction", stream)
		}
	}
}

func TestCodexAutomaticCompactionRuleRegression(t *testing.T) {
	payload := []byte(`{"model":"gpt-6-luna","context_management":[{"type":"compaction","compact_threshold":12000}],"input":[{"role":"user","content":"test"}]}`)
	body, err := codexresponses.ConvertOpenAIResponsesRequestToCodex("gpt-6-luna", payload, true)
	if err != nil {
		t.Fatal(err)
	}
	if !gjson.GetBytes(body, "context_management").Exists() {
		t.Fatal("proxy stripped the requested server-side automatic compaction rule")
	}
}

func TestClaudeNonstreamSSERequiresMessageStop(t *testing.T) {
	start := `data: {"type":"message_start","message":{"id":"msg_test","model":"claude-sonnet-5-5"}}` + "\n\n"
	delta := `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"}}` + "\n\n"
	if err := validateClaudeStreamingResponse([]byte(start + delta)); err == nil {
		t.Fatal("accepted SSE body without message_stop")
	}
	if err := validateClaudeStreamingResponse([]byte(start + delta + "data: {\"type\":\"message_stop\"}\n\n")); err != nil {
		t.Fatal(err)
	}
}
