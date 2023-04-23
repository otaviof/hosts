package main

import (
	"log"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generate and write the final '/etc/hosts' file.",
	Long: `
By inspecting 'hosts.baseDirectory', 'apply' command will compose a unified version to be stored
at '/etc/hosts' file, or another configurable location. Every file name found in 'baseDirectory' is
kept as a comment before it's contents.

On reading the input files, it's applied sanitization to only bring strings that are starting with an
IPv4 or IPv6 address, on IPv6 format we also include OS X interface name aliasing, using '%'
notation.

Please not you may also need the companion of 'sudo' to execute 'hosts apply' action.
`,
	Run: runApplyCmd,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

// runApplyCmd instantiate Hosts and run Apply routine.
func runApplyCmd(_ *cobra.Command, _ []string) {
	hosts := newHosts()
	if err := hosts.Load(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	if err := hosts.Apply(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
