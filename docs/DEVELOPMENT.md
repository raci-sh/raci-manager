# Development

Build, test and release workflow of [raci-manager](../README.md).

## Commands

```bash
make build   # build the raci-manager binary (version only)
make test    # run the tests with the race detector
make cover   # run the tests with the per-package coverage floor
make lint    # run golangci-lint v2
make fmt     # format the code
```

CI runs `make cover`, which fails any package under 70% statement coverage
([`coverage-gate.sh`](../scripts/coverage-gate.sh)). It also builds, vets, tests on Windows, runs
golangci-lint, actionlint and editorconfig-checker, checks commit messages against
[`.commitlintrc.yml`](../.commitlintrc.yml), and checks that the pull request title is a
Conventional Commit.

Two lint rules shape the code. Library packages never print, log or exit; they return errors or
`verify.Finding` values, and only `src/cmd/` and tests may print. The manager never imports
raci-CLI or raci-cloud. See [`AGENTS.md`](../AGENTS.md) for the agent workflow and the
[beads](https://github.com/gastownhall/beads) issue tracker.

## Branching and releases

`develop` is the default branch and collects the changes for the next release;
`main` only moves when a release is cut:

1. Pull requests target `develop` and are **squashed**. The pull request title
	becomes the commit release-please reads, so it decides the release note and
	the version bump; the `pr-title` check enforces Conventional Commits.
2. The **Promote to production** workflow carries the bundle to `main`. Run it
	from `main`:

	```bash
	gh workflow run promote.yml --ref main
	```

	It opens a `chore: promote develop to main` pull request and merges it with a
	**merge commit**, so release-please reads the individual conventional commits
	rather than a single `chore:`.
3. release-please opens its release pull request against `main`; the workflow
	merges it, which tags `vX.Y.Z` and writes `CHANGELOG.md`. The tag runs
	GoReleaser.
4. **Sync release files to develop** opens a pull request carrying
	`CHANGELOG.md` and `.release-please-manifest.json` back to `develop`, and the
	promote workflow merges it. Only those two paths are synced; never edit
	either file by hand.

> [!IMPORTANT]
> Promote `develop` before merging a release pull request by hand. release-please
> stops reading at the last release commit, by date, so a `feat:` or `fix:` left
> on `develop` while a release merges is dropped from the next release.

The `Release changelog` workflow catches a dropped commit. After every Release
Please run, and on each release pull request, it lists the commits `main`
gained since the last release by ancestry and checks that `CHANGELOG.md` names
each one. Its verdict is the `release-changelog` status, which is required on
`main`.

The `merge-method-develop` (squash) and `merge-method-main` (merge) rulesets
enforce both merge methods. The release workflows act as the `raci-sh-app`
GitHub App, so their commits are signed and their pull requests run CI.

After a manager release, bump the CLI:

```bash
go get github.com/raci-sh/raci-manager@vX.Y.Z   # in raci-CLI
```
