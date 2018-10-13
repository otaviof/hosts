package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/otaviof/hosts/pkg/hosts"
)

var rootCmd = &cobra.Command{
	Use: "hosts",
	Long: `### hosts

Helps you to manage and dinamically render '/etc/hosts' file contents, based in small ".host" files.
You can also fetch hosts definitions from the internet, parse, clean-up, and let it be part of your
hosts definitions.`,
}

var configPath string
var dryRun bool

func init() {
	var flags = rootCmd.PersistentFlags()

	flags.StringVar(&configPath, "config", "/usr/local/etc/hosts.yaml", "configuration file path")
	flags.BoolVar(&dryRun, "dry-run", false, "dry-run mode")
}

func getConfig() *hosts.Config {
	var config *hosts.Config
	var err error

	log.Printf("Loading configuration from: '%s'", configPath)
	if config, err = hosts.NewConfig(configPath); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	log.Printf("Validating configuraiton...")
	if err = config.Validate(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	return config
}

func main() {
	var err error

	if err = rootCmd.Execute(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
