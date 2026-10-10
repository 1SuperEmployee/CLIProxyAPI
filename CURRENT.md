# Maintained fork

`se/main` is the maintained branch; inherited `main` is upstream history.
The current personal-fleet artifact is `8.0.23-se-affinity.6-rc1`, built from
`04248f92` with upstream main pinned at `3de4e248`. Latest published stable remains
v8.0.23. This includes reviewed unreleased fixes; do not call it a newer stable.
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
