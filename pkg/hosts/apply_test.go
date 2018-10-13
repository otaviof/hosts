package hosts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var apply *Apply

func TestApplyNewApply(t *testing.T) {
	var config = &Config{Hosts: Hosts{BaseDirectory: "../../test/hosts-dir"}}
	var err error

	apply, err = NewApply(config, true)

	assert.Nil(t, err)
	assert.NotNil(t, apply)

	assert.Equal(t, 2, len(apply.files))
}

func TestApplyExecute(t *testing.T) {
	var err error

	err = apply.Execute()

	assert.Nil(t, err)
}
