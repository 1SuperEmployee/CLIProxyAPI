package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDurableAffinityWriteFailureKeepsRouting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := durableTestSelector(t, path)
	a := &Auth{ID: "a", Provider: "codex", Status: StatusActive}
	b := &Auth{ID: "b", Provider: "codex", Status: StatusActive}
	durableTestPick(t, s, "warm", b)
	if err := os.Rename(path, path+".saved"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		if got := durableTestPick(t, s, "warm", a, b); got.ID != b.ID {
			t.Fatal("disk failure moved healthy conversation")
		}
		durableTestPick(t, s, "new", a, b)
	}
	if s.cache.persistenceError() == nil {
		t.Fatal("write failure not reported by persistence diagnostics")
	}
	if err := os.Rename(path, path+".blocked"); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(path+".saved", path); err != nil {
		t.Fatal(err)
	}
	// Advance retry eligibility without a timing-dependent sleep.
	s.cache.mu.Lock()
	s.cache.persistRetryAt = time.Now().Add(-time.Second)
	s.cache.mu.Unlock()
	durableTestPick(t, s, "warm", a, b)
	if err := s.cache.persistenceError(); err != nil {
		t.Fatal("checkpoint did not recover", err)
	}
	s.Stop()
	s = durableTestSelector(t, path)
	if got := durableTestPick(t, s, "warm", a, b); got.ID != b.ID {
		t.Fatal("recovered checkpoint lost binding")
	}
}
