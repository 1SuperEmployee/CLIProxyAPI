package executor

import (
	"sync"
	"testing"

	claudeauth "github.com/router-for-me/CLIProxyAPI/v8/internal/auth/claude"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/auth"
)

func TestClaudeSetupTokenMetadataFlags(t *testing.T) {
	for _, key := range []string{"skip_account_profile", "is_setup_token", "setup_token"} {
		for _, value := range []any{nil, false, "true", 1, true} {
			auth := &cliproxyauth.Auth{Metadata: map[string]any{key: value}}
			want, _ := value.(bool)
			if got := isClaudeSetupToken(auth, "sk-ant-oat-synthetic"); got != want {
				t.Errorf("flag %s value %v: got %t, want %t", key, value, got, want)
			}
		}
	}
	for _, auth := range []*cliproxyauth.Auth{nil, {}, {Metadata: map[string]any{}}} {
		if isClaudeSetupToken(auth, "sk-ant-oat-synthetic") {
			t.Error("missing metadata classified as setup token")
		}
	}
}

func TestClaudeSetupTokenMetadataConcurrentProfileWrites(t *testing.T) {
	auth := &cliproxyauth.Auth{Metadata: map[string]any{
		"skip_account_profile": false, "is_setup_token": false, "setup_token": false,
	}}
	start := make(chan struct{})
	var workers sync.WaitGroup
	for worker := 0; worker < 16; worker++ {
		workers.Add(1)
		go func(writer bool) {
			defer workers.Done()
			<-start
			for i := 0; i < 100; i++ {
				if writer {
					claudeauth.StoreMetadataString(&auth.Metadata, "account_uuid", "synthetic-account")
					continue
				}
				// All three flags are false, so every reader traverses all three accesses.
				if isClaudeSetupToken(auth, "sk-ant-oat-synthetic") {
					t.Error("profile writes changed setup-token classification")
				}
			}
		}(worker%2 == 0)
	}
	close(start)
	workers.Wait()
}
