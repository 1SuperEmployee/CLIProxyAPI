package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
	"gopkg.in/yaml.v3"
)

func TestReloadConfigIfChanged_WriteDuringReloadIsNotSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	authDir := filepath.Join(tmpDir, "auth")
	if err := os.MkdirAll(authDir, 0o755); err != nil {
		t.Fatalf("failed to create auth dir: %v", err)
	}

	configPath := filepath.Join(tmpDir, "config.yaml")
	writeConfig := func(port int) {
		data, err := yaml.Marshal(&config.Config{
			Port:               port,
			AuthDir:            authDir,
			CredentialInFlight: config.DefaultCredentialInFlightConfig(),
		})
		if err != nil {
			t.Fatalf("failed to marshal config: %v", err)
		}
		if err = os.WriteFile(configPath, data, 0o644); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}
	}

	writeConfig(8080)

	// A management save that lands while the first reload is applying its config
	// must still be picked up by the next trigger instead of matching the hash.
	var appliedPorts []int
	w := &Watcher{
		configPath: configPath,
		authDir:    authDir,
	}
	w.reloadCallback = func(cfg *config.Config) {
		appliedPorts = append(appliedPorts, cfg.Port)
		if len(appliedPorts) == 1 {
			writeConfig(9090)
		}
	}

	w.reloadConfigIfChanged()
	w.reloadConfigIfChanged()
	if len(appliedPorts) != 2 || appliedPorts[1] != 9090 {
		t.Fatalf("expected the write made during reload to be applied next, applied ports %v", appliedPorts)
	}

	w.reloadConfigIfChanged()
	if len(appliedPorts) != 2 {
		t.Fatalf("expected unchanged config to be skipped after catching up, applied ports %v", appliedPorts)
	}
}

func TestReloadConfigIfChanged_LoadsTheBytesItHashed(t *testing.T) {
	tmpDir := t.TempDir()
	authDir := filepath.Join(tmpDir, "auth")
	if err := os.MkdirAll(authDir, 0o755); err != nil {
		t.Fatalf("failed to create auth dir: %v", err)
	}
	configPath := filepath.Join(tmpDir, "config.yaml")
	content := []byte("port: 8317\nremote-management:\n  secret-key: \"$2a$10$abcdefghijklmnopqrstuuabcdefghijklmnopqrstuvwxyz01234\"\n")
	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// A management save truncates the file right after the watcher reads it, as
	// os.Create/os.WriteFile do. Loading must use the bytes that were read, or the
	// reload applies an empty config and drops the management secret.
	original := readConfigFile
	t.Cleanup(func() { readConfigFile = original })
	readConfigFile = func(path string) ([]byte, error) {
		data, err := original(path)
		if err == nil {
			if errTruncate := os.WriteFile(path, nil, 0o600); errTruncate != nil {
				t.Fatalf("failed to truncate config: %v", errTruncate)
			}
		}
		return data, err
	}

	var applied *config.Config
	w := &Watcher{
		configPath:     configPath,
		authDir:        authDir,
		reloadCallback: func(cfg *config.Config) { applied = cfg },
	}
	w.reloadConfigIfChanged()

	if applied == nil {
		t.Fatal("expected the reload to apply the config that was read")
	}
	if applied.Port != 8317 || applied.RemoteManagement.SecretKey == "" {
		t.Fatalf("applied config port=%d secret empty=%v, want the read config with its secret", applied.Port, applied.RemoteManagement.SecretKey == "")
	}
}

func TestAddOrUpdateClientStampsSnapshotBeforeReading(t *testing.T) {
	authDir := t.TempDir()
	path := filepath.Join(authDir, "codex.json")
	if err := os.WriteFile(path, []byte(`{"type":"codex","email":"user@example.com"}`), 0o600); err != nil {
		t.Fatalf("write auth file: %v", err)
	}
	var readAt time.Time
	original := readAuthFile
	t.Cleanup(func() { readAuthFile = original })
	readAuthFile = func(name string) ([]byte, error) {
		data, err := original(name)
		readAt = time.Now()
		return data, err
	}

	w := &Watcher{authDir: authDir, lastAuthHashes: make(map[string]string)}
	w.SetConfig(&config.Config{AuthDir: authDir})
	w.addOrUpdateClient(path)

	if len(w.currentAuths) == 0 {
		t.Fatal("expected the auth file to be synthesized")
	}
	for id, auth := range w.currentAuths {
		if auth.UpdatedAt.After(readAt) {
			t.Fatalf("auth %s snapshot stamped %v after the file read finished", id, auth.UpdatedAt.Sub(readAt))
		}
	}
}
