package hosts

import (
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	configDir  = ".hosts"
	ConfigFile = "hosts.yaml"
	extension  = "host"
)

// Root configuration top level object.
type Root struct {
	Hosts Config `json:"hosts"`
}

// Config primary application configuration.
type Config struct {
	Input  Input    `json:"input"`
	Output []Output `json:"output"`
}

// Input input sectin, for data obtained externally.
type Input struct {
	Sources         []Source         `json:"sources"`
	Transformations []Transformation `json:"transformations"`
}

// Source input source, describes a single URI.
type Source struct {
	Name string `json:"name,omitempty"`
	URI  string `json:"uri"`
	File string `json:"file"`
}

// Transformation describes how data obtained externally will be transformed.
type Transformation struct {
	Name    string `json:"name,omitempty"`
	Search  string `json:"search"`
	Replace string `json:"replace,omitempty"`
}

// Output describes a output file.
type Output struct {
	Name    string `json:"name,omitempty"`
	Path    string `json:"path"`
	Dnsmasq bool   `json:"dnsmasq,omitempty"`
	With    string `json:"with,omitempty"`
	Without string `json:"without,omitempty"`
}

func (o *Output) CompileREs() (*regexp.Regexp, *regexp.Regexp, error) {
	var withRE *regexp.Regexp = nil
	var withoutRE *regexp.Regexp = nil
	var err error

	if o.With != "" {
		if withRE, err = regexp.Compile(o.With); err != nil {
			return nil, nil, err
		}
	}
	if o.Without != "" {
		if withoutRE, err = regexp.Compile(o.Without); err != nil {
			return nil, nil, err
		}
	}
	return withRE, withoutRE, nil
}

// TODO: validate instantiated configuration
func (c *Config) Validate() error {
	return nil
}

// NewConfig creates a new configuration instance based informed file path.
func NewConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	root := &Root{}
	if err = yaml.Unmarshal(data, root); err != nil {
		return nil, err
	}
	return &root.Hosts, nil
}
