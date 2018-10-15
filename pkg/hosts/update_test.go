package hosts

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config *Config
var update *Update
var reSearchReplace map[*regexp.Regexp]string

func TestUpdateNewUpdate(t *testing.T) {
	config, _ = NewConfig("../../configs/hosts.yaml")
	update = NewUpdate(config, true)

	assert.NotNil(t, update)
}

func TestUpdateExecute(t *testing.T) {
	var err error

	err = update.Execute()

	assert.Nil(t, err)
	assert.True(t, len(update.content) > 0)
}

func TestUpdateReCompile(t *testing.T) {
	var err error

	reSearchReplace, err = update.reCompile(config.External[0].Transform)

	assert.Nil(t, err)
}

func TestUpdateTransform(t *testing.T) {
	var reader = bufio.NewReader(bytes.NewReader(update.content))
	var line []byte
	var err error

	for {
		var search *regexp.Regexp
		var replace string

		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		for search, replace = range reSearchReplace {
			// means the matching line should have been skipped
			if replace == "" {
				assert.False(t, search.Match(line))
			}
		}
	}
}
