package hosts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilsStringSliceContains(t *testing.T) {
	var sliceA = []string{"a", "b", "c"}
	var contains bool

	contains = stringSliceContains(sliceA, "b")
	assert.True(t, contains)

	contains = stringSliceContains(sliceA, "d")
	assert.False(t, contains)
}
