package hosts

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config primary configuration file contents
type Config struct {
	Hosts    Hosts      `yaml:"hosts"`
	External []External `yaml:"external"`
}

// Hosts `hosts` configuration block
type Hosts struct {
	BaseDirectory string `yaml:"baseDirectory"`
	Output        string `yaml:"output"`
}

// External `blacklist` configuration block
type External struct {
	URL      string    `yaml:"url"`
	Output   string    `yaml:"output"`
	Mappings []Mapping `yaml:"mappings"`
	Skip     []string  `yaml:"skip"`
}

// Mapping `mappings` block inside `blacklist`
type Mapping struct {
	Search  string `yaml:"search"`
	Replace string `yaml:"replace"`
}

// Validate check if all required configuration fields are present.
func (c *Config) Validate() error {
	if !isDir(c.Hosts.BaseDirectory) {
		return fmt.Errorf("Can't find directory at: '%s'", c.Hosts.BaseDirectory)
	}
	if c.Hosts.Output == "" {
		return fmt.Errorf("output is mandatory: '%s'", c.Hosts.Output)
	}
	return nil
}

// NewConfig creates a new configuration instance baed on yaml input.
func NewConfig(path string) (*Config, error) {
	var config = &Config{}
	var err error

	if err = yaml.Unmarshal(readFile(path), config); err != nil {
		return nil, err
	}

	// expanding base directory
	config.Hosts.BaseDirectory = strings.Replace(
		config.Hosts.BaseDirectory, "~", os.Getenv("HOME"), 1)

	return config, nil
}
