# v45 security and WebSocket candidate, qualification in progress

Current candidate `8.0.45-se-affinity.15-rc2`, code `320482eb5b59`, merges official
v8.0.45 at ad8bdd3f20be with PR6515's interrupt repair and PR6755's streaming-log
close/join/clear repair and regression. RC2 also patches newly announced HTML-library
issues with x/net0.62 and its required crypto/sync/sys/text upgrades. Live remains
accepted v13/v42. Earlier v14 and RC1 results do not qualify RC2.

Superseded RC1 at7c343124 passed102 full Linux packages, nine race packages,
30 interrupt repeats and50 logging repeats. Both target-source scans reported no
vulnerable called symbols; its exact Windows binary scan still flagged two new
x/net HTML advisories and the existing module-wide OpenPGP advisory. RC1 was never
deployed. RC2 is undergoing fresh tests/scans; do not claim a clean scan from RC1.

Retain compress1.18.7, Go1.27.3, Windows CGO/MinGW12 and timetzdata. Both snapshot
fences, SourceHash/order, durable schema1 bindings, nearest-reset selection and the
optional ledger bridge remain. No account flags, client transports, credentials,
console or persistence schema are changed. Late terminal interrupts are silent;
queued-input loss checks still refuse unsent stale work. The old fixed Buzz roster
is stale and must not authorize maintenance. Exact artifact, recovery guard and
post-cutover native acceptance remain pending. Current qualification evidence stays
outside the public source repository.

# Qualified WebSocket candidate, 26 October 2026

Candidate `8.0.42-se-affinity.14-rc1`, code `8db9efad6636`, ports upstream PR6515
at `5d8de6f6e9c218d178fb42efce55732c61e88a7a`. PR6515 is open, not merged.
It repairs local HTTP-fallback interrupt identity, queued-frame claims and terminal
publication, including synthetic prewarm history bounded to256 IDs per connection.
No dependencies, OAuth handling, routing policy, checkpoint schema or console change.

macOS and Linux full suites pass102 packages each; focused handler interrupt race
checks pass50 repeats on each. Windows exact-artifact tests pass331 handler,
1,744 executor and1,088 auth top-level tests. The initial Windows executor run had
one Antigravity credit-hint failure;30 isolated repeats pass on both unchanged v13
and candidate, and full baseline/candidate reruns pass. Its cause remains unproven,
not fixed or silently omitted. Empty-auth Windows startup passes401 on an assigned
loopback port with no provider request. Six bounded pre-stop guards and104 offline
operator tests pass without changing the live service.

Windows Go1.27.3 CGO/MinGW12, `-tags timetzdata`, exact executable SHA256:
`4d5e5e24723c8d9a33514aead01e9aac55d0e07c459f8198860a4e7612cdc8b6`.
The candidate is staged, not active. Live remains accepted v13 below. The matched
v14/v13 transition is staged and read-only validated; no restart is authorised
until a quiet window is coordinated. Client WebSocket flags and Codex account
WebSocket metadata remain unchanged. Native tools, follow-up, cancellation/resume,
account/cache attribution and actual transport/latency acceptance are pending.

# Live v42 security update, 25 October 2026

Source `db965301` integrates official stablev8.0.42 at `bc064340` with preserved
fork routing/snapshot protections and a patched compression dependency. Go1.27.3
is the qualification compiler. Full macOS suite, relevant race checks, repeated
merge-boundary regressions and compilation pass; source vulnerability analysis
finds no vulnerable called symbols. Linux102-package full suite and five relevant race packages pass. The exact Windows
CGO rc2 artifact passed4,141 top-level tests, with zero failures and four explicit
platform/helper skips; isolated empty-auth startup also passed.
Windows rc1 exposed absent IANA timezone data. The rc2 rebuild uses standard
`-tags timetzdata`; the initial rc1 artifact is not deployment-qualified.
Exact Windows73,390,592-byte SHA256:
04e19e9828879227fe51edb436638291c366f87b39c52cdc1a706ba932fb6dd4.
Live activation at13:24 UTC restored169 nonexpired schema-1 routing groups,
kept OAuth and the console, and reconnected the recorder. Independent Astra
review found no source/artifact blocker and repeated relevant regressions and
104 operator tests. Actual official Claude and Codex file tools and exact-native-ID
continuations passed after a separate client-host reboot: eight successful
attempts, one unchanged account per conversation and cache reads. This is bounded
acceptance, not proof of every business delivery or natural compaction path.
Immediate recovery is exact retained acceptedv12, using the v13/v12 Windows
transition. The old fixed SE Buzz fleet roster is stale for the new T3 runtime;
independently establish current consumers/idleness before any later stop, and
never reuse a gateway hold controller for host reboots. Read
`maintenance/reviews/2026-10-25-v8.0.42.md` for evidence and limits.
Earlier accepted releases below are history and retained recovery artifacts.

# Accepted v35 release, 21 October 2026

The personal gateway now runs `8.0.35-se-affinity.12-rc1`, production code
`b49ded518f63d029fd60c249bffbab8cdcd3c480`, merging official v8.0.35
`ae3ee3be303fb07fb8173279923413005bf32177`. `se/main` is the maintained branch.
Later documentation commits do not change the deployed executable.

Exact Windows amd64 CGO artifact, 72,058,368 bytes, SHA256:
`21a47f41fb11623e70573ff2e3a15a341f34140a9d126a334a37e6ebc1d04227`.
Built with Go 1.27.1 and MinGW 12-posix. Independent review, Linux 102 full
packages and 10 race packages, actual Windows 2,410 top-level tests in 15 groups
with no failures/skips, and all 16 executable hashes passed. Synthetic empty-auth
startup and TLS 0/10/0 passed. Three newer unreleased commits and new optional
provider/search/SIWC features were excluded.

At 07:36 UTC the guarded executable-only switch restored all 48 valid account-
binding groups, preserved the console/environment and restored the recorder.
All six worker PIDs, harness hashes/settings and 241 prior saved conversation
identities/configurations/generations matched. Four existing native conversations
passed actual tools and same-thread follow-ups. Two bound Claude/Codex conversations
retained accounts across the switch, with cache reads. A real owner-only message
has sent-log confirmation. The initial post-upgrade sample had 37 attempts and zero
failures; this is a bounded sample, not a future availability guarantee.

The matched operator passes 104 tests and actual Windows pre-stop guards. It also
accepts verified removed retired containers while preserving archive scope, refuses
reactivation and requires optional job runners idle. Earlier live background child
work correctly blocked cutover. No worker/helper restart, login/model/permission
or persistent VPS setting changed. Native child scanner coverage remains separate.

Immediate software rollback is exact accepted v33, production `1289d6bd`:
`b7c5b100dcb849dfadd081c5820afa0a71335b3cb15c3edc2af594baa615df36`.
Use the matched v12/v11 operator and preserve current OAuth, checkpoint, console
and environment. Read-only rollback ValidateOnly passes against the current running
gateway; no rollback was executed. Historical transition pairs do not apply.

Keep both snapshot fences, config ordering/SourceHash, durable schema-1 affinity,
bounded nearest-reset selection, ledger bridge, setup-token synchronization and
strict native SSE/compaction corrections. Read FORK.md, AFFINITY-PERSISTENCE.md,
maintenance/qualification.md and maintenance/reviews/2026-10-20-v8.0.35.md.
Natural compaction, native-child account attribution, downstream WebSockets,
in-flight recovery and client-timed latency retain their separate proof limits.
Same-PC bounded selector samples showed no regression, not an end-to-end speedup.
Private operational receipts, identities and messages stay outside this repository.

# Accepted v33 history

`se/main` is the maintained branch; inherited `main` is upstream history.

18 October: personal-fleet release `8.0.33-se-affinity.11-rc1` was live and accepted
for native Claude Code and Codex HTTP/SSE operation. Production source `1289d6bd`
merges official v8.0.33 at `2422c6b0`. Later commits record tests/docs only.
Exact Windows amd64 CGO executable SHA256:
`b7c5b100dcb849dfadd081c5820afa0a71335b3cb15c3edc2af594baa615df36`.

Read FORK.md, AFFINITY-PERSISTENCE.md and maintenance/qualification.md. The release
retains durable schema-1 account affinity, bounded cached nearest-reset selection,
the optional ledger bridge, config snapshot ordering/SourceHash and the scoped
compaction payload corrections. Upstream's complete SSE parser replaces our
interim parser; strict top-level terminal and buffered message_stop regressions
remain. Both canonical file-digest and observation-time snapshot guards are needed
for their distinct bound-file and unbound/store paths. Rejected snapshots correctly
report applied=false. Independent Astra review cleared source, exact artifact and
operator before cutover.

Linux full suite passes 102 packages and relevant race checks pass 10 packages.
Actual Windows passes 2,237 top-level tests across 14 groups with zero failures/skips.
All 15 executable hashes match. The matched operator passes 101 offline tests plus
actual Windows deadline/PID/pre-stop and bounded-child guards. A separate empty-auth
candidate starts and round-trips TLS settings 0/10/0 without provider requests.
The first Windows launcher referenced a no-test package, aborted, and was corrected;
only the complete rerun is counted.

At 13:45 UTC the bounded gateway-only switch restored all 25 nonexpired checkpoint
groups exactly, preserved console/environment and restored the recorder. All six
worker PIDs, harness/settings and 174 pre-existing saved native records were unchanged.
Four previously bound Claude/Codex conversations retained their accounts through real
tools and follow-ups:31 attempts, zero failures and cache reads. An actual WhatsApp
reply completed afterward. Private identities, messages and receipts stay outside
this repository. No Chief/helper restart or settings change was made for v33.

Immediate software rollback is exact retained `8.0.32-se-affinity.10-rc1`, production
`0611ffeb`, SHA256
`e9b3bfcac078ab9cbb258b1cc018995ae466fc4b8c59f5596eba0a20621ce5d5`.
Use the matched v11/v10 operator; preserve current OAuth, checkpoint, environment
and console. Older transition maps do not apply. No automatic rollback or updater
is installed. Nearest-reset stays enabled for the personal fleet; healthy valid
bindings win. Fresh TLS handshake protection remains inactive; WebSockets, natural
compaction, native-child attribution and in-flight recovery retain separate proof
limits. This qualification does not claim an end-to-end latency improvement.
