package storage

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

// FileDriver struct
type FileDriver struct {
	base string
}

// NewFileDriver constructs new FileDriver with base path
func NewFileDriver() *FileDriver {
	var fs FileDriver
	fs.base = viper.GetString("fs_base")
	if fs.base == "" {
		panic("GODINARY_FS_BASE should be setted")
	}
	return &fs
}

// Write in filesystem a bytearray
func (fs *FileDriver) Write(buf []byte, hash string, prefix string) error {
	dir, newHash := makeFoldersFromHash(hash, fs.base, 3)
	err := os.MkdirAll(prefix+dir, 0744)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(newHash, buf, 0644)
	return err
}

// NewReader produces a file descriptor
func (fs *FileDriver) NewReader(hash string, prefix string) (io.ReadCloser, error) {
	_, newHash := makeFoldersFromHash(hash, fs.base, 3)
	r, err := os.Open(prefix + newHash)
	return r, err
}
