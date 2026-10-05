package auth

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDurableAffinityInvalidRestoreCannotRewriteCheckpoint(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	original := []byte(`{"version":1,"groups":[{"auth_id":"a","aliases":["one"],"expires_at":"2099-01-01T00:00:00Z"},{"auth_id":"","aliases":["two"],"expires_at":"2099-01-01T00:00:00Z"}]}`)
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatal(err)
	}
	c := NewSessionCache(time.Hour)
	c.Stop()
	c.mu.Lock()
	c.statePath = path
	err := c.restoreLocked()
	c.mu.Unlock()
	if err == nil {
		t.Fatal("invalid checkpoint accepted")
	}
	// Deterministically simulate a cleanup tick already queued during failed
	// initialization. Closing stopCh alone does not drain an in-flight tick.
	c.cleanup()
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, original) {
		t.Fatalf("failed restore allowed checkpoint replacement: %v", err)
	}
	if c.persistenceError() == nil {
		t.Fatal("restore error was not latched")
	}
}

func TestDurableAffinityNoMutationDoesNotWriteCheckpoint(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := durableTestSelector(t, path)
	s.cache.Set("existing", "a")
	if err := os.Rename(path, path+".saved"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	// Failed replacements expose unnecessary writes without timestamp sleeps
	// or OS-specific filesystem watchers. These operations change no binding.
	s.cache.GetAndRefresh("missing")
	s.cache.Touch("existing", "wrong-account")
	s.cache.CompareAndDelete("existing", "wrong-account")
	s.cache.Invalidate("missing")
	s.cache.InvalidateAuth("missing-account")
	s.cache.cleanup()
	if err := s.cache.persistenceError(); err != nil {
		t.Fatalf("no-op attempted a checkpoint write: %v", err)
	}
	// A real binding change must still attempt persistence and report failure.
	s.cache.Set("new", "b")
	if s.cache.persistenceError() == nil {
		t.Fatal("binding mutation ignored failed checkpoint write")
	}
}
