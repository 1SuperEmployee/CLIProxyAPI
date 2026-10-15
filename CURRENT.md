# Maintained fork

`se/main` is the maintained branch; inherited `main` is upstream history.
The personal-fleet artifact deployed14October is `8.0.28-se-affinity.8-rc2`,
production source bff9dc70, official stable v8.0.28 at0d18d3b3 plus reviewed fork
patches. ae48adb1 changes tests/docs only. The exact Windows SHA256 is
`86e399ce7665148981232174adad6a0471f3eac0d4730f9fe4cc80bafd05cc6a`.
Read FORK.md, AFFINITY-PERSISTENCE.md and maintenance/qualification.md.

Bounded completed-turn cutover restored23 nonexpired schema-1 bindings, preserved
OAuth and the static console, restored the running recorder and released the six
current active workers without changing their process identities. Durable storage
initialized active. Real Claude Code/Sonnet and Codex/Luna continuations completed
useful tools in their previous native conversations, with one account per thread
and later provider cache reads. This is not mid-turn failure, load or host takeover
acceptance. Private fleet identities, credentials and receipts stay outside this repo.

Immediate software rollback is the exact retained `8.0.23-se-affinity.6-rc1`
artifact, source04248f92, SHA256
`a51b27274cead678ff5ae95cfe358bec812ba9d071261c35a6af509a3e8ed27d`.
The matched bounded v8/v6 operator changes only the executable; preserve current
OAuth, checkpoint, environment and console. Do not run the old v6/v5 panel/XML
procedure against this pair. There is no automatic rollback or scheduled updater.

`routing.nearest-reset` code is present and defaults off. It requires session
affinity; healthy bindings win. Shared quota observations are cached in memory.
Missing/stale new assignments coalesce parallel fixed-endpoint usage GETs within
a two-second total wait, then safely fall back on incomplete data. No poller,
credential owner or model-limit ranking subsystem was added. The personal fleet
opted in14October: new Claude Code/Sonnet and Codex/Luna conversations chose the
expected accounts from live quotas and retained them through useful tools and
follow-ups, with cache reads. Activation retained existing bindings without a
restart. One older Codex conversation recovered a TLS failure through normal
account failover while preserving its native conversation. This is bounded live
choice/continuation acceptance, not a controlled lookup-overhead or load benchmark.

Latest official stable checked14October is v8.0.29. It is reviewed, not integrated
or qualified. Its registration-publication and queued delete/recreate fixes cover
races outside our bounded PR6603 repair. Independently applied upstream regressions
reproduce three publication failures here. Qualify that update before account
prefix changes, deletion/re-enrollment or plugin registration work. Its optional
TLS handshake protection defaults off. Prioritize qualification against the
observed Codex TLS failure; transport and a fix are not yet established.
See qualification for remaining scope.
