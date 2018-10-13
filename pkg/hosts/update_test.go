package hosts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var update *Update

func TestUpdateNewUpdate(t *testing.T) {
	var config *Config
	var err error

	config, _ = NewConfig("../../configs/hosts.yaml")
	update, err = NewUpdate(config, true)

	assert.Nil(t, err)
	assert.NotNil(t, update)
}

func TestUpdateExecute(t *testing.T) {
	var err error

	err = update.Execute()

	assert.Nil(t, err)
}
