package main

import (
	"log"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update host entries from 'external' resource.",
	Long: `
Update will execute a GET request against configured 'URI', and read returned body line by
line, those lines will be transformed accordingly with 'external[n].transform' block, where lines
can be subject to search-and-replace, or skipped, before data is stored.

After 'update' command you also need to execute 'apply' command to consist those changes in the
final '/etc/hosts' file, or other configured location.`,
	Run: runUpdateCmd,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdateCmd(cmd *cobra.Command, args []string) {
	hosts := newHosts()
	if err := hosts.Update(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
