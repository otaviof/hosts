package hosts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

// ConfigFile default configuration file name.
const ConfigFile = "hosts.yaml"

const (
	// configDir default configuration directory.
	configDir = ".hosts"
	// extension default extension name.
	extension = "host"
)

var (
	ErrRequiredAttribute = errors.New("required attribute not informed")
	ErrInvalidRegex      = errors.New("invalid regular-expresssion")
)

// Root configuration top level object.
type Root struct {
	Hosts Config `json:"hosts"` // root
}

// Config primary application configuration.
type Config struct {
	Input  Input    `json:"input"`  // input block, data coming from external sources
	Output []Output `json:"output"` // output block, files that will be created
}

// Input input section, for data obtained externally.
type Input struct {
	Sources         []Source         `json:"sources"`         // slice of sources
	Transformations []Transformation `json:"transformations"` // slice of transformations
}

// Source input source, describes a single URI.
type Source struct {
	Name string `json:"name,omitempty"` // data source name
	URL  string `json:"url"`            // resource url
	File string `json:"file"`           // file name to store collected data
}

// Transformation describes how data obtained externally will be transformed.
type Transformation struct {
	Search  string `json:"search"`            // search regular expresssion
	Replace string `json:"replace,omitempty"` // replace with
}

// CompileRE compiles the regular expression in search attribute.
func (t *Transformation) CompileRE() (*regexp.Regexp, error) {
	if t.Search == "" {
		return nil, nil
	}
	return regexp.Compile(t.Search)
}

// Output describes a output file.
type Output struct {
	Name    string `json:"name,omitempty"`    // output name
	Path    string `json:"path"`              // output file full path
	Dnsmasq bool   `json:"dnsmasq,omitempty"` // format as dnsmasq
	With    string `json:"with,omitempty"`    // with files, regular-expression
	Without string `json:"without,omitempty"` // without files, regular-expression
	Mode    int    `json:"mode,omitempty"`    // file mode
}

// CompileREs compile regular-expressions found in the output instance.
func (o *Output) CompileREs() (*regexp.Regexp, *regexp.Regexp, error) {
	var withRE *regexp.Regexp
	var withoutRE *regexp.Regexp
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

// Validate inspect instantiated configuration to check if required fields are defined, and regular
// expressions are able to be compiled.
func (c *Config) Validate() error {
	if len(c.Input.Sources) > 0 {
		for i, s := range c.Input.Sources {
			if s.File == "" {
				return fmt.Errorf("%w: hosts.input.sources[%d].file", ErrRequiredAttribute, i)
			}
			if s.URL == "" {
				return fmt.Errorf("%w: hosts.input.sources[%d].url", ErrRequiredAttribute, i)
			}
		}
	}
	if len(c.Input.Transformations) > 0 {
		for i, t := range c.Input.Transformations {
			if t.Search == "" {
				return fmt.Errorf("%w: hosts.input.transformations[%d].search",
					ErrRequiredAttribute, i)
			}
			if _, err := t.CompileRE(); err != nil {
				return fmt.Errorf("%w: hosts.input.transformations[%d].search", err, i)
			}
		}
	}

	if len(c.Output) == 0 {
		log.Warn("No output files will be created!")
	} else {
		for i, o := range c.Output {
			if _, _, err := o.CompileREs(); err != nil {
				return fmt.Errorf("%w: hosts.output[%d].(with|without)", err, i)
			}
			if o.Path == "" {
				return fmt.Errorf("%w: hosts.output[%d].path", ErrRequiredAttribute, i)
			}
		}
	}

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
