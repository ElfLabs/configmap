package file

import (
	"io/ioutil"
	"path/filepath"
)

// File implements a File provider.
type File struct {
	path string
}

// Provider returns a file provider.
func Provider(path string) *File {
	return &File{path: filepath.Clean(path)}
}

// ReadBytes reads the contents of a file on disk and returns the bytes.
func (f *File) ReadBytes() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}
