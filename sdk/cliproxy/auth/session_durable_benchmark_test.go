package auth

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v8/sdk/cliproxy/executor"
	log "github.com/sirupsen/logrus"
)

// No network or provider credentials: measure selection and its real filesystem
// checkpoint on the host where the benchmark is executed. Keep session count
// fixed so growth during the timed loop does not distort comparisons.
func BenchmarkDurableSelection(b *testing.B) {
	level := log.GetLevel()
	log.SetLevel(log.WarnLevel)
	defer log.SetLevel(level)
	for _, size := range []int{2, 32, 256} {
		for _, durable := range []bool{false, true} {
			for _, parallel := range []bool{false, true} {
				b.Run(fmt.Sprintf("sessions=%d/durable=%t/parallel=%t", size, durable, parallel), func(b *testing.B) {
					cfg := SessionAffinityConfig{TTL: 24 * time.Hour}
					if durable {
						cfg.StateFile = filepath.Join(b.TempDir(), "affinity.json")
					}
					s := NewSessionAffinitySelectorWithConfig(cfg)
					defer s.Stop()
					auths := []*Auth{{ID: "a", Provider: "codex", Status: StatusActive}, {ID: "b", Provider: "codex", Status: StatusActive}}
					opts := func(i int) cliproxyexecutor.Options {
						return cliproxyexecutor.Options{Headers: http.Header{"Session-Id": {fmt.Sprintf("bench-%d", i)}}}
					}
					expected := make([]string, size)
					for i := range expected {
						a, err := s.Pick(context.Background(), "codex", "model", opts(i), auths)
						if err != nil || a == nil {
							b.Fatalf("seed: %v", err)
							return
						}
						expected[i] = a.ID
					}
					var next atomic.Uint64
					var failures atomic.Uint64
					pick := func() {
						i := int(next.Add(1)-1) % size
						a, err := s.Pick(context.Background(), "codex", "model", opts(i), auths)
						if err != nil || a == nil || a.ID != expected[i] {
							failures.Add(1)
						}
					}
					b.ReportAllocs()
					b.ResetTimer()
					if parallel {
						b.RunParallel(func(pb *testing.PB) {
							for pb.Next() {
								pick()
							}
						})
					} else {
						for i := 0; i < b.N; i++ {
							pick()
						}
					}
					b.StopTimer()
					if n := failures.Load(); n != 0 {
						b.Fatalf("%d routing errors or account changes", n)
					}
				})
			}
		}
	}
}
