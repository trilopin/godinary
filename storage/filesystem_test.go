package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileDriver(t *testing.T) {
	fw := NewFileDriver("base")
	assert.NotNil(t, fw)
	assert.Equal(t, fw.base, "base")
}

func TestWrite(t *testing.T) {
	buf := []byte("CONTENT")
	fw := NewFileDriver("/tmp/.godtmp/")
	err := fw.Write(buf, "aabbccddee", "")
	assert.Nil(t, err)

	buf, err = ioutil.ReadFile("/tmp/.godtmp/aa/bb/cc/aabbccddee")
	assert.Nil(t, err)
	assert.Equal(t, "CONTENT", string(buf))
	os.RemoveAll("/tmp/.godtmp/")

}
