package files

import (
	"errors"
	"fmt"
	"os"
)

// Exists returns true if the given file exists or false if it does not
func Exists(filename string) bool {
	fileStat, err := os.Stat(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}

	if fileStat.IsDir() {
		return false
	}
	return true
}

// ReadAll reads the contents of the given file into a string
func ReadAll(filename string) (string, error) {
	str, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// CreateFileInTmpDir creates a new file in the current tmp directory, if there is a file with the given name
// already it will be replaced
func CreateFileInTmpDir(filename, content string) (string, error) {
	fullFileName := fmt.Sprintf("%s/%s", os.TempDir(), filename)
	if Exists(fullFileName) {
		_ = os.Remove(fullFileName)
	}
	err := os.WriteFile(fullFileName, []byte(content), 0644)
	return fullFileName, err
}
