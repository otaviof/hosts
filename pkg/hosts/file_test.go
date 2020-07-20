package hosts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFile(t *testing.T) {
	f := NewFile("../../test/hosts-dir/00-localhost.host")

	t.Run("read", func(t *testing.T) {
		err := f.Read()
		require.NoError(t, err, "should not error when reading a existing file")
		require.Len(t, f.Content, 4, "expect certain amount of lines in file")
	})

	t.Run("load", func(t *testing.T) {
		file, err := os.Open("../../test/hosts-dir/01-example.host")
		require.NoError(t, err, "should not error reading a test host file")
		defer file.Close()

		err = f.Load(file)
		require.NoError(t, err, "should not error on parsing host file content")
		require.Len(t, f.Content, 5, "expect certain amount of lines in file")
	})
}
