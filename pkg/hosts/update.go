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

// Update read external blacklist contents, and update configured output file.
type Update struct {
	config  *Config
	dryRun  bool
	content []byte
	skipRes []*regexp.Regexp
}

// readExternalURL executes the GET request to read blacklist URL body.
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

// mappings search and replace contents.
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
		for _, mapping = range mappings {
			lineStr = strings.Replace(string(line[:]), mapping.Search, mapping.Replace, 1)
		}

		mapped = append(mapped, []byte(fmt.Sprintf("%s\n", lineStr))...)
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
	log.Printf("Saving blacklist at: '%s'", filePath)
	return ioutil.WriteFile(filePath, content, 0644)
}

// Execute actions to read blacklist upstream and save/print contents.
func (u *Update) Execute() error {
	var external External
	var err error

	for _, external = range u.config.External {
		var content []byte
		var skipREs []*regexp.Regexp

		if content, err = u.readExternalURL(external.URL); err != nil {
			return err
		}

		if content, err = u.mappingSearchAndReplace(content, external.Mappings); err != nil {
			return err
		}

		for _, re := range external.Skip {
			var compiled *regexp.Regexp

			if compiled, err = regexp.Compile(re); err != nil {
				return err
			}
			skipREs = append(skipREs, compiled)
		}

		if content, err = u.skip(content, skipREs); err != nil {
			return err
		}

		if u.dryRun {
			fmt.Printf("%s", string(content))
		} else {
			if err = u.render(content, external.Output); err != nil {
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
