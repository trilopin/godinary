package storage

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// FileDriver struct
type FileDriver struct {
	base string
}

// NewFileDriver constructs new FileDriver with base path
func NewFileDriver() *FileDriver {
	var fs FileDriver
	fs.base = os.Getenv("GODINARY_FS_BASE")
	if fs.base == "" {
		panic("GODINARY_FS_BASE should be setted")
	}
	return &fs
}

// Write in filesystem a bytearray
func (fs *FileDriver) Write(buf []byte, hash string) error {
	dir, newHash := makeFoldersFromHash(hash, fs.base, 3)
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(newHash, buf, 0644)
	return err
}

func (fs *FileDriver) Read(hash string) (io.Reader, error) {
	_, newHash := makeFoldersFromHash(hash, fs.base, 3)
	log.Println(newHash)
	r, err := os.Open(newHash)
	return r, err
}
