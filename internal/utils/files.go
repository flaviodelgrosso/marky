package utils

import (
	"fmt"
	"os"
)

// ReadFile reads a file and returns its content as bytes.
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", path, err)
	}
	return data, nil
}
