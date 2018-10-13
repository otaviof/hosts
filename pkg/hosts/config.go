package hosts

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config primary configuration file contents
type Config struct {
	Hosts     Hosts     `yaml:"hosts"`
	Blacklist Blacklist `yaml:"blacklist"`
}

// Hosts `hosts` configuration block
type Hosts struct {
	BaseDirectory string `yaml:"baseDirectory"`
	Output        string `yaml:"output"`
}

// Blacklist `blacklist` configuration block
type Blacklist struct {
	Output         string         `yaml:"output"`
	AddressMapping AddressMapping `yaml:"addressMapping"`
}

// AddressMapping `addressMapping` block inside `blacklist`
type AddressMapping struct {
	Search  string `yaml:"search"`
	Replace string `yaml:"replace"`
}

// Validate check if all required configuration fields are present.
func (c *Config) Validate() error {
	if !isDir(c.Hosts.BaseDirectory) {
		return fmt.Errorf("Can't find directory at: '%s'", c.Hosts.BaseDirectory)
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
