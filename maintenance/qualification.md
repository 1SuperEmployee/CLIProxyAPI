# Current update qualification

This is the live qualification record, not a task queue. Completed source reviews
remain under `reviews/`. Read CURRENT.md for the accepted artifact.

## Current personal-fleet release, 10 October

Source04248f92 integrates upstream main3de4e248 and adds an opt-in authenticated
GET-only loopback ledger bridge. Latest stable remains v8.0.23. See the finished
main review; this release includes unreleased upstream fixes.

- macOS relevant suites and focused races passed. Full Mac suite passed100 packages,
  with29 no-test packages; discovery multicast failed identically on pristine upstream.
- Networking-enabled Linux full suite passed, including discovery; focused bridge,
  affinity and cancellation/cooldown races passed. No upstream test was weakened.
- Windows CGO server and auth-test binaries built using Go1.27.1 / MinGW-posix GCC12.
  The full auth executable passed958 top-level tests including10 durability tests,
  with tracked public fixtures and temporary synthetic state.
- Three 100-selection Windows benchmark runs per case, four Go processors: serial
  medians at2/32/256 bindings were3.271/2.147/2.402ms; prior v5 measured
  3.159/2.138/2.430ms. This short sample shows comparable selector behavior;
  it does not establish whole-request throughput or improved provider latency.
- Bridge rejection/authentication, method/path restriction, loopback-only destinations,
  redirects, query forwarding, timeout, malformed/oversized JSON and reader failure
  were tested synthetically. The reader owns response schemas; the bridge owns
  bounded authenticated transport, rather than duplicating the reader schema.
- Approved completed-turn cutover restored27 nonexpired bindings, retained OAuth,
  restored the running recorder and released all16 workers without worker restarts.
- Actual Claude Code/Sonnet and Codex/Luna conversations did useful tool/file work
  before and after restart in the same native sessions and same provider accounts,
  with cache reads, completed native turns and delivered follow-ups. Final maintenance
  inventory had16 running/unheld seats,287 histories and no blockers.
- The production static console rendered real quota and ledger data through the
  authenticated same-origin bridge. Its original layout is privately maintained;
  panel auto updates were already disabled. The development server is not required.

Artifact identity `8.0.23-se-affinity.6-rc1`, source04248f92:

```sh
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-posix GOOS=windows GOARCH=amd64 \
  go build -p 2 \
  -ldflags '-s -w -X main.Version=8.0.23-se-affinity.6-rc1 -X main.Commit=04248f92' \
  -o cli-proxy-affinity-v6-rc1.exe ./cmd/server
```

Server SHA256 `a51b27274cead678ff5ae95cfe358bec812ba9d071261c35a6af509a3e8ed27d`.
Auth-test SHA256 `096db80a03bcf9234490262156bdb1296b5e4109ba3c948fa5d851924a6300cb`.
Immediate software rollback retains v5 and the prior panel; v8.0.16 remains the older
fully accepted baseline. This check does not prove forced continuation/failover,
recovery during an active provider turn, full host takeover or broad load behavior.
No provider credentials or private records were supplied to the builder or this repo.

## Previous v8.0.23 qualification

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
## Shared-quota selection candidate

13 October extension: twenty accounts and six simultaneous distinct conversations
share one usage GET per account for both Claude and Codex. After passive reset
preference changes, all healthy follow-ups retain their account and a new
conversation uses the newer reading with no new GET. One unavailable usage endpoint
keeps the entire twenty-account baseline pool eligible, backs off, and returns to
quota ranking when ordinary inference supplies fresh data. Crossing a reset while
waiting cannot fabricate available quota. Focused race checks passed.

Upstream PR6562 was reviewed for edge cases. Its most-remaining policy and bounded
unknown exploration differ from this nearest-weekly-reset policy, so the selector
was not imported. Its header-arrival timestamp approach was adapted: per-field
capture times prevent long streams and partial events from freshening old quota,
and older completed streams cannot overwrite newer observations. Logging/auth
capture tests and focused races passed. This remains source qualification, not
real twenty-account provider traffic or live first-response timing.

Reference: https://github.com/router-for-me/CLIProxyAPI/pull/6562

The eea3ba35 Windows CGO auth executable passed focused quota, session-affinity,
manager-affinity and durability tests on the actual Windows host without stopping
inference. Server/auth builds used Go1.27.1 and MinGW-posix GCC12. Linux auth quota
races passed. The Linux full suite exposed the existing intermittent late-interrupt
test; it is retained as a qualification limit rather than skipped.

Independent source review identified a relative Codex reset edge. Shared
Reset-After-Seconds now becomes Reset-At using that field's own capture timestamp,
unless an equally/newer absolute reset is present. Conservative snapshot age is
preserved. Deterministic tests cover older unrelated fields, older absolute resets
and equally recent provider absolute values. Focused auth races passed after the
fix; the earlier Windows executable predates this follow-up.

The usage reporter starts after account selection. Its elapsed/TTFT fields do not
include the new quota wait. Actual first-response qualification must time from the
client's request, with cold/warm choices distinguished; unchanged ledger TTFT alone
cannot prove no regression.

11 October, `se/quota-selection-policy`. This is source qualification only; the
maintained release and live gateway have not changed. The earlier pure helper is
now connected to Manager and the existing session-affinity fallback through the
disabled-by-default `routing.nearest-reset` option. Session affinity is required.
Home and handled plugin scheduling remain authoritative. Mixed-provider or
non-OAuth pools retain baseline selection. A single comparable candidate needs no
quota lookup. Quota ranking follows priority, weight eligibility and the existing
cold Codex transport preference; healthy bindings still outrank all preferences.

### Implemented behavior

- Parse default shared weekly/short windows from Claude and Codex observations and
  read-only usage bodies. Validate units, dates, finite values and required fields.
  Identify Codex windows by duration, not primary/secondary names. Do not rank
  model, overage, credit or additional-limit namespaces. Shared denial constrains
  headroom without creating a new scheduler cooldown path.
- Ordinary responses continue updating the existing Auth quota observation.
  Normalize newer observations lazily when an assignment needs comparison. Cache
  reads do not extend timestamps; crossing a reset invalidates a reading.
- One bounded in-memory entry per account, capped at 256. Capacity overflow retains
  baseline selection. Credential version/registration changes invalidate prior
  reads, pending results and failure backoff. Carried old passive observations are
  not adopted after replacement. Future timestamps cannot poison successful reads.
- Only stale/missing comparative candidates trigger usage GETs, in parallel and
  coalesced per account. One waiter leaving does not cancel another's lookup.
  Failures briefly back off. No poller or independent OAuth refresh/persistence.
- GET targets are fixed to Claude `/api/oauth/usage` and Codex `/backend-api/wham/usage`.
  Existing executors prepare credentials and transports, including CodexAutoExecutor.
  Redirects, error bodies and bodies above 1 MiB are rejected. No quota body,
  credential, account email or prompt is logged by this feature.
- No Manager or explicit affinity pick lock is held across the quota network wait.
  After waiting, Manager gathers candidates again; selection still atomically
  rechecks the explicit binding before narrowing the fallback pool. A late reading
  cannot replace a newer passive reading or a replacement credential's entry.
- Candidate defaults are 15m freshness, 5% preferred reserve, a 2s total wait and
  30s failure backoff. SDK configuration exposes these values. They are candidate
  settings to observe before live acceptance, not measured provider tail guarantees.
  Equal resets and incomplete data retain baseline selection; all-low reserve is a
  preference fallback, not an invented hard provider rejection.

### Source checks

Passed: full auth, service, executor and executor-helper packages; focused race
checks covering quota normalization/cache/Manager routing, existing affinity and
Fable-only rejection; server compile. Synthetic checks cover fresh passive/no GET,
weekly-only and short windows, overlapping cold starts with one lookup per account
and one explicit binding, cancellation, parallel deadline, failure backoff,
reset boundary, replacement, old lookup/new passive ordering, binding/disabled
changes during a lookup, transport preference and baseline rotation.

The real ClaudeExecutor and CodexAutoExecutor were exercised through Manager with
synthetic HTTP transports. Each fetched once per account and used zero GETs on a
same-session follow-up. This is runtime code-path proof with fake credentials and
responses, not actual provider or Windows deployment proof. Lookup deadlines do
not propagate into inference contexts. The service constructor test proves the
feature changes selection only when configured and preserves the default.

`go test ./...` still fails only `internal/discovery.TestAdvertiserAndBrowser_Integration`
waiting for a local multicast advertisement. Discovery code is unchanged and the
earlier qualification reproduced that failure on pristine upstream. All other
packages passed. The affected packages and focused race checks were rerun after
the final source guards; do not call the whole suite clean.

Codex field interpretation was checked against the current primary source:
https://github.com/openai/codex/blob/main/codex-rs/codex-api/src/rate_limits.rs
The default shared header family and optional windows are distinct from additional
named families. The existing upstream quota-observation snapshot remains the
passive source; this candidate does not introduce a provider-wide scope catalog.

### Remaining acceptance

Compare proposed cold choices with current actual choices without changing live
routing. Then qualify Windows/runtime compatibility, useful work/follow-up affinity,
observed lookup latency and failure recovery before controlled enablement. No live
cap test, forced continuation, failover/load/outage acceptance or performance claim
is implied by these synthetic checks. The explicit-binding concurrency proof does
not establish stronger concurrent semantics for the existing LCP-only matcher.
