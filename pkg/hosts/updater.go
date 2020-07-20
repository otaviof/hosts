package hosts

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Updater executes the update process of external data sources. Also handles the expected
// transformation of data, informed in configuration file.
type Updater struct {
	client      http.Client
	transformer *Transformer
}

// Get executes a http.Get against informed URI. The data obtained goes through the transformer rules
// before being returned by this method.
func (u *Updater) Get(uri string) ([]byte, error) {
	response, err := u.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	log.Debugf("Returned status code '%d'", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-expected status '%d'", response.StatusCode)
	}

	payload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return u.transformer.Transform(payload)
}

// NewUpdater instantiate a updater with a transformer instance.
func NewUpdater(t *Transformer) *Updater {
	return &Updater{transformer: t}
}
