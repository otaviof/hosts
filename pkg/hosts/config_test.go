package hosts

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigNewConfig(t *testing.T) {
	var config *Config
	var err error

	config, err = NewConfig("../../configs/hosts.yaml")

	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.Hosts.BaseDirectory, fmt.Sprintf("%s/.hosts", os.Getenv("HOME")))
}

func TestConfigValidate(t *testing.T) {
	var config = &Config{}
	var err error

	err = config.Validate()

	assert.NotNil(t, err)
}
