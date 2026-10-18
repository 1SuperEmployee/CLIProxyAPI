# Maintained fork

`se/main` is the maintained branch; inherited `main` is upstream history.

18 October: personal-fleet release `8.0.32-se-affinity.10-rc1` is accepted for
native Claude Code and Codex HTTP/SSE operation. Production source is `0611ffeb`,
on official v8.0.32 at `911f49cb`. Later commits record qualification/docs only.
Exact Windows CGO executable SHA256:
`e9b3bfcac078ab9cbb258b1cc018995ae466fc4b8c59f5596eba0a20621ce5d5`.

Read FORK.md, AFFINITY-PERSISTENCE.md and maintenance/qualification.md. The release
preserves durable affinity, nearest-reset selection and the optional ledger bridge.
Config integration keeps snapshot ordering plus the upstream stale-save fence.
Reviewed narrow repairs detect truncated native Claude streams and preserve
supported compaction request forms. The metadata-lock correction remains.
Independent review found and corrected missing snapshot hashes and false terminal
recognition before deployment. Tests retain those regressions.

Linux full suite passes 102 packages and relevant race checks pass 10 packages.
Actual Windows passes 2,000 top-level tests across 11 groups, no failures/skips.
All 12 server/test hashes match. The matched operator passes 101 offline tests and
actual Windows deadline, PID and pre-stop validation. The separate empty-auth
candidate started and round-tripped TLS settings without provider requests.

Bounded completed-turn cutover restored all 21 nonexpired schema-1 binding groups,
kept the console/environment and restored the recorder. Six workers retained their
process identities through the gateway switch. Separately reloaded client helpers
preserved all 174 saved native records/configurations/generations and harness
settings. Real native tools/follow-ups passed; the final sample has 41 attempts,
zero failures and cache reads. A recently bound conversation retained its account
across restart. Older probes whose aliases had expired establish native resumption,
not retention of an expired binding. Actual WhatsApp sends completed after helper
repair. Private identities, logs, receipts and credentials remain outside this repo.

Immediate software rollback is exact retained `8.0.28-se-affinity.8-rc2`, production
`bff9dc70`, SHA256
`86e399ce7665148981232174adad6a0471f3eac0d4730f9fe4cc80bafd05cc6a`.
Use the matched v10/v8 operator, preserving current OAuth, checkpoint, environment
and console. Do not invoke historical v8/v6 or panel/XML rollback procedures.
No automatic rollback or scheduled updater is installed.

Nearest-reset defaults off in generic code; the personal fleet keeps its existing
opt-in. Healthy valid bindings remain authoritative. Quota snapshots are cached;
only missing/stale new assignments coalesce bounded reads, with a two-second total
wait and ordinary eligible fallback. No polling service or new credential owner.
The optional fresh TLS handshake bound remains inactive in the personal fleet.
WebSockets, natural long-context compaction, native-child attribution and in-flight
recovery retain their explicit qualification limits. A small same-PC selector
benchmark shows no median regression, not an end-to-end latency guarantee.
