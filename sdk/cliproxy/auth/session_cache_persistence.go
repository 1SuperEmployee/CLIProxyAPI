package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// A host-local routing checkpoint. No tokens, request bodies, or message history.
type sessionCheckpoint struct {
	Version int                      `json:"version"`
	Groups  []sessionCheckpointGroup `json:"groups"`
}

type sessionCheckpointGroup struct {
	AuthID    string    `json:"auth_id"`
	Aliases   []string  `json:"aliases"`
	ExpiresAt time.Time `json:"expires_at"`
}

type sharedSessionCache struct {
	cache *SessionCache
	lock  *os.File
	refs  int
}

var persistentSessionCaches = struct {
	sync.Mutex
	byPath map[string]*sharedSessionCache
}{byPath: make(map[string]*sharedSessionCache)}

// Share one cache across hot-reloaded selectors. The OS lock excludes other processes.
func acquirePersistentSessionCache(path string, ttl time.Duration) (*SessionCache, func(), error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, err
	}
	persistentSessionCaches.Lock()
	defer persistentSessionCaches.Unlock()
	shared := persistentSessionCaches.byPath[path]
	if shared == nil {
		if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return nil, nil, err
		}
		lock, errLock := lockSessionState(path + ".lock")
		if errLock != nil {
			return nil, nil, fmt.Errorf("lock affinity checkpoint: %w", errLock)
		}
		cache := NewSessionCache(ttl)
		cache.mu.Lock()
		cache.statePath = path
		err = cache.restoreLocked()
		if err == nil {
			cache.persistLocked()
			err = cache.persistErr
		}
		cache.mu.Unlock()
		if err != nil {
			cache.Stop()
			_ = lock.Close()
			return nil, nil, err
		}
		shared = &sharedSessionCache{cache: cache, lock: lock}
		persistentSessionCaches.byPath[path] = shared
	}
	shared.refs++
	shared.cache.mu.Lock()
	shared.cache.ttl = ttl
	shared.cache.mu.Unlock()
	var once sync.Once
	release := func() {
		once.Do(func() {
			persistentSessionCaches.Lock()
			defer persistentSessionCaches.Unlock()
			shared.refs--
			if shared.refs == 0 {
				shared.cache.Stop()
				// Serialize with an in-flight cache mutation before releasing the OS lock.
				shared.cache.mu.Lock()
				shared.cache.persistErr = errors.New("affinity checkpoint closed")
				_ = shared.lock.Close()
				shared.cache.mu.Unlock()
				delete(persistentSessionCaches.byPath, path)
			}
		})
	}
	return shared.cache, release, nil
}

func (c *SessionCache) persistenceError() error {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.persistErr != nil {
		return fmt.Errorf("session affinity persistence unavailable: %w", c.persistErr)
	}
	return nil
}

func (c *SessionCache) restoreLocked() error {
	f, err := os.Open(c.statePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	var state sessionCheckpoint
	d := json.NewDecoder(io.LimitReader(f, 64<<20))
	d.DisallowUnknownFields()
	if err = d.Decode(&state); err != nil {
		return fmt.Errorf("decode affinity checkpoint: %w", err)
	}
	var extra any
	if err = d.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("trailing data in affinity checkpoint")
	}
	if state.Version != 1 {
		return errors.New("unsupported affinity checkpoint version")
	}
	seen := make(map[string]bool)
	now := time.Now()
	for _, group := range state.Groups {
		if group.AuthID == "" || len(group.Aliases) == 0 || len(group.Aliases) > maxStableSessionAliases+1 || group.ExpiresAt.IsZero() {
			return errors.New("invalid affinity checkpoint group")
		}
		for _, alias := range group.Aliases {
			if alias == "" || seen[alias] {
				return errors.New("invalid or duplicate affinity checkpoint alias")
			}
			seen[alias] = true
		}
		if len(seen) > c.maxEntries {
			return errors.New("affinity checkpoint exceeds cache capacity")
		}
		if now.Before(group.ExpiresAt) {
			c.replaceAliasGroupsLocked(group.AuthID, group.ExpiresAt, group.Aliases)
		}
	}
	return nil
}

// Called under the cache mutex after complete mutations, before Pick can return.
// Any failure latches the cache closed for routing until an operator repairs it.
func (c *SessionCache) persistLocked() {
	if c.statePath == "" || c.persistErr != nil {
		return
	}
	state := sessionCheckpoint{Version: 1}
	now := time.Now()
	for el := c.evictionOrder.Front(); el != nil; el = el.Next() {
		group := c.groups[el.Value.(string)]
		if now.Before(group.expiresAt) {
			state.Groups = append(state.Groups, sessionCheckpointGroup{group.authID, group.aliases, group.expiresAt})
		}
	}
	data, err := json.Marshal(state)
	if err == nil {
		err = writeSessionCheckpoint(c.statePath, data)
	}
	if err != nil {
		c.persistErr = err
	}
}

func writeSessionCheckpoint(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".affinity-*.tmp")
	if err != nil {
		return err
	}
	name := f.Name()
	defer func() { _ = f.Close(); _ = os.Remove(name) }()
	if _, err = f.Write(data); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return replaceSessionCheckpoint(name, path)
}
