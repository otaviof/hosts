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

You can generate multiple output files, and specify the dnsmasq format as well. Therefore, from a
single source-of-authority you can include and exclude parts to generate those files.
`,
}

var (
	// logLevel flag to hold log verbosity level.
	logLevel int
	// baseDir flag to hold path to base directory.
	baseDir string
	// dryRun flag to indicate dry-run mode.
	dryRun bool
)

func init() {
	flags := rootCmd.PersistentFlags()

	flags.IntVar(&logLevel, "log-level", int(log.InfoLevel), "log verbosity, from 0 to 6")
	flags.StringVar(&baseDir, "base-dir", "", "configuration file path")
	flags.BoolVar(&dryRun, "dry-run", false, "dry-run mode")
}

// newConfig returns a hosts.Config instance based on parameters informed.
func newConfig() (*hosts.Config, error) {
	hosts.SetLogLevel(logLevel)

	if baseDir == "" {
		var err error
		if baseDir, err = hosts.DefaultConfigDir(); err != nil {
			log.Fatalf("error finding default base-dir: %s", err)
		}
	}
	configPath := path.Join(baseDir, hosts.ConfigFile)
	log.Debugf("Using configuration file at '%s'", configPath)

	return hosts.NewConfig(configPath)
}

// newHosts instantiate the application and configuration.
func newHosts() *hosts.Hosts {
	cfg, err := newConfig()
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
