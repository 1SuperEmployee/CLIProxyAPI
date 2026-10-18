## v32 qualification, 18 October 2026

Candidate production source 0611ffeb merges official v8.0.32 at911f49cb into the
qualified v29 branch. Config reloads keep one snapshot and one apply lock; its
SourceHash also survives LoadConfigData, preserving the management stale-save
fence. Independent review found and corrected the missing snapshot hash and a
false terminal detector that could accept nested tool JSON.

Scoped corrections from upstream PR6708 cover native Claude stream truncation,
failed usage and single terminal SSE errors. Nonstream SSE conversion also requires
message_stop. PR6716 supplies the narrow Claude compaction-injection and Codex
compaction-rule compatibility repairs. Other-provider PR6708 changes are excluded.
The previously qualified setup-token metadata lock remains. No client, model,
permission, account policy, WebSocket setting or timeout is changed by this source.

Mac server build and full executor package pass. Config/watcher/management/Claude
handler/Codex translator suites pass; independent overlay regressions pass under
race detection. Linux full suite passes102packages and relevant race checks pass10.
Actual Windows passes2,000top-level tests across11groups, zero failures/skips. All12
server/test hashes match the build. Empty-auth Windows startup and TLS setting
0/10/0 round trip pass without provider traffic or a live-service change. The
matched v10/v8 operator passes101offline tests and Windows deadline/PID/pre-stop
ValidateOnly guards; independent Astra review clears the source/artifact/operator.

Exact server8.0.32-se-affinity.10-rc1 SHA256:
`e9b3bfcac078ab9cbb258b1cc018995ae466fc4b8c59f5596eba0a20621ce5d5`.
Same-PC serial durable-selection medians at2/32/256sessions were3.091/2.037/2.273ms,
versus3.157/2.082/2.280ms on v28. One v32 sample at256sessions took4.644ms.
This small synthetic check finds no median regression; it is not an end-to-end
latency or load guarantee. Fresh cutover gates and real worker acceptance remain
pending.
Do not treat source qualification as live deployment. Private receipts stay outside
this repository. The live personal fleet still runs v28 until the guarded cutover.

# Current update qualification

This is the live qualification record, not a task queue. Completed source reviews
remain under `reviews/`. Read CURRENT.md for the accepted artifact.

## v8.0.29 qualified and staged, deployment deferred, 14 October

Candidate0405a543 integrates official5fc580b1 and preserves the fork snapshot guard
inside the refactored locked update path. It adds one production correction:
Claude setup-token boolean reads now use the existing metadata lock. The same
read/write race reproduces on official v8.0.29 and the maintained v8.0.28 build;
the new focused regression fails on pristine upstream and passes100 race repeats
after repair. No new lock or executor-wide serialization was added.

The final Mac full suite passes except the unchanged multicast discovery test.
All eight focused Mac race packages pass. Linux full suite passes102 packages;
seven relevant auth/service/watcher/config/transport race packages pass. An extra
full downstream-WebSocket handler race run exposes the retained late-interrupt
boundary described below; it is not a clean eight-package Linux race pass.
Windows artifact `8.0.29-se-affinity.9-rc1` was built from0405a543 using
Go1.27.1 and MinGW-posix GCC12. Server SHA256 is
`029c96652ceabef9c1c43127af9709132a1cebf5d03ab2a348d19f354087dc33`.
Nine server/test executable hashes were verified after transfer to both the Mac
and the actual Windows host. Windows passed 1511 top-level tests with no failures
or skips: 1004 auth, 221 service, 76 watcher, 136 config, 49 token-store, 19 Codex
auth, four targeted TLS helper tests and two targeted executor TLS tests. The first
six packages ran in full; the last two ran only the new TLS cases. Tests used
tracked public fixtures and temporary synthetic state, never live credentials.
An initial PowerShell launcher split unquoted dotted test flags and executed zero
tests; that run is not counted. Quoting the flags fixed the launcher without
changing tests or executables.

A separate actual Windows candidate process with an empty synthetic auth directory
started on its own loopback port. The v8 management root TLS setting round-tripped
0, 10, 0 using a process-local synthetic management credential. No provider request
or live configuration was involved. The test process was stopped afterward; the
live gateway PID was unchanged throughout.

The exact server executable is staged unreferenced beside the running artifact.
The matched private v9/v8 operator passes 101 offline tests. Actual Windows script
parse, expired-deadline, wrong-PID and bounded pre-stop ValidateOnly guards passed
without changing the gateway PID. The retained v8 artifact is the candidate's
software rollback; this validates that transition, not a live rollback drill.
The owner deferred restart during active use. No workers were held, no service
configuration changed and the optional TLS bound remains inactive. Production is
still v8.0.28; live v9 specialist/continuation acceptance remains pending the
owner's restart window.

Independent review recommends this candidate for current HTTP/SSE clients within
the limitations below. The proposed root requests.tls-handshake-timeout-seconds
value is 10, to be enabled and read back only after an approved cutover. This bounds
each fresh TLS handshake after TCP/proxy CONNECT, not total retries, TCP setup,
response waiting or an established stream. Synthetic Codex HTTP tests prove its
scope; a real silent-provider ten-second cutoff has not been observed live.

Two test-only ordering repairs retain their assertions: the batch-registration
fixture no longer assumes worker order and joins cleanup; the HTTP cancellation
fixture observes its executor's cancellation before checking the flag. Independent
race stress passes100 repetitions of each. Real Codex HTTP silent-handshake and
successful Chrome-uTLS SSE tests cover bounded setup, retry/socket cleanup and
streaming beyond the handshake deadline. No provider credentials are used.

The retained downstream-WebSocket issue can emit a synthetic incomplete event for
an unknown interrupt during the terminal-frame/stream-close gap. Relevant production
files are unchanged from the accepted v8.0.28 fork; official v28..v29 changes do not
touch them. Pristine Mac repetitions did not reproduce it, which does not invalidate
the observed Linux failure. Current HTTP/SSE client acceptance does not establish
general WebSocket interruption correctness. Keep that path gated before adoption;
do not weaken its test or claim it repaired by this update.

## v8.0.28 source and live cutover, 13-14 October

Source integration follows eea3ba35 and pins official stable 0d18d3b3.
The qualified v8 artifact was deployed14October; the personal fleet subsequently
opted in to nearest-reset with the bounded live checks below. The default remains off.
See reviews/2026-10-13-v8.0.28.md for source changes and merge resolution.

macOS build passed. Auth/service/executor/logging focused race checks passed,
including twenty-account quota cases and existing durability/affinity. The final
full Mac suite passes all packages except the unchanged multicast discovery test.
Two catalog-dependent executor fixtures failed identically on pristine v8.0.28;
temporary explicit model definitions preserve their original assertions and pass
three repeats. Initial full-suite runs also exposed intermittent late-interrupt
test behavior; twenty isolated repeats and untouched upstream's full handler
package passed. Do not call that interruption path repaired or fully accepted.

Independent review found and verified the relative-reset timestamp correction.
It then independently reproduced and reviewed the bounded reload/auth-snapshot
repairs adapted from PR6603. Root and independent runs each passed fifty repeats
of Issue6559 after both repairs, with zero alias/prefix failures. Independent ten
race repeats of the new snapshot regressions passed. Final Mac races across auth,
service, watcher and config passed three repeats; server build passed. Full Mac
qualification still has the unchanged discovery failure. No failing test was
removed or weakened. See the finished upstream review for the excluded PR scope.

The13October checks below were synthetic. The14October bounded cutover and
actual-worker acceptance are recorded separately below. Enabled quota choice is
accepted there within explicit scope; incremental lookup timing remains open.
Do not infer either from source or default-policy benchmarks.

Production candidate source is bff9dc70. Networking-enabled Linux full suite passed
102 packages; focused auth/service/watcher/config/logging races passed. Windows CGO
server and six test packages built with Go1.27.1 and MinGW-posix GCC12. The actual
Windows host passed 1001 auth, 19 Codex auth, 49 token-store and 72 watcher top-level
tests. Config passed 131 after restoring its original working-directory layout,
config.example.yaml and tracked public contract fixture. No assertion was changed.
The full service run found the positive snapshot fixture clock ambiguity described
in the review. The test-only correction and explicit equal-time stale regression
are in ae48adb1, with no production source change after bff9dc70. Ten Mac and Linux
race repeats passed. Final Windows runs pass all 1002 auth and 212 service top-level
tests. Together with the other four packages this is 1485 top-level Windows tests,
zero failures and zero skips. Standalone tests use tracked public fixtures and
temporary synthetic state, never live auth or the live checkpoint.
The concurrent-save Issue6559 integration also passed ten actual Windows repeats.
An initial extra command targeted the API package instead of the service package
and ran zero tests; it was not counted as validation. The corrected command checks
the expected ten-pass count in addition to exit status.

A separate Windows candidate process with empty synthetic auth and its own
loopback port started successfully. Its v8 routing setting round-tripped false,
true, false through normal management GET/PUT. The live gateway PID was unchanged;
no provider traffic or production configuration was involved.

Artifact `8.0.28-se-affinity.8-rc2`, production source bff9dc70:

```sh
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-posix GOOS=windows GOARCH=amd64 \
  go build -p 2 \
  -ldflags '-s -w -X main.Version=8.0.28-se-affinity.8-rc2 -X main.Commit=bff9dc70' \
  -o cli-proxy-affinity-v8-rc2.exe ./cmd/server
```

Server SHA256 `86e399ce7665148981232174adad6a0471f3eac0d4730f9fe4cc80bafd05cc6a`.
Final auth-test SHA256, fixture source ae48adb1:
`68a060335e3672daf5decdb79c7020f41fd1b60513006081007b375af96c845e`.
Final service-test SHA256:
`ecfbd6ea69e4d87e28fdb067a97d1e1721ee28e752ad4c334333ca864319a2e7`.
PC copies match. The exact server artifact became the running executable14October.

Same Windows host, four Go processors, three 100-selection runs per case:
default serial durable-selection medians at 2/32/256 bindings were
3.157/2.187/2.510ms versus current v6's 3.229/2.601/2.703ms. All serial/concurrent
cases passed without account drift. One v6 256-binding sample was unusually short,
0.003ms; retain that variance rather than claim a controlled speed improvement.
These samples use the default affinity policy. They do not measure enabled quota
lookup, provider latency or complete client first-response time.

Independent review found no source blocker after the relative-reset guard,
bounded reload/auth-snapshot repair and explicit Windows fixture ordering. It
does not by itself prove live useful work, account attribution, lookup overhead or
outage recovery. The following live acceptance has narrower, explicit scope.

## Live acceptance and remaining scope, 14 October

Independent Astra source/operator review found no remaining source blocker after
two release-operator fixes: hash exactly the XML snapshot being replaced and detect
all numbered worker containers. The matching private operator passed94 tests and
actual Windows parse/expired-deadline/wrong-PID/bounded ValidateOnly guards.
Fresh active/parked roster, native-history and complete runtime lifecycle checks
passed before and during holds. Cutover retained all23 nonexpired bindings,
unchanged OAuth/console and six worker identities. Persistence initialized active.
Actual Claude Code/Sonnet and Codex/Luna useful tool work continued previous native
sessions; requests stayed on one account per conversation and later read cache.
Provider timing varies; these observations are not a controlled latency comparison.
Forced outages, cache retention guarantees and enabled quota overhead remain open.

The subsequent targeted v8 management PUT used raw JSON `true`. A GET confirmed
nearest-reset enabled, with session/subagent affinity, TTL and baseline selector
unchanged. Activation restarted no gateway or worker, retained all27 existing
nonexpired alias mappings immediately and left the recorder connected.
Fresh shared quota observations covered all four active accounts. New real
Claude Code/Sonnet and Codex/Luna conversations chose the expected eligible weekly
reset preference, then retained those accounts through tools/follow-ups:15 gateway
attempts, zero failures, with cache reads. Existing Claude continued on its account.
All four tested new/existing conversations completed without pending tools/delivery.

The existing Codex conversation's first attempt failed after155.210seconds in an
upstream TLS handshake/Windows socket read. Correlated trace logs classify the
failure; the gateway's500 is not evidence that the provider returned500. The same
request retried on another account8.892ms after that attempt ended and completed
the same native conversation. This is normal error-driven failover under unchanged
default `session-affinity-rebind-on:any`, not healthy-binding eviction or a155-second
quota lookup. Independent review accepted that explanation and found no activation
blocker. Exact low-level cause and transport remain unestablished; `wsarecv` alone
does not establish WebSocket use. No policy tightening or unrelated worker change
was introduced to hide the failure.

These are bounded account-choice/continuation checks. Usage GET round trips were
288–777ms, not isolated selection overhead. Relay-based reply timings mix dispatch,
model/tools and cross-machine clocks. Client-timed cold/warm incremental overhead,
load and controlled outage acceptance remain open. Concurrent personal-fleet work
outside these tested conversations was not restarted or declared idle.

Official v8.0.29 was released14October. Its optional TLS handshake bound defaults0.
Independent review reproduced three model-publication regressions on this candidate
using upstream tests. The queued delete/recreate race is a separate repair from
our PR6603 slices too. No evidence makes these a new regression in our artifact;
unchanged-account cutover is accepted. Qualify v8.0.29 next before account prefix
changes, deletion/re-enrollment or plugin registration. Do not claim all auth races
are repaired here. See the [official release](https://github.com/router-for-me/CLIProxyAPI/releases/tag/v8.0.29).
Prioritize qualifying the observed Codex connection path and a bounded handshake;
the new shared uTLS option is not yet demonstrated to fix this incident.

## Previous personal-fleet release, 10 October

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

Windows/runtime compatibility, live quota comparison, useful work/follow-up affinity
and bounded personal-fleet enablement are recorded above. Controlled client-timed
lookup overhead, live cap/load tests and forced-outage acceptance remain open.
The recovered TLS incident is evidence of one error-driven retry, not exhaustive
failure-recovery acceptance. These synthetic checks do not imply performance or
capacity guarantees. The explicit-binding concurrency proof does not establish
stronger concurrent semantics for the existing LCP-only matcher.
