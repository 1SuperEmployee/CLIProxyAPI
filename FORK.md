# Super Employee maintenance

The accepted release carries a small, opt-in durability patch on CLIProxyAPI v8.0.13.
This candidate branch merges upstream v8.0.16; it is not yet a deployed release.
Upstream remains https://github.com/router-for-me/CLIProxyAPI. Preserve its MIT
license and attribution. Keep provider integrations and ordinary fixes upstream.

## Source and releases

- `se/main` is our maintained branch. The inherited `main` is upstream history,
  not our deployment branch.
- `v8.0.13-se-affinity.3` identifies the deployed source commit
  `79919de730fff9571254c8a70aa91e7cc637ff22`.
- Releases `.1` and `.2` have an unsafe persistence error path. Use them only
  with `CLIPROXY_AFFINITY_STATE_FILE` empty; never restore their old durable XML.
- Later documentation commits do not imply a new deployed executable.
- `origin` points to this fork; `upstream` points to router-for-me/CLIProxyAPI.
- Our change is described in [AFFINITY-PERSISTENCE.md](AFFINITY-PERSISTENCE.md).
  Keep custom changes separate so an equivalent upstream implementation can
  replace them after verification.

Only source, synthetic tests, build instructions and release artifacts belong
here. Live configuration, OAuth files, routing checkpoints, session transcripts
and operational logs stay outside the repository and release assets.

## Build the Windows gateway

The deployed executable was built on Linux using Go 1.27.1 and Debian Bookworm's
`gcc-mingw-w64-x86-64-posix`. CGO must stay enabled to retain plugin support.
Use a clean checkout of the release tag, then:

```sh
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-posix GOOS=windows GOARCH=amd64 \
  go build -p 2 \
  -ldflags '-s -w -X main.Version=8.0.13-se-affinity.3 -X main.Commit=79919de' \
  -o cli-proxy-affinity-v3.exe ./cmd/server
sha256sum cli-proxy-affinity-v3.exe
```

The third release's executable SHA256 is
`d20c7472528abe3664b3b715840778d48533f4f302ba310e9d2978c44b8dbf46`.
The memory-only .2 fallback has SHA256
`833ad767269322d902206d654e709c5797120f96ef7d3f4f226a44a7d668c99f`.
This identifies the tested artifact; a different compiler or build environment
may produce different bytes. A successful rebuild is not deployment validation.

## Upgrade and rollback

1. Fetch the specific upstream release into a candidate branch based on `se/main`.
   Inspect routing, auth refresh, persistence and streaming changes. Preserve
   published history; never rewrite a deployed release tag.
2. Merge the release into the candidate. Reconcile our patch and review the diff.
   Do not copy a new official executable over the custom gateway.
3. Run the tests in AFFINITY-PERSISTENCE.md, compile the server, and run the durable
   affinity tests on Windows. Synthetic tests require no provider credentials.
4. Record the source commit, compiler, build command and executable checksum.
   Retain the previous executable and service configuration.
5. After active turns finish, test actual Claude Code and Codex workers: useful
   work, follow-up tools, native conversation continuity, account attribution,
   and retention after a service restart with alternative accounts available.
6. Promote only after those checks pass. If they fail, restore the last accepted
   executable and its reviewed compatible service configuration, preserving current
   auth and conversation state. A future candidate falling back to .3 may retain
   persistence only when checkpoint compatibility is verified. A rollback to .1/.2
   must disable persistence. Back up routing state before any schema change.

Review upstream weekly for relevant fixes, and sooner when a provider integration
breaks. Existing healthy conversations keep their account ahead of any future
reset-time preference. A provider cache hit is separate from an affinity hit.

## Maintenance check

The operator maintaining the personal fleet owns this check. It is manual: this
document does not create a scheduled job or imply that somebody is watching it.
Run it weekly, before a gateway upgrade, and when provider behavior changes:

```sh
git status --short
git fetch --quiet origin
git fetch --quiet --no-tags upstream main
gh release list -R router-for-me/CLIProxyAPI --limit 5
git diff --stat v8.0.13 origin/se/main
git diff --name-only v8.0.13-se-affinity.3 origin/se/main
```

Choose a published upstream tag after reviewing its release notes. Fetch that tag
explicitly; inspect the endpoint-to-endpoint diff, then its ancestry-path log. This
checkout began shallow, so an unrestricted history range can include unrelated old
side-branch history. Do not interpret that as thousands of new changes.

```sh
git fetch --quiet --no-tags upstream tag v8.0.16
git diff --stat v8.0.13 v8.0.16
git log --ancestry-path --no-merges --oneline v8.0.13..v8.0.16
git merge-tree --write-tree origin/se/main v8.0.16
```

These example tags identify the 6 October review; choose and record the tags for
each later check. `merge-tree` checks textual integration without switching the
checkout or live service. A clean result is not runtime or semantic validation.

Before accepting a candidate, rerun the full Go suite, the focused race tests in
AFFINITY-PERSISTENCE.md, and a server build. Run the Windows-specific durability
tests with the release compiler/CGO settings, including startup under a blocked
checkpoint replacement. Measure `BenchmarkDurableSelection` on the same PC and
compare with the accepted release under similar load. Then perform the bounded
real-worker checks above. Tests use temporary synthetic state, never live auth or
the live affinity checkpoint.

The 6 October candidate passed the full Go suite, focused affinity/auth-persistence
race tests, a Windows CGO server/test build and all ten Windows durability tests.
An unmodified upstream checkout reproduces a catalog-test cleanup failure: it tries
to restore the Devin catalog with Home mode enabled, which now disables that catalog.
The candidate corrects that test-only flag. Three-repeat isolated and race checks
and the full suite pass after the correction. No production behavior was changed
to silence the test.

The candidate Windows artifact is built from merge commit a4aa7fb5, with Go1.27.1,
MinGW-posix GCC12 and CGO enabled. Its version is8.0.16-se-affinity.4-rc1; SHA256 is
4209a2dbbaf59fd60050c814eb259a8d2ad66cc1e9d4ff21bfca65249b07aeef.
The later test-cleanup commit changes no production source. This is a candidate,
not an accepted service release or a new deployed tag.

On the same Windows host, three runs of100 selections gave median serial durable
selection times of3.463ms at2 bindings,2.149ms at32, and2.357ms at256. The accepted
release measured3.253/2.274/2.410ms in the preceding check. Concurrent cases passed
without account drift. These are short selector/checkpoint samples, not throughput
or whole-request latency guarantees.

Actual Claude/Codex canaries and a completed-turn gateway restart are still required
before release acceptance. Upstream's tool-alias store for shortened Claude
continuations is process-local and distinct from our durable account bindings.
Its missing-state response requests full-history replay; the real installed client
must demonstrate its response to that condition where applicable. A saved account
binding alone does not establish continuation compatibility. Credential persistence
also changed upstream; rollback must preserve the newest auth files, not restore
an old credential snapshot.

An update record must name the deployed tag, candidate tag, reviewed relevant
changes, test results, artifact checksum, rollback compatibility and decision.
Keep machine names, account labels, raw test logs and fleet task state in private
operations records. This public repository is for source and generic procedures.

GitHub Actions are disabled for this fork initially. Inherited release and
container-publishing workflows are not our release procedure. Releases are
manual, and no workflow deploys to the live gateway. Enable automation only
after adapting it to this documented scope.

The releases have unit/race tests, Windows filesystem/process tests and
real completed-turn restart evidence. It does not establish exactly-once
in-flight tool execution, every power-loss outcome, or large-fleet performance.
