package main

import (
	"log"

	"github.com/spf13/cobra"

	hosts "github.com/otaviof/hosts/pkg/hosts"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update blacklist entries.",
	Long:  ``,
	Run:   runUpdateCmd,
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
