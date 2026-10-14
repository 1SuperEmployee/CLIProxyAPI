# Maintained fork

`se/main` is the maintained branch; inherited `main` is upstream history.
The current personal-fleet artifact is `8.0.23-se-affinity.6-rc1`, built from
`04248f92` with upstream main pinned at `3de4e248`. Latest published stable remains
v8.0.28 as of 13 October; the running artifact remains on the reviewed v8.0.23
base and main fixes. Do not call a fetched release a deployment.
Read [FORK.md](FORK.md) and [current qualification](maintenance/qualification.md).

The Windows cutover on 10 October retained 27 nonexpired bindings without changing
OAuth. Actual Claude Code/Sonnet and Codex/Luna useful work and post-restart
same-native-session/account follow-ups passed with cache reads. All16 workers were
released without restarting them. The optional authenticated loopback ledger bridge
supports the privately maintained static console. Neither the reader nor the console
is on the inference path.

Immediate software recovery retains the prior v5 executable and panel. The older
v8.0.16 artifact remains the earlier fully accepted baseline. Forced continuation,
failover and broad load/outage paths remain unproved by this bounded release check.
No automatic provider updater or scheduled reviewer is installed. Live configuration,
provider credentials, checkpoints and private fleet records never belong here.

`se/quota-selection-policy` is an unaccepted development branch with an opt-in
shared-quota preference integrated into Manager/session affinity. `routing.nearest-reset`
defaults off and requires session affinity. Shared Claude/Codex observations are
normalized in memory; stale cold-selection candidates share bounded usage GETs.
There is no poller, separate credential refresh or durable quota store. Healthy
bindings remain authoritative and network waits occur before affinity pick locks.
After waiting, eligibility is gathered again and the atomic binding check remains.
Supplemental model limits retain existing rejection/cooldown handling. Synthetic
source qualification passed apart from the unchanged local discovery failure;
see [qualification](maintenance/qualification.md). This branch is not merged or
deployed; observation and controlled host/runtime qualification remain open.

13 October source qualification adds twenty-account/six-conversation tests for
both providers, one failed quota endpoint with passive recovery, and a reset during
lookup. Arrival timestamps now age passive quota readings instead of stream
completion; older completed streams cannot replace newer observations. The
capture-time guard follows the reviewed idea in upstream PR6562, without adopting
its different most-remaining policy. No live routing change is implied.

`se/review-v8.0.28` integrates official stable 0d18d3b3 after that quota candidate.
The durable checkpoint schema remains unchanged; optional upstream rebinding and
new tool-search/body-limit/heartbeat options retain their default behavior.
Read maintenance/reviews/2026-10-13-v8.0.28.md for source review, fixture corrections,
the independently checked bounded PR6603 reload/auth-snapshot repair, and unresolved
interruption-test timing. Root and independent fifty-repeat Issue6559 checks pass.
Linux full suite passed 102 packages. Windows passed 1485 top-level tests across
six packages without skips; the isolated Windows server's v8 setting round-trip
passed. Candidate artifact is `8.0.28-se-affinity.8-rc2`, production source bff9dc70;
ae48adb1 changes tests/docs only. Exact checksums and timing limits are in the
qualification record. Live choice/timing, a matching bounded release operator,
real useful-work follow-ups and enablement remain open. se/main is unchanged.
This candidate is not a deployed release.
