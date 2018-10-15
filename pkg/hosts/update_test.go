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

func TestUpdateNewUpdate(t *testing.T) {
	config, _ = NewConfig("../../configs/hosts.yaml")
	update = NewUpdate(config, true)

	assert.NotNil(t, update)
}

func TestUpdateExecute(t *testing.T) {
	var err error

	err = update.Execute()

	assert.Nil(t, err)
}

func TestUpdateMapping(t *testing.T) {
	var reader = bufio.NewReader(bytes.NewReader(update.content))
	var line []byte
	var lineStr string
	var mapping Mapping
	var err error

	for {
		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		lineStr = string(line[:])
		for _, mapping = range config.External[0].Mappings {
			assert.NotSubset(t, mapping.Search, lineStr)
			assert.Subset(t, mapping.Replace, lineStr)
		}
	}
}

func TestUpdateSkip(t *testing.T) {
	var reader = bufio.NewReader(bytes.NewReader(update.content))
	var skipREs []*regexp.Regexp
	var line []byte
	var err error

	for _, re := range config.External[0].Skip {
		compiled, _ := regexp.Compile(re)
		skipREs = append(skipREs, compiled)
	}

	for {
		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		for _, re := range skipREs {
			assert.False(t, re.Match(line))
		}
	}
}
