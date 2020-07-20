package hosts

import (
	"fmt"
	"regexp"
)

type Host struct {
	Address string
	Host    string
}

const parseRE = `^([a-z0-9.:%]+)\s+(\w+.*?)$`

var re = regexp.MustCompile(parseRE)

func NewHost(entry string) (*Host, error) {
	matches := re.FindStringSubmatch(entry)
	if len(matches) != 3 {
		return nil, fmt.Errorf("informed entry is not a valid host entry")
	}
	return &Host{Address: matches[1], Host: matches[2]}, nil
}
