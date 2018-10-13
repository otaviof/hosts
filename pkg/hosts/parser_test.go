package hosts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var parser *Parser

func TestParserNewParser(t *testing.T) {
	var err error

	parser, err = NewParser("../../test/hosts-dir/00-localhost.host")

	assert.Nil(t, err)
	assert.NotNil(t, parser)
}

func TestParserIngest(t *testing.T) {
	var err error

	err = parser.Ingest()

	assert.Nil(t, err)
	assert.True(t, len(parser.Contents) > 1)
}
