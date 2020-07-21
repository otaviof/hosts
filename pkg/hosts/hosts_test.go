package hosts

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHosts(t *testing.T) {
	cfg := newConfig(t)
	hosts := NewHosts(cfg, testBaseDir, false)

	// cleaning up blocks file before testing
	_ = os.RemoveAll(path.Join(testBaseDir, cfg.Input.Sources[0].File))

	t.Run("load", func(t *testing.T) {
		err := hosts.Load()
		require.NoError(t, err, "should not error loading test files")
		require.Len(t, hosts.files, 2, "should have two files instantiated")
	})

	t.Run("update", func(t *testing.T) {
		err := hosts.Update()
		require.NoError(t, err, "should not error updating external sources")
	})

	t.Run("apply", func(t *testing.T) {
		err := hosts.Apply()
		require.NoError(t, err, "should not error updating applying output")
	})
}
