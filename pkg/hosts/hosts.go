package hosts

import (
	"bytes"
	"path"

	log "github.com/sirupsen/logrus"
)

// Hosts represents the application itself, implementing the endpoints of sub-commaands.
type Hosts struct {
	logger  *log.Entry // logger
	cfg     *Config    // application configuration
	baseDir string     // base directory path
	files   []*File    // instantiated files
	dryRun  bool       // dry-run flag
}

// Load read dot-host files in base directory.
func (h *Hosts) Load() error {
	h.logger.Info("Loading host files...")
	files, err := dirGlob(h.baseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		h.logger.Infof("Reading file '%s'", file)
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
	t, err := NewTransformer(h.cfg.Input.Transformations)
	if err != nil {
		return err
	}

	h.logger.Infof("Updating external data sources (%d)", len(h.cfg.Input.Sources))
	for _, source := range h.cfg.Input.Sources {
		logger := h.logger.WithFields(log.Fields{
			"name":    source.Name,
			"URI":     source.URL,
			"file":    source.File,
			"dry-run": h.dryRun,
		})

		u := NewUpdater(t)
		logger.Infof("Updating external source")
		payload, err := u.Get(source.URL)
		if err != nil {
			return err
		}

		filePath := path.Join(h.baseDir, source.File)
		logger.Infof("Saving data at '%s'", filePath)
		f := NewFile(filePath)
		if err = f.Load(bytes.NewReader(payload)); err != nil {
			return err
		}
		if h.dryRun {
			logger.Info("Dry-run mode, file not saved.")
			logger.Tracef("%s", payload)
		} else {
			if err = f.Save(); err != nil {
				return err
			}
		}
		logger.Info("Done")
	}
	return nil
}

// Apply render output file based on data on base-dir.
func (h *Hosts) Apply() error {
	r := NewRender(h.files, false)
	for _, output := range h.cfg.Output {
		if err := r.Output(output); err != nil {
			return err
		}
	}
	return nil
}

// NewHosts instante hosts application primary object.
func NewHosts(cfg *Config, baseDir string, dryRun bool) *Hosts {
	return &Hosts{
		logger:  log.WithField("component", "hosts"),
		cfg:     cfg,
		baseDir: baseDir,
		files:   []*File{},
		dryRun:  dryRun,
	}
}
