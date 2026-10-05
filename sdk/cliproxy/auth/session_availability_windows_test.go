package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDurableAffinityWindowsStartupReaderKeepsRestoredBindings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	a := &Auth{ID: "a", Provider: "codex", Status: StatusActive}
	b := &Auth{ID: "b", Provider: "codex", Status: StatusActive}
	s := durableTestSelector(t, path)
	durableTestPick(t, s, "warm", b)
	s.Stop()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	s = durableTestSelector(t, path)
	if s.stateErr != nil || s.cache.statePath != path {
		t.Fatal("initial save failure discarded the restored cache", s.stateErr)
	}
	if s.cache.persistenceError() == nil {
		t.Fatal("test did not exercise initial replacement failure")
	}
	if got := durableTestPick(t, s, "warm", a, b); got.ID != b.ID {
		t.Fatal("startup reader lost the saved account")
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	s.cache.mu.Lock()
	s.cache.persistRetryAt = time.Now().Add(-time.Second)
	s.cache.mu.Unlock()
	durableTestPick(t, s, "warm", a, b)
	if err := s.cache.persistenceError(); err != nil {
		t.Fatal("startup save did not recover", err)
	}
	s.Stop()
	s = durableTestSelector(t, path)
	if got := durableTestPick(t, s, "warm", a, b); got.ID != b.ID {
		t.Fatal("restart after recovered startup save lost account")
	}
}

func TestDurableAffinityWindowsReaderDoesNotStopRouting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := durableTestSelector(t, path)
	a := &Auth{ID: "a", Provider: "codex", Status: StatusActive}
	b := &Auth{ID: "b", Provider: "codex", Status: StatusActive}
	durableTestPick(t, s, "warm", b)
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	// Hold the reader across selection deterministically. Windows rejects replacement.
	for i := 0; i < 20; i++ {
		if got := durableTestPick(t, s, "warm", a, b); got.ID != b.ID {
			t.Fatal("reader changed account")
		}
		durableTestPick(t, s, "new", a, b)
	}
	if s.cache.persistenceError() == nil {
		t.Fatal("test did not exercise the Windows replacement failure")
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	s.cache.mu.Lock()
	s.cache.persistRetryAt = time.Now().Add(-time.Second)
	s.cache.mu.Unlock()
	durableTestPick(t, s, "warm", a, b)
	if err := s.cache.persistenceError(); err != nil {
		t.Fatal("reader closed but checkpoint did not recover", err)
	}
}
