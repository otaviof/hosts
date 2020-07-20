package hosts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTransformer(t *testing.T) {
	payload := `
data
to-replace
data
to-skip
data
`
	expected := `
data
replaced
data
data
`
	transformations := []Transformation{
		{
			Search:  "^to-replace$",
			Replace: "replaced",
		},
		{
			Search: "^to-skip$",
		},
	}
	transformer, err := NewTransformer(transformations)
	require.NoError(t, err, "should not error on compiling test regexps")

	t.Run("transform", func(t *testing.T) {
		transformed, err := transformer.Transform([]byte(payload))
		require.NoError(t, err, "should not error on applying transformations")
		require.Equal(t, expected, string(transformed), "transformations should match expected")
	})
}
