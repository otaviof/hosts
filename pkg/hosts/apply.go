package hosts

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// Apply constructs a new '/etc/hosts' file.
type Apply struct {
	config   *Config  // app configuration
	files    []string // list of '.host' files
	etcHosts []string // final contents of '/etc/hosts'
	dryRun   bool     // dry-run mode
}

// print shows 'etcHosts' contents in stdout.
func (a *Apply) print() {
	for _, line := range a.etcHosts {
		fmt.Printf("%s\n", line)
	}
}

// render save 'etcHosts' conents in configured output file.
func (a *Apply) render() error {
	var content []byte
	var line string

	log.Printf("Saving contents at: '%s'", a.config.Hosts.Output)

	for _, line = range a.etcHosts {
		for _, b := range []byte(fmt.Sprintf("%s\n", line)) {
			content = append(content, b)
		}
	}

	return ioutil.WriteFile(a.config.Hosts.Output, content, 0644)
}

// Execute steps to load and render/print output file.
func (a *Apply) Execute() error {
	var file string
	var parser *Parser
	var err error

	log.Printf("Files: '%#v'", a.files)

	for _, file = range a.files {
		var line string

		if parser, err = NewParser(file); err != nil {
			return err
		}
		if err = parser.Ingest(); err != nil {
			return err
		}

		a.etcHosts = append(a.etcHosts, fmt.Sprintf("### %s", file))
		for _, line = range parser.Contents {
			a.etcHosts = append(a.etcHosts, line)
		}
	}

	if a.dryRun {
		a.print()
	} else {
		if err = a.render(); err != nil {
			return err
		}
	}

	return nil
}

// dirGlob search for '.host' files in configured location.
func dirGlob(dirPath string) ([]string, error) {
	var files []string
	var err error

	if files, err = filepath.Glob(filepath.Join(dirPath, "*.host")); err != nil {
		return nil, err
	}

	return files, nil
}

// NewApply create a new Apply object, by inspecting base directory.
func NewApply(config *Config, dryRun bool) (*Apply, error) {
	var apply = &Apply{config: config, dryRun: dryRun}
	var err error

	log.Printf("Inspecting base-directory: '%s'", config.Hosts.BaseDirectory)

	if apply.files, err = dirGlob(config.Hosts.BaseDirectory); err != nil {
		return nil, err
	}

	return apply, nil
}
