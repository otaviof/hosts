package hosts

import (
	"bytes"
	"path"

	log "github.com/sirupsen/logrus"
)

// Hosts represents the application itself, implementing the endpoints of sub-commaands.
type Hosts struct {
	cfg     *Config // application configuration
	baseDir string  // base directory path
	files   []*File // instantiated files
}

// Load read dot-host files in base directory.
func (h *Hosts) Load() error {
	log.Info("Loading host files...")
	files, err := dirGlob(h.baseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		log.Infof("Reading file '%s'", file)
		f := NewFile(file)
		if err = f.Read(); err != nil {
			return err
		}
		h.files = append(h.files, f)
	}
	return nil
}

// Update external sources of data, save parsed content in dot-host files.
func (h *Hosts) Update() error {
	log.Infof("Updating external data sources (%d)", len(h.cfg.Input.Sources))
	t, err := NewTransformer(h.cfg.Input.Transformations)
	if err != nil {
		return err
	}

	for _, source := range h.cfg.Input.Sources {
		logger := log.WithFields(log.Fields{
			"name": source.Name,
			"URI":  source.URI,
			"file": source.File,
		})

		u := NewUpdater(t)
		logger.Infof("Updating external source")
		payload, err := u.Get(source.URI)
		if err != nil {
			return err
		}

		filePath := path.Join(h.baseDir, source.File)
		log.Debugf("Saving data at '%s'", filePath)
		f := NewFile(filePath)
		if err = f.Load(bytes.NewReader(payload)); err != nil {
			return err
		}
		if err = f.Save(); err != nil {
			return err
		}
		logger.Info("Done")
	}
	return nil
}

// Apply render output file based on data on base-dir.
func (h *Hosts) Apply() error {
	r := NewRender(h.files)
	for _, output := range h.cfg.Output {
		if err := r.Output(output); err != nil {
			return err
		}
	}
	return nil
}

// NewHosts instante hosts application primary object.
func NewHosts(cfg *Config, baseDir string) *Hosts {
	return &Hosts{
		cfg:     cfg,
		baseDir: baseDir,
		files:   []*File{},
	}
}
