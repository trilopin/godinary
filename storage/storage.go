package storage

import (
	"bytes"
	"io"
)

var (
	// StorageDriver is global struct for persistence driver
	StorageDriver Driver
)

// Driver is the interface for saving images
type Driver interface {
	Write(buf []byte, hash string) error
	Read(hash string) (io.Reader, error)
}

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
