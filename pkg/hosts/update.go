package hosts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

// Update read external resource contents, and update configured output files.
type Update struct {
	config  *Config
	content []byte
	dryRun  bool
}

// readExternalURL execute a GET request to read external resource URL body.
func (u *Update) readExternalURL(URL string) ([]byte, error) {
	var client http.Client
	var resp *http.Response
	var err error

	log.Printf("External resoure URL: '%s'", URL)
	if resp, err = client.Get(URL); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status-code is not okay: '%d'", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

// mappingSearchAndReplace search and replace contents.
func (u *Update) mappingSearchAndReplace(content []byte, mappings []Mapping) ([]byte, error) {
	var mapped []byte
	var reader = bufio.NewReader(bytes.NewReader(content))
	var err error

	log.Printf("Replacing output accordingly with 'mapping'...")
	for {
		var line []byte
		var lineStr string
		var mapping Mapping

		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// applying the search and replace per line
		for _, mapping = range mappings {
			lineStr = strings.Replace(string(line[:]), mapping.Search, mapping.Replace, 1)
			line = []byte(lineStr)
		}

		mapped = append(mapped, line...)
		mapped = append(mapped, []byte("\n")...)
	}

	return mapped, nil
}

// skip lines that are matching one of the configured regular expressions.
func (u *Update) skip(content []byte, res []*regexp.Regexp) ([]byte, error) {
	var reader = bufio.NewReader(bytes.NewReader(content))
	var kept []byte
	var err error

	log.Printf("Skipping output matching 'skip'...")
	for {
		var line []byte
		var skipLine = false

		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for _, re := range res {
			if re.Match(line) {
				skipLine = true
				log.Printf("[SKIP] regexp='%s', line='%s'", re, line)
				break
			}
		}

		if !skipLine {
			kept = append(kept, line...)
			kept = append(kept, []byte("\n")...)
		}
	}

	return kept, nil
}

// render write contents to output file.
func (u *Update) render(content []byte, name string) error {
	var filePath = path.Join(u.config.Hosts.BaseDirectory, name)
	log.Printf("Saving external resource data at: '%s'", filePath)
	return ioutil.WriteFile(filePath, content, 0644)
}

// Execute actions to read upstream resources and save/print contents.
func (u *Update) Execute() error {
	var external External
	var err error

	for _, external = range u.config.External {
		var skipREs []*regexp.Regexp
		var compiled *regexp.Regexp

		// reading external data
		if u.content, err = u.readExternalURL(external.URL); err != nil {
			return err
		}

		// applying the initial search and replace, using strings
		if u.content, err = u.mappingSearchAndReplace(u.content, external.Mappings); err != nil {
			return err
		}

		// compiling regular expressions to be used in `skip` method
		for _, re := range external.Skip {
			if compiled, err = regexp.Compile(re); err != nil {
				return err
			}
			skipREs = append(skipREs, compiled)
		}

		// skipping lines using regular expresions
		if u.content, err = u.skip(u.content, skipREs); err != nil {
			return err
		}

		// saving or printing contens
		if u.dryRun {
			fmt.Printf("%s", string(u.content))
		} else {
			if err = u.render(u.content, external.Output); err != nil {
				return err
			}
		}
	}

	return nil
}

// NewUpdate creates a new update instance, by initializing .
func NewUpdate(config *Config, dryRun bool) *Update {
	return &Update{config: config, dryRun: dryRun}
}
