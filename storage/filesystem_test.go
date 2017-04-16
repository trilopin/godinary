package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileDriver(t *testing.T) {
	os.Setenv("GODINARY_FS_BASE", "BASE")
	fw := NewFileDriver()
	assert.NotNil(t, fw)
	assert.Equal(t, fw.base, "BASE")
	os.Setenv("GODINARY_FS_BASE", "")
}

// TODO: replace this fixed jpeso route
func TestWrite(t *testing.T) {
	buf := []byte("CONTENT")
	os.Setenv("GODINARY_FS_BASE", "/tmp/.godtmp/")
	fw := NewFileDriver()
	err := fw.Write(buf, "aabbccddee")
	assert.Nil(t, err)

	buf, err = ioutil.ReadFile("/tmp/.godtmp/aa/bb/cc/aabbccddee")
	assert.Nil(t, err)
	assert.Equal(t, "CONTENT", string(buf))
	os.RemoveAll("/tmp/.godtmp/")

}
