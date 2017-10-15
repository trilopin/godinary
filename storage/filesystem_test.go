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

// func TestWriteFail(t *testing.T) {
// 	buf := []byte("CONTENT")
// 	fw := NewFileDriver("/root/")
// 	err := fw.Write(buf, "aabbccddee", "")
// 	assert.NotNil(t, err)
// }

func TestNewReader(t *testing.T) {
	fw := NewFileDriver("/fakedir/")
	r, err := fw.NewReader("aabbccddee", "/tmp/")
	assert.Equal(t, err.Error(), "open /fakedir//tmp/aa/bb/cc/aabbccddee: no such file or directory")
	assert.Nil(t, r)
}
