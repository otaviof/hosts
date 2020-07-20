package hosts

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Render represents the render component that format collected data into desired output files.
type Render struct {
	files []*File // host files available
}

// formatterFn signature for functions used on formatting.
type formatterFn func(string, []string) string

// selectFiles based on with and without regular expressions.
func (r *Render) selectFiles(withRE, withoutRE *regexp.Regexp) []*File {
	selected := []*File{}
	for _, f := range r.files {
		if withRE == nil && withoutRE == nil {
			selected = append(selected, f)
			continue
		}

		if withRE != nil && !withRE.MatchString(f.Name()) {
			continue
		} else if withoutRE != nil && withoutRE.MatchString(f.Name()) {
			continue
		}
		selected = append(selected, f)
	}
	return selected
}

// etcHostsFormatter format input for /etc/hosts.
func (r *Render) etcHostsFormatter(address string, hostnames []string) string {
	return fmt.Sprintf("%s %s\n", address, strings.Join(hostnames, " "))
}

// dnsmasqFormatter format input for dnsmasq.
func (r *Render) dnsmasqFormatter(address string, hostnames []string) string {
	lines := []string{}
	for _, hostname := range hostnames {
		lines = append(lines, fmt.Sprintf("address=/%s/%s\n", hostname, address))
	}
	return strings.Join(lines, "\n")
}

// loopFilesContent loop through files and their content, marking each file with a header, and
// contents informed to formmater function. It returns a complete payload with headers and formatted
// output file content.
func (r *Render) loopFilesContent(files []*File, fn formatterFn) []byte {
	payload := []byte{}
	for _, f := range files {
		payload = append(payload, []byte(fmt.Sprintf("### %s\n", f.Name()))...)
		for _, h := range f.Content {
			hostnames := strings.Split(h.Host, " ")
			data := fn(h.Address, hostnames)
			payload = append(payload, []byte(data)...)
		}
	}
	return payload
}

// Output render payload into desired output file.
func (r *Render) Output(output Output) error {
	logger := log.WithFields(log.Fields{
		"name":    output.Name,
		"path":    output.Path,
		"with":    output.With,
		"without": output.Without,
	})

	logger.Debugf("Compiling regular expressions")
	withRE, withoutRE, err := output.CompileREs()
	if err != nil {
		return err
	}

	fn := r.etcHostsFormatter
	if output.Dnsmasq {
		logger.Debugf("Using DNSMasq formatter")
		fn = r.dnsmasqFormatter
	}

	selectedFiles := r.selectFiles(withRE, withoutRE)
	logger.Debugf("Files selected: '%#v'", selectedFiles)
	payload := r.loopFilesContent(selectedFiles, fn)
	logger.Debugf("File size '%d' bytes", len(payload))
	return ioutil.WriteFile(output.Path, payload, 0644)
}

// NNewRender instantiate a render by informing all files instances.
func NewRender(files []*File) *Render {
	return &Render{files: files}
}
