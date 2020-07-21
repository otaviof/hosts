package hosts

import (
	"errors"
	"fmt"
	"regexp"
)

// Host represent a single host file line, with address and host section.
type Host struct {
	Address   string // ip address
	Hostnames string // hostname and aliases
}

const parseRE = `^([a-z0-9.:%]+)\s+(\w+.*?)$`

var re = regexp.MustCompile(parseRE)

var invalidHostEntryErr = errors.New("hosts line does not match expected format")

// NewHost instantiate a host line entry.
func NewHost(entry string) (*Host, error) {
	matches := re.FindStringSubmatch(entry)
	if len(matches) != 3 {
		return nil, fmt.Errorf("%w: '%s'", invalidHostEntryErr, entry)
	}
	return &Host{Address: matches[1], Hostnames: matches[2]}, nil
}
