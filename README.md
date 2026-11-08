# raci-manager

The engine of [RACI.sh](https://raci.sh): RACI for agents. It records who is **R**esponsible, **A**ccountable, **C**onsulted and **I**nformed for activities on resources, per realm.

`raci-manager` is a Go library. It holds the RACI model, the hash-chained audit log, canonical JSON and digests, the `.raci` materialised view, constraints, inference rules and verification. The [`raci` CLI](https://github.com/raci-sh/raci-CLI) imports it, and so will the cloud service later. The library has no dependency on any cloud, so everything works offline with just the CLI, local files and GitHub.

The specification is in [raci-sh/raci-spec](https://github.com/raci-sh/raci-spec). RACI.sh is a companion to [SteerSpec](https://github.com/SteerSpec), but you can use it without SteerSpec.

> Status: scaffold. The types and the log hash chain work. Replay, inference, constraint evaluation and the rest of verification are still to be built.

## Install

```sh
go get github.com/raci-sh/raci-manager
```

```go
import (
	"github.com/raci-sh/raci-manager/src/eventlog"
	"github.com/raci-sh/raci-manager/src/verify"
)

report, err := verify.Log("sh.raci", eventlog.Open(".raci.jsonl"))
```

## Packages

| Package | What it does |
| --- | --- |
| `src/model` | Roles, parties, activities, resources, entries |
| `src/canonical` | RFC 8785 canonical JSON |
| `src/digest` | `blake3:<hex>` digests |
| `src/event` | Audit event envelope and payloads |
| `src/eventlog` | Per-realm JSONL log with hash-chain verification |
| `src/view` | The `.raci` materialised view |
| `src/constraint` | Strictness constraints and the strict preset |
| `src/inference` | Inference rule vocabulary (map, add-issuer, strip) |
| `src/verify` | Verification findings and report |

## Development

```sh
make test cover lint build
```

`CLAUDE.md` covers the branching model: `develop` is the default branch and takes squash merges. Releases reach `main` through the *Promote to production* workflow, as merge commits, and release-please tags them.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
