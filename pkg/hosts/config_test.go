package hosts

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const testBaseDir = "../../test/hosts-dir"

func newConfig(t *testing.T) *Config {
	cfg, err := NewConfig(path.Join(testBaseDir, "hosts.yaml"))
	require.NoError(t, err, "should not error on parsing test configuration")
	return cfg
}

func TestNewConfig(t *testing.T) {
	cfg := newConfig(t)

	t.Run("validate", func(t *testing.T) {
		err := cfg.Validate()
		require.NoError(t, err, "should not error on validation")
	})
}
