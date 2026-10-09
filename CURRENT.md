# Maintained fork

`se/main` is the maintained branch; inherited `main` is upstream history.
The accepted personal-fleet artifact is `v8.0.16-se-affinity.4-rc1`, built from
`a4aa7fb5c56d48a5737fb7334d2b290080157146`. Later documentation commits do not
change that executable. Read [FORK.md](FORK.md) before building or updating.

The [9 October upstream review](maintenance/reviews/2026-10-09-v8.0.23.md)
compares upstream `v8.0.16` with `v8.0.23`. `se/review-v8.0.23` is an isolated
integration candidate, not a deployed release. Mac package/race checks and a
server build passed. Windows, full-suite and actual-client qualification remain
pending; an upstream cache-mark defect was reproduced offline.

No automatic updater or scheduled reviewer is installed. No live configuration,
provider credentials or private fleet records belong in this repository.
