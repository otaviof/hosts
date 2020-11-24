package hosts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var ErrNonExpectedStatus = errors.New("non-expected status returned")

// Updater executes the update process of external data sources. Also handles the expected
// transformation of data, informed in configuration file.
type Updater struct {
	logger      *log.Entry   // logger
	client      http.Client  // http client
	transformer *Transformer // transformer instance
}

// Get executes a http.Get against informed URI. The data obtained goes through the transformer rules
// before being returned by this method.
func (u *Updater) Get(uri string) ([]byte, error) {
	u.logger.Debugf("Making a request on '%s'", uri)
	response, err := u.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	u.logger.Debugf("Returned status code '%d'", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrNonExpectedStatus, response.StatusCode)
	}

	payload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return u.transformer.Transform(payload)
}

// NewUpdater instantiate a updater with a transformer instance.
func NewUpdater(t *Transformer) *Updater {
	return &Updater{
		logger:      log.WithField("component", "updater"),
		transformer: t,
	}
}
