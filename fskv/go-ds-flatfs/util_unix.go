//go:build !windows

package flatfs

import (
	"os"
)

func tempFileOnce(dir, pattern string) (*os.File, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.CreateTemp(dir, pattern)
}

func readFileOnce(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
