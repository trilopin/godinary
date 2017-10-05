package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewFileDriver(t *testing.T) {
	viper.Set("fs_base", "base")
	fw := NewFileDriver()
	assert.NotNil(t, fw)
	assert.Equal(t, fw.base, "base")
	viper.Set("fs_base", "")
}

// TODO: replace this fixed jpeso route
func TestWrite(t *testing.T) {
	buf := []byte("CONTENT")
	viper.Set("fs_base", "/tmp/.godtmp/")
	fw := NewFileDriver()
	err := fw.Write(buf, "aabbccddee", "")
	assert.Nil(t, err)

	buf, err = ioutil.ReadFile("/tmp/.godtmp/aa/bb/cc/aabbccddee")
	assert.Nil(t, err)
	assert.Equal(t, "CONTENT", string(buf))
	os.RemoveAll("/tmp/.godtmp/")

}
