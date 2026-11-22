Use conventional commit format for all commit messages.

Format: `<type>(<scope>): <description>`

Types: feat, fix, docs, chore, refactor, test, ci, perf, build

Scopes are package names under `src/` (`model`, `canonical`, `digest`, `event`, `eventlog`,
`view`, `constraint`, `inference`, `verify`, `cmd`), or `deps`, `ci`, `readme`.

Examples:
- `fix(eventlog): reject a line that breaks the hash chain`
- `feat(view): replay entry events into the .raci view`
- `docs(readme): add the package reference`

Do not use generic messages like "Update file.go". Always include a type prefix.
