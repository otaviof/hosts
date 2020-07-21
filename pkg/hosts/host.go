package hosts

import (
	"fmt"
	"regexp"
)

// Host represent a single host file line, with address and host section.
type Host struct {
	Address   string // ip address
	Hostnames string // hostname and aliases
}

// parseRE regular expression to parse host lines.
const parseRE = `^([a-z0-9.:%]+)\s+(\w+.*?)$`

// re compiled regular-expression.
var re = regexp.MustCompile(parseRE)

// NewHost instantiate a host line entry.
func NewHost(entry string) (*Host, error) {
	matches := re.FindStringSubmatch(entry)
	if len(matches) != 3 {
		return nil, fmt.Errorf("informed entry is not a valid host entry")
	}
	return &Host{Address: matches[1], Hostnames: matches[2]}, nil
}
