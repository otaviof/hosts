package hosts

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// Transformer represents the transformations applied to data. It holds the regular-expressions
// compiled, and apply them against each line replacing data or skipping lines.
type Transformer struct {
	logger             *log.Entry                // logger
	transformations    []Transformation          // list of transformations
	regularExpressions map[*regexp.Regexp]string // map of compiled regexp and replace
}

// applyREs apply pre-compiled regular-expressions against informed payload. In case of payload meant
// to be skipped, it returns nil.
func (t *Transformer) applyREs(payload []byte) []byte {
	for search, replace := range t.regularExpressions {
		if !search.Match(payload) {
			continue
		}
		t.logger.Tracef("Line matches regular-expression '%s' (replace '%s'): '%s'",
			search.String(), replace, payload)
		if replace == "" {
			return nil
		}
		payload = search.ReplaceAll(payload, []byte(replace))
	}
	return payload
}

// Transform apply transformations to each line of the payload, returning the result.
func (t *Transformer) Transform(payload []byte) ([]byte, error) {
	r := bufio.NewReader(bytes.NewReader(payload))
	transformed := []byte{}
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if line = t.applyREs(line); line != nil {
			line = append(line, []byte("\n")...)
			transformed = append(transformed, line...)
		}
	}
	return transformed, nil
}

// compileREs compile all regular-expressions found on informed transformations, saving them in the
// internal representation. A shared map holds compiled regular expression and the replace string.
func (t *Transformer) compileREs() error {
	for _, transformation := range t.transformations {
		re, err := transformation.CompileRE()
		t.logger.Debugf("Compiling regular-expression for search '%s', replace '%s'",
			transformation.Search, transformation.Replace)
		if err != nil {
			return err
		}
		t.regularExpressions[re] = transformation.Replace
	}
	return nil
}

// NewTransformer instantiate a new transformer by compiling and preparing regular-expressions.
func NewTransformer(transformations []Transformation) (*Transformer, error) {
	t := &Transformer{
		logger:             log.WithField("component", "transformer"),
		transformations:    transformations,
		regularExpressions: map[*regexp.Regexp]string{},
	}
	return t, t.compileREs()
}
