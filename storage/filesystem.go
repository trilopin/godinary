package storage

import (
	"io"
	"io/ioutil"
	"os"
)

// FileDriver struct
type FileDriver struct {
	base string
}

// NewFileDriver constructs new FileDriver with base path
func NewFileDriver(base string) *FileDriver {
	var fs FileDriver
	fs.base = base
	return &fs
}

// Init does nothing in this implementation
func (fs *FileDriver) Init() error {
	return nil
}

// Write in filesystem a bytearray
func (fs *FileDriver) Write(buf []byte, hash string, prefix string) error {
	dir, newHash := makeFoldersFromHash(hash, fs.base+prefix, 3)
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(newHash, buf, 0644)
	return err
}

// NewReader produces a file descriptor
func (fs *FileDriver) NewReader(hash string, prefix string) (io.ReadCloser, error) {
	_, newHash := makeFoldersFromHash(hash, fs.base+prefix, 3)
	r, err := os.Open(newHash)
	return r, err
}
