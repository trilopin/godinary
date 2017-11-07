package storage

import (
	"bytes"
	"io"
)

// Driver is the interface for saving images
type Driver interface {
	Init() error
	Write(buf []byte, hash string, prefix string) error
	NewReader(hash string, prefix string) (io.ReadCloser, error)
}

// makeFoldersFromHash compute new path in n folders and prefix based on current path
func makeFoldersFromHash(path string, prefix string, n int) (string, string) {
	var newPath bytes.Buffer
	newPath.WriteString(prefix)
	for i := 0; i < n; i++ {
		newPath.WriteString(path[i*2 : i*2+2])
		newPath.WriteString("/")
	}
	dir := newPath.String()
	newPath.WriteString(path)
	return dir, newPath.String()
}
