package main

import (
	"log"

	hosts "github.com/otaviof/hosts/pkg/hosts"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Generate a `/etc/hots` formated file in configured location.",
	Long:  `Gather all '.host' files and compose a unified version of it.`,
	Run:   runApplyCmd,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func runApplyCmd(cmd *cobra.Command, args []string) {
	var config = getConfig()
	var apply *hosts.Apply
	var err error

	if apply, err = hosts.NewApply(config, dryRun); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	if err = apply.Execute(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
