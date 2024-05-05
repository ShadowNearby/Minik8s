package utils

import (
	"io"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	content, err := io.ReadAll(fd)

	if err != nil {
		return nil, err
	}
	return content, nil
}
