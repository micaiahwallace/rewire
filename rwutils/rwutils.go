package rwutils

import (
	"bufio"
	"os"
)

// NextLine returns next text line from scanner
func NextLine(scanner *bufio.Scanner) (string, error) {
	scanner.Scan()
	str := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return str, nil
}

// FileExists returns true if filename is a path to a vaild file
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// MakePath creates directory path if any part doesn't exist
func MakePath(path string) error {
	return os.MkdirAll(path, 0755)
}

// MakePaths creates a list of directory paths with MakePath
func MakePaths(paths []string) []error {
	errs := make([]error, 0)
	for _, path := range paths {
		err := MakePath(path)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
