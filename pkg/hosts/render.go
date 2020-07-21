package hosts

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Render represents the render component that format collected data into desired output files.
type Render struct {
	logger *log.Entry // logger
	files  []*File    // host files available
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

		name := f.Name()
		if withRE != nil && !withRE.MatchString(name) {
			r.logger.Debugf("Skipping file '%s' by not matching with clause ('%s')",
				name, withRE.String())
			continue
		} else if withoutRE != nil && withoutRE.MatchString(name) {
			r.logger.Debugf("Skipping file '%s' by matching without clause ('%s')",
				name, withoutRE.String())
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
		r.logger.Debugf("Processing file: '%s'", f.Name())
		payload = append(payload, []byte(fmt.Sprintf("### %s\n", f.Name()))...)
		for _, h := range f.Content {
			hostnames := strings.Split(h.Hostnames, " ")
			data := fn(h.Address, hostnames)
			payload = append(payload, []byte(data)...)
		}
	}
	return payload
}

// Output render payload into desired output file.
func (r *Render) Output(output Output) error {
	logger := r.logger.WithFields(log.Fields{
		"name":    output.Name,
		"path":    output.Path,
		"with":    output.With,
		"without": output.Without,
		"mode":    output.Mode,
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
	logger.Debugf("Amount of files selected: '%d'", len(selectedFiles))
	payload := r.loopFilesContent(selectedFiles, fn)
	logger.Debugf("File size '%d' bytes", len(payload))

	mode := int(0o600)
	if output.Mode > 0 {
		mode = output.Mode
	}
	logger.Infof("Writting file '%s' (%d bytes)", output.Path, len(payload))
	return ioutil.WriteFile(output.Path, payload, os.FileMode(mode))
}

// NewRender instantiate a render by informing all files instances.
func NewRender(files []*File) *Render {
	return &Render{
		logger: log.WithField("component", "render"),
		files:  files,
	}
}
