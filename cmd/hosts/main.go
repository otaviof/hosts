package main

import (
	"path"

	hosts "github.com/otaviof/hosts/pkg/hosts"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use: "hosts",
	Long: `### hosts

Helps you to manage and dinamically render '/etc/hosts' file contents, based in small ".host" files.
You can also fetch hosts definitions from the internet, parse, clean-up, and let it be part of your
hosts definitions.

You can generate multiple output files, and specify dnsmasq format as well. Therefore, from a single
source-of-authority you can include and exclude parts to generate those files.
`,
}

var (
	logLevel int
	baseDir  string
	dryRun   bool
)

func init() {
	flags := rootCmd.PersistentFlags()

	flags.IntVar(&logLevel, "log-level", int(log.InfoLevel), "configuration file path")
	flags.StringVar(&baseDir, "base-dir", "", "configuration file path")
	flags.BoolVar(&dryRun, "dry-run", false, "dry-run mode")
}

// newHosts instantiate the application and configuration.
func newHosts() *hosts.Hosts {
	hosts.SetLogLevel(logLevel)

	if baseDir == "" {
		var err error
		if baseDir, err = hosts.DefaultConfigDir(); err != nil {
			log.Fatalf("error finding default base-dir: %s", err)
		}
	}
	configPath := path.Join(baseDir, hosts.ConfigFile)
	log.Debugf("Using configuration file at '%s'", configPath)

	cfg, err := hosts.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error instantiating config: '%v'", err)
	}
	if err = cfg.Validate(); err != nil {
		log.Fatalf("error validating configuration: '%v'", err)
	}

	return hosts.NewHosts(cfg, baseDir, dryRun)
}

func main() {
	var err error

	if err = rootCmd.Execute(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
