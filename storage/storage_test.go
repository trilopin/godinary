package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var fileCases = []struct {
	hash   string
	prefix string
	n      int
	dirs   string
	path   string
}{
	{"aabbccdd", "/tmp/", 3, "/tmp/aa/bb/cc/", "/tmp/aa/bb/cc/aabbccdd"},
	{"aabbccdd", "/tmp/", 2, "/tmp/aa/bb/", "/tmp/aa/bb/aabbccdd"},
}

func TestMakeFoldersFromHash(t *testing.T) {
	for _, test := range fileCases {
		dirs, path := makeFoldersFromHash(test.hash, test.prefix, test.n)
		assert.Equal(t, dirs, test.dirs)
		assert.Equal(t, path, test.path)
	}
}
