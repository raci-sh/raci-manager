<div align="center">

<img src="https://raci.sh/assets/logo.svg" alt="RACI.sh logo" width="96" height="96">

# RACI.sh Manager

**The engine of RACI.sh: a Go library that records who is Responsible, Accountable, Consulted and Informed for activities on resources, per realm, in a hash-chained audit log.**

[![CI](https://github.com/raci-sh/raci-manager/actions/workflows/ci.yml/badge.svg)](https://github.com/raci-sh/raci-manager/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/raci-sh/raci-manager?color=161B22&labelColor=161B22)](https://github.com/raci-sh/raci-manager/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/raci-sh/raci-manager.svg)](https://pkg.go.dev/github.com/raci-sh/raci-manager)
[![Go](https://img.shields.io/github/go-mod/go-version/raci-sh/raci-manager?color=161B22&labelColor=161B22)](go.mod)
[![License](https://img.shields.io/badge/License-Apache%202.0-161B22?labelColor=161B22)](LICENSE)

[Install](#install) · [Where to go](#where-to-go) · [Governance](#governance) · [Security](#security) · [Community](#community)

</div>

[RACI.sh](https://raci.sh) is RACI for agents. The manager holds the RACI model, the hash-chained
audit log, canonical JSON and digests, the `.raci` materialised view, constraints, inference rules
and verification. It depends on no cloud service: everything works offline with the CLI, local
files and GitHub. Every package under `src/` that has statements keeps at least 70% statement
coverage, which CI enforces with [`scripts/coverage-gate.sh`](scripts/coverage-gate.sh).

> [!NOTE]
> Status: scaffold. The types and the log hash chain work. Replay, inference, constraint
> evaluation and the rest of verification are still to be built.

## Install

```bash
go get github.com/raci-sh/raci-manager@latest
```

Verify the hash chain of a realm's audit log:

```go
package main

import (
	"fmt"
	"os"

	"github.com/raci-sh/raci-manager/src/eventlog"
	"github.com/raci-sh/raci-manager/src/verify"
)

func main() {
	report, err := verify.Log("com.example", eventlog.Open(eventlog.FileName))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	for _, f := range report.Findings {
		fmt.Printf("%s %s: %s\n", f.Severity, f.Code, f.Detail)
	}
	if report.HasErrors() {
		os.Exit(1)
	}
}
```

Most users want [`raci`](https://github.com/raci-sh/raci-CLI), the command line built on this
library. Releases also ship a `raci-manager` binary with a `version` command.

## Where to go

| I want to | Go to |
| --------- | ----- |
| Record and check RACI from a terminal | [raci-CLI](https://github.com/raci-sh/raci-CLI) |
| See how the manager fits RACI.sh | [Architecture](docs/REFERENCE.md#architecture) |
| Find the package for a task | [Packages](docs/REFERENCE.md#packages) |
| Read the specification | [raci-spec `SPEC.md`](https://github.com/raci-sh/raci-spec/blob/main/SPEC.md) |
| Read the Go API documentation | [pkg.go.dev](https://pkg.go.dev/github.com/raci-sh/raci-manager) |
| Build, test and lint | [Development](docs/DEVELOPMENT.md#commands) |
| Understand branching and releases | [Branching and releases](docs/DEVELOPMENT.md#branching-and-releases) |

## Governance

Open source software by Consulare. RACI.sh is a companion to [SteerSpec](https://github.com/SteerSpec),
and its realms match SteerSpec realms, but it does not require SteerSpec.

Apache 2.0, see [LICENSE](LICENSE) and [NOTICE](NOTICE). Pull requests target `develop` and use
[Conventional Commits](https://www.conventionalcommits.org/). The workflow is in
[Branching and releases](docs/DEVELOPMENT.md#branching-and-releases).

## Security

Do not open a public issue for a vulnerability. Report it privately through GitHub: on this
repository, choose **Security**, then **Report a vulnerability**.

## Community

Questions and bug reports go to the
[issue tracker](https://github.com/raci-sh/raci-manager/issues).
