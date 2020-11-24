package hosts

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

// File reprents a dot-host file under application base directory.
type File struct {
	filePath string  // full path to file
	Content  []*Host // host entries
}

// Name file name.
func (f *File) Name() string {
	return path.Base(f.filePath)
}

// Read file data from file-system.
func (f *File) Read() error {
	file, err := os.Open(f.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return f.Load(file)
}

// Load parse informed content and load into file instance.
func (f *File) Load(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		h, err := NewHost(scanner.Text())
		if err != nil {
			log.Tracef("Skipping: '%s'", err)
			continue
		}
		f.Content = append(f.Content, h)
	}
	return nil
}

// Save write file to disk using hosts notation.
func (f *File) Save() error {
	payload := []byte{}
	for _, h := range f.Content {
		line := fmt.Sprintf("%s %s\n", h.Address, h.Hostnames)
		payload = append(payload, []byte(line)...)
	}
	return ioutil.WriteFile(f.filePath, payload, 0o600)
}

// NewFile instantiate a file by path.
func NewFile(filePath string) *File {
	return &File{filePath: filePath}
}
