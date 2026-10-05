# Super Employee maintenance

This fork carries a small, opt-in durability patch on CLIProxyAPI v8.0.13.
Upstream remains https://github.com/router-for-me/CLIProxyAPI. Preserve its MIT
license and attribution. Keep provider integrations and ordinary fixes upstream.

## Source and releases

- `se/main` is our maintained branch. The inherited `main` is upstream history,
  not our deployment branch.
- `v8.0.13-se-affinity.2` identifies the deployed source commit
  `051b65101bfbbc3df4a9c401599b5f5f971c7999`.
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
  -ldflags '-s -w -X main.Version=8.0.13-se-affinity.2 -X main.Commit=051b651' \
  -o cli-proxy-affinity-v2.exe ./cmd/server
sha256sum cli-proxy-affinity-v2.exe
```

The second release's executable SHA256 is
`833ad767269322d902206d654e709c5797120f96ef7d3f4f226a44a7d668c99f`.
The previous release's hash is
`62ffc6129e9c882c4e1cb1ae9490078a9efcf6bf338a1cd19b24c967736125cb`.
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
6. Promote only after those checks pass. If they fail, restore the previous
   executable/configuration while preserving the current auth store. Back up
   routing state before any schema change and verify rollback compatibility.

Review upstream weekly for relevant fixes, and sooner when a provider integration
breaks. Existing healthy conversations keep their account ahead of any future
reset-time preference. A provider cache hit is separate from an affinity hit.

GitHub Actions are disabled for this fork initially. Inherited release and
container-publishing workflows are not our release procedure. Releases are
manual, and no workflow deploys to the live gateway. Enable automation only
after adapting it to this documented scope.

The releases have unit/race tests, Windows filesystem/process tests and
real completed-turn restart evidence. It does not establish exactly-once
in-flight tool execution, every power-loss outcome, or large-fleet performance.
