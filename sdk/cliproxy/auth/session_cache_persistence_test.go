package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
)

func durableTestSelector(t *testing.T, path string) *SessionAffinitySelector {
	t.Helper()
	s := NewSessionAffinitySelectorWithConfig(SessionAffinityConfig{StateFile: path, TTL: time.Hour})
	t.Cleanup(s.Stop)
	return s
}

func durableTestPick(t *testing.T, s *SessionAffinitySelector, id string, auths ...*Auth) *Auth {
	t.Helper()
	a, err := s.Pick(context.Background(), "codex", "test-model", cliproxyexecutor.Options{Headers: http.Header{"Session-Id": {id}}}, auths)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestDurableAffinityRestartFailoverAndRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	a := &Auth{ID: "account-a", Provider: "codex", Status: StatusActive}
	b := &Auth{ID: "account-b", Provider: "codex", Status: StatusActive}
	s := durableTestSelector(t, path)
	if got := durableTestPick(t, s, "conversation", b); got.ID != b.ID {
		t.Fatal(got.ID)
	}
	s.Stop()
	s = durableTestSelector(t, path)
	if got := durableTestPick(t, s, "conversation", a, b); got.ID != b.ID {
		t.Fatalf("restart moved binding to %s", got.ID)
	}
	// An unavailable saved credential must still fail over and persist the replacement.
	if got := durableTestPick(t, s, "conversation", a); got.ID != a.ID {
		t.Fatal(got.ID)
	}
	s.Stop()
	s = durableTestSelector(t, path)
	if got := durableTestPick(t, s, "conversation", a, b); got.ID != a.ID {
		t.Fatal(got.ID)
	}
}

func TestDurableAffinityAliasesExpiryAndInvalidation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := durableTestSelector(t, path)
	s.cache.SetAliases("a", "parent", "alias")
	s.cache.Set("expired", "a")
	s.cache.mu.Lock()
	old := s.cache.entries["expired"]
	s.cache.replaceAliasGroupsLocked("a", time.Now().Add(-time.Hour), []string{"expired"}, old)
	s.cache.persistLocked()
	s.cache.mu.Unlock()
	s.Stop()
	s = durableTestSelector(t, path)
	if got, ok := s.cache.Get("alias"); !ok || got != "a" {
		t.Fatal("lost alias")
	}
	if _, ok := s.cache.Get("expired"); ok {
		t.Fatal("revived expired binding")
	}
	s.cache.CompareAndDelete("parent", "a")
	s.Stop()
	s = durableTestSelector(t, path)
	if _, ok := s.cache.Get("parent"); ok {
		t.Fatal("revived invalidated alias")
	}
	if got, ok := s.cache.Get("alias"); !ok || got != "a" {
		t.Fatal("deleted surviving alias")
	}
	s.InvalidateAuth("a")
	s.Stop()
	s = durableTestSelector(t, path)
	if s.cache.Len() != 0 {
		t.Fatal("revived removed account")
	}
}

func TestDurableAffinityConcurrentColdRequestsAndReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := durableTestSelector(t, path)
	other := durableTestSelector(t, path)
	if s.cache != other.cache {
		t.Fatal("reload has independent writer")
	}
	a := &Auth{ID: "a", Provider: "codex", Status: StatusActive}
	b := &Auth{ID: "b", Provider: "codex", Status: StatusActive}
	start := make(chan struct{})
	results := make(chan string, 64)
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			got, err := s.Pick(context.Background(), "codex", "test-model", cliproxyexecutor.Options{Headers: http.Header{"Session-Id": {"cold"}}}, []*Auth{a, b})
			if err != nil {
				results <- "ERROR"
				return
			}
			results <- got.ID
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	var selected string
	for got := range results {
		if selected == "" {
			selected = got
		}
		if got == "ERROR" || got != selected {
			t.Fatal("split cold binding", selected, got)
		}
	}
	s.Stop()
	if got := durableTestPick(t, other, "cold", a, b); got.ID != selected {
		t.Fatal("reload lost binding")
	}
}

func TestDurableAffinityCorruptionAndWriteFailureFailClosed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	s := durableTestSelector(t, path)
	a := &Auth{ID: "a", Provider: "codex", Status: StatusActive}
	opts := cliproxyexecutor.Options{Headers: http.Header{"Session-Id": {"session"}}}
	if got, err := s.Pick(context.Background(), "codex", "model", opts, []*Auth{a}); err == nil || got != nil {
		t.Fatal("corrupt file accepted")
	}
	data, _ := os.ReadFile(path)
	if string(data) != "broken" {
		t.Fatal("overwrote corrupt state")
	}
	path = filepath.Join(t.TempDir(), "state.json")
	s = durableTestSelector(t, path)
	// Replacing a file with a directory simulates a write/replace failure on all hosts.
	if err := os.Rename(path, path+".saved"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	if got, err := s.Pick(context.Background(), "codex", "model", opts, []*Auth{a}); err == nil || got != nil {
		t.Fatal("unsaved binding routed upstream")
	}
}

func TestDurableAffinityAbruptProcessExitAndExclusiveWriter(t *testing.T) {
	if mode := os.Getenv("SE_AFFINITY_TEST_CHILD"); mode != "" {
		s := NewSessionAffinitySelectorWithConfig(SessionAffinityConfig{StateFile: os.Getenv("SE_AFFINITY_TEST_PATH"), TTL: time.Hour})
		if mode == "lock" {
			if s.stateErr == nil {
				os.Exit(7)
			}
			os.Exit(0)
		}
		if s.stateErr != nil {
			os.Exit(8)
		}
		s.cache.SetAliases("account-b", "one", "two")
		if s.cache.persistenceError() != nil {
			os.Exit(9)
		}
		os.Exit(0) // No Stop/defer/graceful-shutdown flush.
	}
	path := filepath.Join(t.TempDir(), "state.json")
	run := func(mode string) {
		cmd := exec.Command(os.Args[0], "-test.run=^TestDurableAffinityAbruptProcessExitAndExclusiveWriter$")
		cmd.Env = append(os.Environ(), "SE_AFFINITY_TEST_CHILD="+mode, "SE_AFFINITY_TEST_PATH="+path)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("child: %v %s", err, out)
		}
	}
	run("write")
	s := durableTestSelector(t, path)
	if got, ok := s.cache.Get("two"); !ok || got != "account-b" {
		t.Fatal("abrupt exit lost binding")
	}
	run("lock")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var state sessionCheckpoint
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatal(err)
	}
	if len(state.Groups) != 1 {
		t.Fatal("unexpected checkpoint groups")
	}
}
