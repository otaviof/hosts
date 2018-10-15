package hosts

import (
	"io/ioutil"
	"log"
	"os"
)

// readFile Wrap up a ioutil call, using fatal log in case of error.
func readFile(path string) []byte {
	log.Printf("Reading file: '%s'", path)

	if !exists(path) {
		log.Fatalf("Can't find file: '%s'", path)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return fileBytes
}

// isDir Check if informed path is a directory, boolean return.
func isDir(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// exists Check if path exists, boolean return.
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// stringSliceContains checks if a slice contiains a string.
func stringSliceContains(slice []string, str string) bool {
	var sliceStr string

	for _, sliceStr = range slice {
		if str == sliceStr {
			return true
		}
	}

	return false
}
