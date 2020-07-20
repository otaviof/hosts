package hosts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUpdater(t *testing.T) {
	transformer, _ := NewTransformer([]Transformation{})
	u := NewUpdater(transformer)

	t.Run("get", func(t *testing.T) {
		payload, err := u.Get("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
		require.NoError(t, err, "should not error on reaching out to external server")
		require.True(t, len(payload) > 10, "expect payload to not be empty")
	})
}
