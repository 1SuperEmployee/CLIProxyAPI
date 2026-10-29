## Accepted RC3 ranking correction, 29 October 2026

RC3 layers a small remaining-weekly-allowance/time-to-reset ranking change onto
qualified v45 RC2. Source `bc8e05ae808b`; exact 74,332,672-byte executable SHA256
`78dd53fa199997b95cf6e0e724d11702d3bbca1027de5787083c99329ba9c201`.
Go 1.27.3, Windows amd64 CGO, MinGW 12-posix and timetzdata are unchanged.

Full auth/SDK Linux race packages pass; focused quota ranking/concurrent fleet
checks pass twenty repetitions. The twenty-account/six-conversation fixture now
prefers 97% remaining over 73h to 9% over 48-67h and verifies unchanged bindings,
shared lookup coalescing and passive updates. A repeat exposed the inherited
fixture's mixed real/fake clock near a second boundary. The test-only controlled
clock correction passed review; preserve the initial failure log. Actual Windows
full auth/SDK packages pass 1,389 top-level tests, no failures/skips. Three binaries
and 2,189 source fixture hashes match; empty-auth startup passes 401 with zero
provider requests. Prior RC2 4,923-test/full integration results below remain
identified as RC2 results. No broad unchanged integration campaign was repeated.

No dependencies/imports changed. Exact RC3 binary scanning reports only the prior
GO-2026-5932 legacy OpenPGP finding, qualified by the reviewed unimported legacy
packages and maintained replacement. 112 offline operator tests and six actual
Windows pre-stop guards pass. Independent Astra review found no blocker.

The approved 08:24 UTC gateway-only switch retained 458 nonexpired checkpoint groups
and fresh runtime logs proved active persistence. Console bytes stayed identical,
the recorder reconnected, and existing official Claude/Codex native IDs passed
incremental streaming, file tools and same-account follow-ups with cache reads.
Codex normal Stop/same-ID resume passed. No worker, credentials or transport flags
changed. Immediate recovery is exact accepted v13; current data stays in place.
Long-term distribution, all business delivery and actual WebSocket adoption remain
outside this bounded acceptance.

## Qualified v45 candidate, 27 October2026, not live

Server `8.0.45-se-affinity.15-rc2`, source 171759b6977f, combines official v8.0.45
at ad8bdd3f with the retained routing/snapshot/compaction protections, upstream
PR6515's interrupt repair and PR6755's close/join/clear streaming-log repair.
The latter reproduced under the race detector. Keep both regressions and retire
equivalent fork deltas after upstream incorporates them.

Build: Go 1.27.3, Windows amd64 CGO, MinGW 12-posix and timetzdata. Exact 74,332,672-byte
executable SHA256:
`b3abcba0b493fb641b8cfd2cc58b67f5d1d8ac5159aa5053ec69fdd30d2c995c`.

RC2 production code passes the 102-package ordinary Linux suite. The nine relevant
race-package rerun passes, as do 30 keepalive-upload repetitions, 30 interrupt
repetitions and 50 logging repetitions. Initial failures remain recorded: inherited
catalog assumptions, changed zstd encoder framing, JSON diagnostic wording, and
fake upload completion before receipt. The corrections preserve the assertions.
The Windows-only SIWC tilde fixture additionally needs USERPROFILE as well as HOME;
test-only 78a3f805 makes that one-line correction, and its Linux fixture passes ten
race repetitions. The unchanged server remains built at 171759b6. Frozen fixture
source 78a3f805 differs from that server source in this one test line only.

Actual Windows passes 4,923 top-level tests in 14 complete packages, with zero
failures. Four explicit skips cover two unreadable-file/dir fixtures and two
helper-process entry points. All 15 transferred executable hashes and 2,189 source
fixture hashes match. A separate empty-auth, local-model candidate starts on an
assigned loopback port, enforces 401 authentication and exits afterward, with zero
provider requests. No live config or provider credential is mounted into the test.

After two newly published HTML advisories, RC2 updates x/net 0.62 and its required
crypto/sync/sys/text versions. compress 1.18.7 remains. Linux and Windows server-
source scans report zero called/imported vulnerabilities and one required-module
warning. Exact Windows PE scanning removes both HTML warnings but still reports
GO-2026-5932, legacy golang.org/x/crypto/openpgp. The actual Windows import graph
excludes every flagged package and uses the maintained ProtonMail fork recommended
by the advisory. Preserve this source/binary distinction; the binary scan is not
empty. Initial RC1 source results alone were not used to accept the dependency
inventory. Interrupted/OOM/noexec scans are not accepted evidence.

Independent Astra review found and corrected an inherited maintenance acceptance
gap. Retained checkpoint JSON and401 do not prove that durable routing loaded.
The matched operator distinguishes checkpoint retention from runtime acceptance;
the separate read-only startup gate requires a new process and fresh version/active-
persistence log bytes with the retained unexpired-group count. It records Running
dependent services before a possible stop failure. Explicit stopped-service recovery
preserves current OAuth/checkpoint/console and does not queue unattended rollback.
112 offline operator tests pass; the eight new startup-proof tests also pass on
Windows. Real-PC unchanged-PID and old-log rejection checks pass without a restart.

Live remains accepted v13/v42. Current T3 roles replace the old fixed Buzz roster as
the maintenance inventory. Coordinate a quiet window, warn Su for Usman and recheck
current consumers immediately before any stop. Preserve healthy account bindings,
console, credentials and client settings. v45 changes negotiated downstream response
compression, including SSE, and upstream response decoding even while request
compression stays off. Post-cutover acceptance must therefore observe actual native
incremental text, tools, normal Stop/same-ID resume, account/cache continuity and a
fresh recorder feed. Pre-upgrade native results are a baseline, not acceptance of
this candidate. No WebSocket account/client flag is enabled by this qualification.

## Accepted v35 qualification, 20–21 October 2026

Code `b49ded51` merges stable officialv8.0.35 (`ae3ee3be`) into the acceptedv33
fork. Source/artifact qualification completed on 20 October; the personal gateway
has run `8.0.35-se-affinity.12-rc1` since 21 October, 07:36 UTC. Three newer unreleased main commits are reviewed and
excluded. Independent review is in `reviews/2026-10-20-v8.0.35.md`.

No blocking conflict was found in the native Claude Code/Codex HTTP/SSE paths.
Keep both bound-digest and observation-time fences, SourceHash/config ordering,
durable schema1 affinity, bounded nearest-reset, optional ledger bridge, setup-token
metadata synchronization, strict buffered message_stop and native compaction fixes.
The new incomplete-file watcher and latest-record Modify paths complement those
protections. Codex proactive/rejected-token refresh retains version fences. New
optional SIWC/search/provider paths are not enabled by this update.

Linux full suite passes102 packages;29 packages have no tests. All10 targeted race
packages pass. The isolated builder used Go1.27.1, MinGW12-posix, CGO enabled and
2CPU/4GiB limits. Its2062 source files matched the archived source tree with no
extra/missing/changed source files. Actual Windows passes2,410 top-level tests in
15 groups, zero failures/skips:13 complete packages plus recorded relevant executor
and helps selections. All16 transferred server/test hashes match. A
separate candidate process with an empty synthetic auth directory starts on a
private test loopback port and passes TLS0/10/0. No provider requests or real
configuration/credential files were used; the live gateway identity stayed intact.

Exact72,058,368-byte Windows artifact `8.0.35-se-affinity.12-rc1` SHA256:
`21a47f41fb11623e70573ff2e3a15a341f34140a9d126a334a37e6ebc1d04227`.
The matched operator candidate changes only the three-line executable/hash
transition map, recovering to acceptedv33 SHA256
`b7c5b100dcb849dfadd081c5820afa0a71335b3cb15c3edc2af594baa615df36`.
Independent executable-map review, 104 operator tests and actual Windows
pre-stop guards pass. The roster also recognises removed retired containers after
complete enumeration, preserves archive scope and refuses busy optional job runners.
The quiet-window gates passed before and during holds. The executable-only cutover
restored all 48 valid binding groups, console/environment and recorder. All six
worker PIDs/harness settings and 241 saved native identities/configs/generations
matched. Four existing conversations passed tools and follow-ups; two bound native
Claude/Codex conversations retained accounts with cache reads. A real owner-only
message has sent-log confirmation. The first 37 post-upgrade attempts had zero
failures. Read-only rollback ValidateOnly passes with the gateway unchanged; no
rollback was executed. No worker/helper restart or login/model/permission change.

A same-PC serial durable-selector/checkpoint benchmark used three samples of100
selections per2/32/256-binding case, GOMAXPROCS2. Acceptedv33 medians were
3.162/2.157/2.341ms; candidatev35 medians3.124/2.083/2.270ms, zero routing errors.
No regression appeared in the bounded samples. Small differences do not establish
end-to-end speed or throughput gains. The fake-auth candidate exited and its
listener is closed. No additional stress test is needed for this candidate.

One source-level SIWC hardening opportunity is documented in the review: its
explicit binding write bypasses our lastSavedAt marker for an unbound snapshot.
The normal active watcher is bound-digest protected and no reachable concurrent
unbound live snapshot path was found. It does not justify a new fleet patch now.
If that inactive combination is introduced later, add one marker/regression repair.

The live checks establish bounded completed-turn maintenance and native
continuation. They do not prove real OAuth rotation, natural compaction, all-seat
business delivery, native-child account attribution, general in-flight recovery or
end-to-end latency. Private operational receipts do not belong in this public
repository.

## Accepted v33 qualification, 18 October 2026

Production source 1289d6bd integrates official v8.0.33 at 2422c6b0, published during
v32 closure. Relevant upstream changes cover Claude preserved/adaptive thinking,
short per-minute cooldowns, auth file snapshots, model registration/exclusions and
complete SSE parsing. Default mDNS still rejects the production loopback listener;
no discovery configuration change was needed.

Independent review caught a digest-rejected snapshot returning applied=true; the
existing stale-snapshot sentinel now preserves the public contract. Both digest
and observation-time fences are retained under the auth mutation gate, since their
coverage differs. A retained regression caught upstream's missing message_stop
requirement in buffered SSE validation; it remains required. Top-level completion,
SourceHash, compaction and payload-final-barrier regressions pass.

Linux full suite 102 packages and relevant race suite 10 packages pass. Mac server
build and focused/full relevant suites pass. Actual Windows 2,237 top-level tests in
14 groups pass with no failures/skips; all 15 server/test hashes match. The initial
launcher referenced provider/claude, a package with no standalone tests, and aborted.
It was corrected to use the real parent thinking tests and precheck all executables;
only the complete rerun is accepted. Empty-auth startup and TLS 0/10/0 pass without
provider requests or live settings changes. Matched v11/v10 operator passes 101 offline
tests and actual Windows six guard and nine bounded-child checks.

Exact 8.0.33-se-affinity.11-rc1 Windows CGO SHA256:
`b7c5b100dcb849dfadd081c5820afa0a71335b3cb15c3edc2af594baa615df36`.
Go 1.27.1, MinGW-posix GCC 12. Immediate rollback is exact accepted v32, production
0611ffeb, SHA256
`e9b3bfcac078ab9cbb258b1cc018995ae466fc4b8c59f5596eba0a20621ce5d5`.

Independent source and exact artifact/operator clearance preceded the 13:45 UTC
bounded switch. Before/held gates passed, all 25 complete nonexpired checkpoint
objects were restored exactly, console/environment retained and recorder restored.
All six worker identities/harness/settings and 174 pre-existing saved records were
unchanged. Four already-bound native Claude/Codex sessions retained their accounts
through tools and follow-ups:31 attempts, zero failures and 2,412,230 cache-read tokens.
Both Chiefs' WhatsApp status calls passed; an actual reply afterward has completed
request, real ID and sent-log evidence. Final native idle check passes. A read while
the delivery turn was active correctly refused readiness; it was not a persistent
mismatch. No Chief/helper restart or live setting change was made for v33.

Private receipts/identities stay outside the public repo. This establishes bounded
completed-turn maintenance and native continuation, not natural compaction, child
attribution, in-flight recovery, downstream-WebSocket acceptance or a latency gain.
Front-door pre-execution rejections and partial token detail on rejected buffered
SSE remain observability follow-ups. No automatic updater or rollback is installed.

## Historical v32 qualification, 18 October 2026

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
race detection. Linux full suite passes 102 packages and relevant race checks pass10.
Actual Windows passes2,000top-level tests across11groups, zero failures/skips. All12
server/test hashes match the build. Empty-auth Windows startup and TLS setting
0/10/0 round trip pass without provider traffic or a live-service change. The
matched v10/v8 operator passes 101 offline tests and Windows deadline/PID/pre-stop
ValidateOnly guards; independent Astra review clears the source/artifact/operator.

Exact server8.0.32-se-affinity.10-rc1 SHA256:
`e9b3bfcac078ab9cbb258b1cc018995ae466fc4b8c59f5596eba0a20621ce5d5`.
Same-PC serial durable-selection medians at2/32/256sessions were3.091/2.037/2.273ms,
versus3.157/2.082/2.280ms on v28. One v32 sample at256sessions took4.644ms.
This small synthetic check finds no median regression; it is not an end-to-end
latency or load guarantee. Fresh cutover gates and bounded real worker acceptance subsequently passed as recorded below.

18 October deployment acceptance: exact artifact started with active persistence and
21 restored groups; all original checkpoint group objects survived. The six active
workers kept their PIDs through the switch. Subsequent caller-helper reloads retained
174 saved records/configurations/generations and all harness/settings. Four existing
native conversations passed tools/follow-ups. Their expired gateway aliases were
new assignments; they are not active-binding-retention proof. A separate recent
bound conversation kept the same account across restart. The final bounded ledger
sample has 41 attempts, zero failures and 2,829,131 cache-read tokens; real WhatsApp
sends also completed. This accepts normal HTTP/SSE completed-turn operation, not
natural compaction, downstream WebSockets, general outage recovery or load behavior.
A helper-client restart sequencing deviation is preserved in private fleet records.

Do not treat source qualification as live deployment. Private receipts stay outside
this repository. The personal fleet now runs the qualified v32 artifact.

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
