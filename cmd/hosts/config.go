package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Validate configuration file.",
	Long: `
Config sub-command aims to validate "hosts.yaml" configuration file without modifying or loading data
during this process. The location used is informed by "--base-dir".

All attributes in "hosts.yaml" are documented at: https://github.com/otaviof/hosts#configuration
	`,
	Run: runConfigCmd,
}

// validate flag to indicate if validation should take place.
var validate bool

func init() {
	flags := configCmd.PersistentFlags()

	flags.BoolVar(&validate, "validate", false, "validate configuration contents")

	rootCmd.AddCommand(configCmd)
}

// runConfigCmd instantiate a new configuration, which validates the base directory and existence of
// "hosts.yaml", and then excute the configuration validation.
func runConfigCmd(_ *cobra.Command, _ []string) {
	cfg, err := newConfig()
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	if validate {
		if err = cfg.Validate(); err != nil {
			log.Fatalf("[ERROR] Invalid configuration: %v", err)
		}
		log.Infof("Configuration file is valid!")
	}
}
