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

// fetchBlacklist executes the GET request to read blacklist URL body.
func (u *Update) fetchBlacklist() error {
	var client http.Client
	var resp *http.Response
	var err error

	log.Printf("Reading URL: '%s'", u.config.External.URL)
	if resp, err = client.Get(u.config.External.URL); err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status-code is not okay: '%d'", resp.StatusCode)
	}

	if u.content, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	return nil
}

// mappings search and replace contents.
func (u *Update) mappings() error {
	var mapped []byte
	var reader = bufio.NewReader(bytes.NewReader(u.content))
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
			return err
		}
		for _, mapping = range u.config.External.Mappings {
			lineStr = strings.Replace(string(line[:]), mapping.Search, mapping.Replace, 1)
		}

		mapped = append(mapped, []byte(fmt.Sprintf("%s\n", lineStr))...)
	}

	u.content = mapped
	return nil
}

// skip lines that are matching one of the configured regular expressions.
func (u *Update) skip() error {
	var reader = bufio.NewReader(bytes.NewReader(u.content))
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
			return err
		}

		for _, re := range u.skipRes {
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

	u.content = kept
	return nil
}

// print out body contents, after mapping and skipping.
func (u *Update) print() {
	fmt.Printf("%s", string(u.content))
}

// render write contents to output file.
func (u *Update) render() error {
	var filePath = path.Join(u.config.Hosts.BaseDirectory, u.config.External.Output)
	log.Printf("Saving blacklist at: '%s'", filePath)
	return ioutil.WriteFile(filePath, u.content, 0600)
}

// Execute actions to read blacklist upstream and save/print contents.
func (u *Update) Execute() error {
	var err error

	if err = u.fetchBlacklist(); err != nil {
		return err
	}
	if err = u.mappings(); err != nil {
		return err
	}
	if err = u.skip(); err != nil {
		return err
	}

	if u.dryRun {
		u.print()
	} else {
		if err = u.render(); err != nil {
			return err
		}
	}

	return nil
}

// NewUpdate creates a new update instance, by initializing .
func NewUpdate(config *Config, dryRun bool) (*Update, error) {
	var update = &Update{config: config, dryRun: dryRun}
	var err error

	for _, skipRe := range config.External.Skip {
		var compiled *regexp.Regexp

		if compiled, err = regexp.Compile(skipRe); err != nil {
			return nil, err
		}
		update.skipRes = append(update.skipRes, compiled)
	}
	log.Printf("%s", update.content)

	return update, nil
}
