package hosts

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// Parser for a '.host' files
type Parser struct {
	filePath    string
	ipv4Re      *regexp.Regexp
	ipv6Re      *regexp.Regexp
	ipv6AliasRe *regexp.Regexp
	Contents    []string
}

// loadRegexp load the regular expressions that are used to identify addresses.
func (p *Parser) loadRegexp() error {
	var ipv6ReStr = `^([0-9a-f]|:){1,4}(:([0-9a-f]{0,4})*){1,7}`
	var err error

	if p.ipv4Re, err = regexp.Compile(`^(\d{1,3}\.){3}\d{1,3}`); err != nil {
		return err
	}

	if p.ipv6Re, err = regexp.Compile(ipv6ReStr); err != nil {
		return err
	}

	if p.ipv6AliasRe, err = regexp.Compile(ipv6ReStr + `(%\w+)`); err != nil {
		return err
	}

	return nil
}

// extract investigate line to extract ipv4 and ipv6 addresses.
func (p *Parser) extract(line []byte) error {
	var lineStr = string(line)

	if !p.ipv4Re.Match(line) && !p.ipv6Re.Match(line) && !p.ipv6AliasRe.Match(line) {
		return fmt.Errorf("cannot find ipv4/ipv6 address in line '%s'", lineStr)
	}

	p.Contents = append(p.Contents, lineStr)
	return nil
}

// Ingest read host file and parse it's contents.
func (p *Parser) Ingest() error {
	var file *os.File
	var reader *bufio.Reader
	var err error

	log.Printf("Inpecting host file: '%s'", p.filePath)

	if file, err = os.Open(p.filePath); err != nil {
		return err
	}
	defer file.Close()

	reader = bufio.NewReader(file)
	for {
		var line []byte

		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if err = p.extract(line); err != nil {
			log.Printf("[WARN] %s", err)
		}
	}

	return nil
}

// NewParser creates a new instance of Parser, given a dot host file input.
func NewParser(filePath string) (*Parser, error) {
	var parser = &Parser{filePath: filePath}
	var err error

	if !exists(filePath) {
		return nil, fmt.Errorf("can not find '%s' file", filePath)
	}

	if err = parser.loadRegexp(); err != nil {
		return nil, err
	}

	return parser, nil
}
