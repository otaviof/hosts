package hosts

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigNewConfig(t *testing.T) {
	var config *Config
	var err error

	config, err = NewConfig("../../configs/hosts.yaml")

	log.Printf("Config: '%#v'", config)

	assert.Nil(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, config.Hosts.BaseDirectory, fmt.Sprintf("%s/.hosts", os.Getenv("HOME")))
	assert.True(t, config.Hosts.Output != "")
}

func TestConfigValidate(t *testing.T) {
	var config = &Config{}
	var err error

	err = config.Validate()

	assert.NotNil(t, err)
}
