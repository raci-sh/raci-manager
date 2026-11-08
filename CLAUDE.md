# CLAUDE.md

Guidance for AI coding agents working in this repository.

## What this repo is

`raci-manager` is the RACI.sh engine: a Go library that the `raci` CLI ([raci-sh/raci-CLI](https://github.com/raci-sh/raci-CLI)) imports as a module, and that the future cloud service (`raci-cloud`) will import too. Validation, hashing, the audit log, the `.raci` view and verification live here; consumers only present results.

The normative source is the specification in [raci-sh/raci-spec](https://github.com/raci-sh/raci-spec) (`SPEC.md`, and the `sh.raci` realm in `rules/`). Package docs cite its sections. When code and spec disagree, the spec wins; file a bead in raci-spec if the spec is wrong.

Everything must work offline, with just the CLI, local files and GitHub. Never add a dependency on a cloud service.

<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:1105d646 -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files

**Architecture in one line:** issues live in a local Dolt DB; sync uses `refs/dolt/data` on your git remote; `.beads/issues.jsonl` is a passive export. See https://github.com/gastownhall/beads/blob/main/docs/core-concepts/sync-concepts.md for details and anti-patterns.

## Agent Context Profiles

The managed Beads block is task-tracking guidance, not permission to override repository, user, or orchestrator instructions.

- **Conservative (default)**: Use `bd` for task tracking. Do not run git commits, git pushes, or Dolt remote sync unless explicitly asked. At handoff, report changed files, validation, and suggested next commands.
- **Minimal**: Keep tool instruction files as pointers to `bd prime`; use the same conservative git policy unless active instructions say otherwise.
- **Team-maintainer**: Only when the repository explicitly opts in, agents may close beads, run quality gates, commit, and push as part of session close. A current "do not commit" or "do not push" instruction still wins.

## Session Completion

This protocol applies when ending a Beads implementation workflow. It is subordinate to explicit user, repository, and orchestrator instructions.

1. **File issues for remaining work** - Create beads for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **Handle git/sync by active profile**:
   ```bash
   # Conservative/minimal/default: report status and proposed commands; wait for approval.
   git status

   # Team-maintainer opt-in only, unless current instructions forbid it:
   git pull --rebase
   git push
   git status
   ```
5. **Hand off** - Summarize changes, validation, issue status, and any blocked sync/commit/push step

**Critical rules:**
- Explicit user or orchestrator instructions override this Beads block.
- Do not commit or push without clear authority from the active profile or the current user request.
- If a required sync or push is blocked, stop and report the exact command and error.
<!-- END BEADS INTEGRATION -->

## Build & Test

```sh
make build   # ./raci-manager (version only; the library is the product)
make test    # go test -race ./...
make cover   # tests with the per-package coverage floor (scripts/coverage-gate.sh, 70%)
make lint    # golangci-lint v2
make fmt
```

## Architecture

All Go code is under `src/`. Public packages are imported as `github.com/raci-sh/raci-manager/src/<pkg>`.

| Package | SPEC.md | Role |
| --- | --- | --- |
| `model` | §1 | Roles (`"AR"`, canonical A>R>C>I order), parties, activities, resources, entries |
| `canonical` | §4.1 | RFC 8785 canonical JSON (via `gowebpki/jcs`) |
| `digest` | §4.1 | `blake3:` + 64 hex digests, the all-zero sentinel |
| `event` | §4.1 | Event envelope, payload types, `Seal` / `ComputeHash` |
| `eventlog` | §4.2 | Per-realm JSONL log: `Append`, `Read`, `VerifyChain` |
| `view` | §7 | The `.raci` materialised view and its canonical ordering; `Replay` (stub) |
| `constraint` | §2 | Constraint names, the strict preset, conflict pairs |
| `inference` | §3.1 | Inference rule vocabulary and structural validation |
| `verify` | §7.1 | Finding codes, the verification report |
| `internal/version` | | Build info injected via ldflags |
| `cmd/raci-manager` | | Tiny binary: `version` |

Rules enforced by `.golangci.yml`:

- Library packages never print, log or exit (`forbidigo`); they return errors or `verify.Finding` values. Only `src/cmd/` and tests may.
- The manager never imports `raci-CLI` or `raci-cloud` (`depguard`). Dependencies point from consumers to the manager.

## Conventions

- Conventional Commits. `.beads/hooks/commit-msg` checks locally; `lint-commits` and `pr-title` check in CI.
- Tests sit beside the code (`*_test.go`, `testdata/`). Every package with statements must reach the coverage floor.
- Apache 2.0. Copyright: Consulare and contributors (`LICENSE`, `NOTICE`).

## Branching and releases

- The default branch is **`develop`**. Feature branch → PR into `develop`, **squash only** (the PR title becomes the commit, so it must be a conventional commit).
- `develop` → `main` happens only through the **Promote to production** workflow (`promote.yml`, run from `main`: `gh workflow run promote.yml --ref main`), which merges with a **merge commit only**. It then merges release-please's release PR, then the sync-back PR.
- release-please (`release-please.yml`) runs on `main`, tags `vX.Y.Z` and writes `CHANGELOG.md`; the tag runs goreleaser (`release.yml`). `sync-release-files.yml` carries `CHANGELOG.md` and `.release-please-manifest.json` back to `develop`. Never hand-edit either file.
- Rulesets `merge-method-develop` (squash) and `merge-method-main` (merge) enforce the merge methods. Release workflows use the `raci-sh-app` GitHub App (`RACI_APP_ID`, `RACI_APP_PRIVATE_KEY`).
- After a manager release, bump the CLI: `go get github.com/raci-sh/raci-manager@vX.Y.Z` in raci-CLI.
- Agents open PRs; they never merge them.
