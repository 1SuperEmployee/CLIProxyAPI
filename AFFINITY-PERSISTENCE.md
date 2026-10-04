# Host-local durable affinity

Small opt-in patch on upstream v8.0.13. Set the service environment variable
`CLIPROXY_AFFINITY_STATE_FILE` to an absolute path outside the auth directory.
Existing `routing.session-affinity` must be enabled. Its configured TTL still
applies, including time while the gateway is stopped. An empty environment
variable leaves upstream in-memory behavior intact.

The versioned JSON checkpoint contains account IDs, bounded session aliases and
absolute expiry times. It contains no OAuth tokens or request/response bodies.
Protect the parent directory using the service account's filesystem permissions.
For our Windows service this directory inherits the protected runtime ACL.

Each completed cache mutation is synchronously checkpointed through a temporary
file, file sync and replacement. Windows uses MoveFileEx with replacement and
write-through flags; Unix also syncs the parent directory. Replacement must stay
on the same filesystem. A process exit does not require a graceful shutdown flush.
This does not claim survival of every filesystem, disk or physical power failure.

An OS lock prevents a second process from owning the same checkpoint. Selectors
created by configuration reload share a reference-counted cache within the process.
Corrupt/unsupported state and write failures fail subsequent selection rather
than silently changing accounts. Repair the disk/path or restore a valid checkpoint
with the service stopped, then restart. Do not delete the checkpoint as routine
recovery; doing so deliberately discards bindings.

The existing credential eligibility, priority, model namespaces, alias handling,
expiry, invalidation and failover rules remain in effect. Concurrent cold selection
for the same cache key is serialized across lookup and binding using fixed lock
stripes; unrelated keys may briefly share a stripe, but network inference occurs
after selection and is not held under those locks.

Scope: explicit/stable harness session bindings in SessionCache. The optional
content-derived Merkle/LCP matcher remains in memory. Our Claude Code and Codex
workers supply explicit identities. This is not provider-side prompt-cache
persistence, multi-host replication, or a reset-aware account scheduler.

Snapshot writes scale with the number of live bindings. This is intended for the
small personal fleet; benchmark before adopting it for a large shared gateway.

The second patch saves only actual mutations and refreshes a healthy binding once
during selection. It still checkpoints synchronously before returning the account.
A failed restore latches persistence closed before a queued cleanup can run, so
partially loaded state cannot overwrite the original damaged checkpoint.

Windows amd64 measurements on 4 October 2026, four Go processors, three runs of
100 selections per case. Values are medians of average serial selection time;
they exclude model inference and completion-time TTL refresh:

| Saved sessions | First patch | Second patch |
| --- | --- | --- |
| 2 | 6.267 ms | 3.172 ms |
| 32 | 4.148 ms | 2.023 ms |
| 256 | 5.916 ms | 2.296 ms |

Concurrent selection benchmarks also passed without account drift. These small
samples describe this Windows host, not a latency SLA or maximum throughput.
Benchmark with `go test ./sdk/cliproxy/auth -run '^$' -bench
'^BenchmarkDurableSelection$' -benchtime=100x -count=3 -cpu=4`. Temporary synthetic
state is used; no network inference or provider credentials are needed.

Validation commands:

```sh
go test -race ./sdk/cliproxy/auth -run 'TestDurableAffinity|TestSessionCache|TestSessionAffinity|TestManagerSessionAffinity'
go test ./sdk/cliproxy/... ./internal/config/...
go build -o ../builds/cli-proxy-affinity-durable-mac ./cmd/server
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go test -c ./sdk/cliproxy/auth -o ../builds/affinity-auth-test.exe
```

Run the Windows test executable on the gateway PC against temporary synthetic
state. Then verify two real completed-turn conversations remain on a non-default
account across a gateway service restart. Keep the previous executable and wrapper XML for rollback. The first and second
custom releases share checkpoint version 1, so rollback between them retains
durable bindings. Returning to the official executable discards this custom
affinity persistence; neither rollback requires OAuth changes.
