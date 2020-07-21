package hosts

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// SetLogLevel set the log level based on parameter.
func SetLogLevel(level int) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.Level(level))
}

// DefaultConfigDir returns the default path to application's base directory.
func DefaultConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dirPath := path.Join(home, configDir)
	if info, err := os.Stat(dirPath); err != nil {
		return "", err
	} else if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", dirPath)
	}
	return dirPath, nil
}

// DefaultConfigPath returns the full path to default configuration file location, in users's home.
func DefaultConfigPath() (string, error) {
	baseDir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}

	filePath := path.Join(baseDir, ConfigFile)
	if _, err := os.Stat(filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

// dirGlob search for '.host' files in configured location.
func dirGlob(dir string) ([]string, error) {
	var files []string
	var err error

	pattern := fmt.Sprintf("*.%s", extension)
	if files, err = filepath.Glob(filepath.Join(dir, pattern)); err != nil {
		return nil, err
	}
	return files, nil
}
