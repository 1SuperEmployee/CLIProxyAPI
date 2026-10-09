# Maintained fork

`se/main` is the maintained branch; inherited `main` is upstream history.
The accepted personal-fleet artifact is `v8.0.16-se-affinity.4-rc1`, built from
`a4aa7fb5c56d48a5737fb7334d2b290080157146`. Later documentation commits do not
change that executable. Read [FORK.md](FORK.md) before building or updating.

The [9 October upstream review](maintenance/reviews/2026-10-09-v8.0.23.md)
compares upstream `v8.0.16` with `v8.0.23`. `se/review-v8.0.23` is an isolated
integration candidate. Its verified `8.0.23-se-affinity.5-rc1` artifact is now
running in the personal fleet at the operator's request. Actual Claude Code and
Codex useful work and same-native-session follow-ups passed on 9 October. This
was a bounded normal-work check; deliberate continuation/failover branches and
a further warm-binding restart were not exercised. The last fully accepted
rollback remains the v8.0.16 artifact above. Follow
[current qualification](maintenance/qualification.md)
for broader test results, environment limits and remaining actual-client checks.
An upstream cache-mark defect was reproduced offline.

No automatic updater or scheduled reviewer is installed. No live configuration,
provider credentials or private fleet records belong in this repository.
