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

// render write contents to output file.
func (u *Update) render(content []byte, name string) error {
	var filePath = path.Join(u.config.Hosts.BaseDirectory, name)
	log.Printf("Saving external resource data at: '%s'", filePath)
	return ioutil.WriteFile(filePath, content, 0644)
}

// reCompile compile the regular expressions in `transform` block.
func (u *Update) reCompile(t []Transform) (map[*regexp.Regexp]string, error) {
	var reSearchReplace = make(map[*regexp.Regexp]string)
	var searchReplace Transform
	var compiled *regexp.Regexp
	var err error

	for _, searchReplace = range t {
		if compiled, err = regexp.Compile(searchReplace.Search); err != nil {
			return nil, err
		}
		reSearchReplace[compiled] = searchReplace.Replace
	}

	return reSearchReplace, nil
}

// transform content to either skip a line, or repliace contents based in regular expressions.
func (u *Update) transform(content []byte, t []Transform) ([]byte, error) {
	var reSearchReplace map[*regexp.Regexp]string
	var reader = bufio.NewReader(bytes.NewReader(content))
	var transformed []byte
	var err error

	if reSearchReplace, err = u.reCompile(t); err != nil {
		return nil, err
	}

	for {
		var line []byte
		var search *regexp.Regexp
		var replace string
		var skip = false

		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for search, replace = range reSearchReplace {
			if search.Match(line) {
				// replace is empty, therefore we are skipping this line
				if replace == "" {
					skip = true
					break
				}
				// replacing content entries
				line = search.ReplaceAll(line, []byte(replace))
			}
		}

		if !skip {
			transformed = append(transformed, line...)
			transformed = append(transformed, []byte("\n")...)
		}
	}

	return transformed, nil
}

// Execute actions to read upstream resources and save/print contents.
func (u *Update) Execute() error {
	var external External
	var err error

	for _, external = range u.config.External {
		if u.content, err = u.readExternalURL(external.URL); err != nil {
			return err
		}

		if u.content, err = u.transform(u.content, external.Transform); err != nil {
			return err
		}

		// saving or printing content
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

// NewUpdate creates a new update instance, by initializing.
func NewUpdate(config *Config, dryRun bool) *Update {
	return &Update{config: config, dryRun: dryRun}
}
