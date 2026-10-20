# Qualified v35 candidate, 20 October 2026

Stable officialv8.0.35 (`ae3ee3be`) is merged at production-source commit
`b49ded518f63d029fd60c249bffbab8cdcd3c480`. Independent source review found no
blocking conflict in the native Claude Code/Codex HTTP/SSE fleet paths. Three
newer unreleased main commits are deliberately excluded. This is a source/artifact-qualified
candidate, not a claim that the running gateway has been upgraded.

`8.0.35-se-affinity.12-rc1` WindowsCGO artifact SHA256:
`21a47f41fb11623e70573ff2e3a15a341f34140a9d126a334a37e6ebc1d04227`.
Go1.27.1 / MinGW12-posix; exact code commit above,72,058,368bytes. Linux full suite
passes102 packages,29 packages have no tests; all10 targeted race packages pass.
Actual Windows qualification passed2,410 top-level tests in15 groups with no
failures/skips; all16 server/test executable hashes match. Separate empty-auth
startup and synthetic TLS0/10/0 pass, with live gateway PID unchanged. Independent
review of the three-line operator mapping change,101 offline tests and actual
Windows parse pass; no production transition or ValidateOnly was invoked.
The synthetic candidate exited and its test listener is closed. A bounded same-PC
serial durable-selector comparison at2/32/256 bindings had v33 medians
3.162/2.157/2.341ms and v35 medians3.124/2.083/2.270ms, with zero routing errors.
These short samples show no regression, not an end-to-end speed improvement. Read
`maintenance/reviews/2026-10-20-v8.0.35.md` and `maintenance/qualification.md`.

Keep both snapshot fences, config ordering/SourceHash, durable schema1 affinity,
bounded nearest-reset selection, ledger bridge, setup-token synchronization and
strict native SSE/compaction corrections. Upstream watcher/incomplete-JSON and
latest-record Modify protections complement those fixes; they do not replace them.
New optional search/SIWC/provider features are not enabled by this update.

Live remains the acceptedv33 artifact below. Do not run a v12 transition using a
historical v11/v10 map. Candidate recovery must point to exact acceptedv33
SHA256b7c5b100dcb849dfadd081c5820afa0a71335b3cb15c3edc2af594baa615df36
and preserve current credentials, checkpoint, environment, console and worker
settings. No gateway or worker has been restarted for this candidate.

# Maintained fork

`se/main` is the maintained branch; inherited `main` is upstream history.

18 October: personal-fleet release `8.0.33-se-affinity.11-rc1` is live and accepted
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
