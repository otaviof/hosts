package main

import (
	"log"

	"github.com/spf13/cobra"

	hosts "github.com/otaviof/hosts/pkg/hosts"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update host entries from 'external' resource.",
	Long: `
Update will execute a GET request agains configured 'external.URL', and read returned body line by
line to extract host entries. Those entires can be use 'mapping's, in order to search and replace
string. You can also configure 'skip', which corresponds to a list of regular-expressions for lines
to skip from payload to be stored.

After 'update' command you also need to execute 'apply' command to consist those changes in the
final '/etc/hosts' file.`,
	Run: runUpdateCmd,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdateCmd(cmd *cobra.Command, args []string) {
	var config = getConfig()
	var update *hosts.Update
	var err error

	if update, err = hosts.NewUpdate(config, dryRun); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	if err = update.Execute(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
