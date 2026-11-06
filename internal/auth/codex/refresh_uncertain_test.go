package codex

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
)

type refreshLostBody struct{}

func (*refreshLostBody) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (*refreshLostBody) Close() error             { return nil }

func TestRefreshDoesNotReplayUncertainExchange(t *testing.T) {
	for _, stage := range []string{"transport", "body", "json", "missing_access"} {
		t.Run(stage, func(t *testing.T) {
			resetCodexRefreshGroupForTest()
			t.Cleanup(resetCodexRefreshGroupForTest)
			var calls atomic.Int32
			auth := &CodexAuth{httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				calls.Add(1)
				if stage == "transport" {
					return nil, errors.New("lost response with fake-credential-secret")
				}
				var body io.ReadCloser = io.NopCloser(strings.NewReader(`{"invalid":`))
				if stage == "body" {
					body = &refreshLostBody{}
				}
				if stage == "missing_access" {
					body = io.NopCloser(strings.NewReader(`{}`))
				}
				return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: body, Request: req}, nil
			})}}
			_, err := auth.RefreshTokensWithRetry(context.Background(), "synthetic-uncertain-refresh", 3)
			if err == nil || !strings.Contains(err.Error(), "refresh_outcome_unknown") {
				t.Fatalf("wrong error: %v", err)
			}
			if strings.Contains(err.Error(), "fake-credential-secret") {
				t.Fatal("transport details leaked")
			}
			if calls.Load() != 1 {
				t.Fatalf("uncertain refresh made %d calls, want one", calls.Load())
			}
		})
	}
}
