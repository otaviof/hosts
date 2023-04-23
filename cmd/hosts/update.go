package main

import (
	"log"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update host entries from 'external' resource.",
	Long: `
Update will execute a GET request against configured 'URI', and read returned body line by line,
those lines will be transformed as 'hosts.input.transformations' block describes, where lines
can be subject to search-and-replace, or skipped.

After 'update' command you also need to execute 'apply' to consist those changes in the output files.
`,
	Run: runUpdateCmd,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// runUpdateCmd instantiate Hosts, and run update routine.
func runUpdateCmd(_ *cobra.Command, _ []string) {
	hosts := newHosts()
	if err := hosts.Update(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
