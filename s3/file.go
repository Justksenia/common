package s3

import (
	"bytes"
	"io"
)

type File struct {
	name string
	path string
	data io.Reader
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Data() ([]byte, error) {
	return io.ReadAll(f.data)
}

func NewFile(name, path string, data []byte) *File {
	return &File{
		name: name,
		path: path,
		data: bytes.NewReader(data),
	}
}
