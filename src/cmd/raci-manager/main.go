// Command raci-manager prints the build information of the RACI.sh manager
// library. The user-facing tool is the raci CLI (raci-sh/raci-CLI).
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/raci-sh/raci-manager/src/internal/version"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "raci-manager",
		Short:         "RACI.sh manager library",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "raci-manager %s (commit %s, built %s)\n", version.Version, version.Commit, version.Date)
		},
	})
	return root
}

// run executes the command line and returns the process exit code.
func run(args []string, stdout, stderr io.Writer) int {
	root := newRootCmd()
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(stderr, "Error:", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
