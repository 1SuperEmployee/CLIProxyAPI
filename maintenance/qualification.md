# Current update qualification

This is the live qualification record, not a task queue. Completed source reviews
remain under `reviews/`. Read CURRENT.md for the accepted artifact.

## Candidate and official state

The candidate is `se/review-v8.0.23`, source
`5d688e9748e108aa990159b9a8e8762361080bb6`. Official metadata checked on
9 October at 22:21 UTC still identifies v8.0.23 as the latest stable release.
Main matches its commit `d318bcc3afb9ea8782862f5e8cdb541afdc9dfa5`.
All seven stable release notes after the accepted v8.0.16 base were included.
Issue 6491 remains open. Its null-cache-mark defect was independently reproduced;
actual native-client impact remains unknown. See the finished source review.

The last fully accepted rollback remains `8.0.16-se-affinity.4-rc1`. The
`8.0.23-se-affinity.5-rc1` candidate was subsequently deployed at the operator's
request. Actual-client useful work and same-thread follow-ups subsequently passed
on 9 October, within the scope recorded below. Offline checks did not change
the live service. No automatic updater is installed.

## Checks and limits

- Five relevant Mac package suites, focused affinity race tests and the Mac
  server build passed in the source review.
- The broader Mac suite passed 100 packages, with 29 packages having no tests.
  Its only failing package was `internal/discovery`, where the live multicast
  discovery test timed out. An unmodified v8.0.23 checkout reproduced that failure
  with the same toolchain. The discovery source is unchanged by our patch.
  This does not establish the precise environmental cause or count as a full pass.
- The extra Mac protocol race check passed `internal/runtime/executor`.
  `sdk/api/handlers/openai` could not link because the Mac ran out of disk space.
  That check is incomplete, not a test failure or a passing result. Further builds
  moved to a bounded Linux builder.
- The first Linux full suite ran with container networking disabled. It failed
  `TestPionMediaRelayBridgesAudioAndDataChannel` in `internal/client/codex/live`.
  That package is unchanged from v8.0.23. The networking-enabled synthetic recheck
  passed. The networking-enabled full suite passed 101 packages, with 29 packages
  having no tests. A fresh uncached verbose discovery integration run also passed
  on Linux; it was not skipped. No upstream test was disabled or changed.
- Linux Go 1.27.1 focused affinity race checks passed. Protocol race checks passed
  both executor and OpenAI handler packages, covering interruptions, missing
  continuation state and effort directives.

## Windows candidate artifact

Built from the candidate commit above on Linux with Go 1.27.1, MinGW-posix GCC 12
and CGO enabled. Server identity is `8.0.23-se-affinity.5-rc1`, commit `5d688e97`.
Both server and auth test package compiled successfully. Source and synthetic tests
only; no configuration or credentials were supplied to the builder.

```sh
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-posix GOOS=windows GOARCH=amd64 \
  go build -p 2 \
  -ldflags '-s -w -X main.Version=8.0.23-se-affinity.5-rc1 -X main.Commit=5d688e97' \
  -o cli-proxy-affinity-v5-rc1.exe ./cmd/server
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-posix GOOS=windows GOARCH=amd64 \
  go test -p 2 -c ./sdk/cliproxy/auth -o affinity-v5-rc1-tests.exe
```

SHA256:

- Server: `cd72dc10db89d58338daa1c3b6e2bf9774bcb30f9c65257909b228caefad89f7`.
- Auth tests: `157507f354d8efbd6ea00d94a4e8bfd23ef92a850a3daf726ddafb3ef15acb8e`.

These are candidate bytes, deployed and checked with real native useful work and
follow-ups as recorded below. Broader fault-path acceptance remains separate. Both PC
copies passed SHA256 verification. The full Windows
auth package passed all 954 top-level tests, including all ten durability tests.
It used temporary synthetic state, never live auth or checkpoint.
The Windows full auth package requires two tracked public JSON fixtures under
`internal/home/testdata/` and a working directory mirroring `sdk/cliproxy/auth`.
Running the standalone test executable without that layout failed only
`TestConcurrencyDispatchFixture`; the ten durability tests passed. Correcting that
layout produced the full 954-test pass without changing or disabling a test. Output
is captured through Python subprocesses so PowerShell's native-stderr handling
does not treat ordinary test logging as a terminating script error.

On the same Windows host, four Go processors, three runs of 100 selections per
case. Serial durable selection medians, in milliseconds:

| Saved sessions | Accepted v8.0.16 build | v8.0.23 candidate |
| --- | --- | --- |
| 2 | 3.315 | 3.246 |
| 32 | 2.190 | 2.127 |
| 256 | 2.520 | 2.341 |

Both full serial/concurrent benchmark matrices passed without account drift.
No selection regression was detected in this short sample. Active inference was
left running; artifact staging had finished before the compared runs. These are
selector/checkpoint measurements, not proof of lower whole-request latency,
cache-hit rate, provider throughput or interruption correctness.

## Rollback review

Upstream's new credential version and rejected-token bookkeeping do not introduce
a persisted OAuth-file schema change in the reviewed provider token writers or
FileTokenStore path. Credential version initializes in memory on load. This is
source-review evidence, not a live downgrade drill. Preserve the freshest OAuth
files during software rollback. Never restore old credentials with an old binary.
The custom affinity checkpoint remains schema 1; its implementation is unchanged.
The usage subscriber's serialization is unchanged, including latency/first-token
millisecond fields. This supports recorder compatibility; improved upstream
first-token semantics still need observation on actual client traffic.

## Deployment checkpoint

The bounded cutover restored all 26 nonexpired schema-1 bindings. Previously idle
workers were held and released with their native processes unchanged. The service
selected the checksum-verified candidate; the earlier executable and service XML
remain available. No provider enrollment, credential restoration or worker migration
was performed. Private HTTPS access returned the expected unauthenticated 401 from
the Mac and a Linux host. This proves transport/startup, not a completed model turn.

The operator check caught Windows stopping the dependent usage recorder along with
the gateway. It was restored. The private restart operator now restores services
that were running before stopping their gateway dependency. This is an operator
repair; gateway source and candidate bytes remain unchanged. A subsequent bounded
restart verified automatic recorder return and retained all 26 bindings again.
The usage feed reconnected with a fresh heartbeat. Real specialist useful-work
checks subsequently passed as recorded below.

## Native useful-work checkpoint, 9 October

The actual Claude Code and Codex specialists resumed their existing native
conversations after the earlier gateway deployment/restart, completed useful
file/research work, and amended/read back their artifacts in a same-thread
follow-up. Native transcript IDs matched relay associations and usage metadata.
Sonnet 5.5 and Luna/high were retained. All 26 model requests succeeded, with one
account per conversation across model/tool requests and follow-ups. Both follow-ups
had positive provider cache reads. Both conversations ended cleanly without pending
tools or relay messages. Gateway/recorder processes and the recorder connection
remained unchanged. No credentials, permissions, services or routing were changed.

The test conversations were old and idle, without a matching binding at baseline.
This proves stable renewed bindings during work, not multi-day account retention.
Turn timing was observed but tasks differed, so it is not a controlled latency
comparison. Private native histories, prompts, IDs, account labels and receipts
remain outside this public repository.

## Remaining acceptance

Useful work, ordinary follow-up tools, account attribution, model/effort retention,
provider cache reads and native-history resumption after deployment passed. The
source/synthetic/Windows checks above remain separate evidence. No further restart
was induced with the test's newly warm bindings. Deliberately forced shortened
Claude continuation/missing-alias recovery and opaque cross-account failover are
still unproved. Do not interrupt active inference to turn those unknowns into
passing claims. The integration remains a draft; no source promotion or rebuild
was performed for this evidence-only pass.
