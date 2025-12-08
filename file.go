package mykonf

import (
	"os"

	"github.com/knadh/koanf/providers/file"
)

type File struct {
	*file.File
}

func Provider(path string) File {
	return File{File: file.Provider(path)}
}

func (f File) ReadBytes() (b []byte, err error) {
	b, err = f.File.ReadBytes()
	if err != nil {
		return nil, err
	}
	return []byte(os.ExpandEnv(string(b))), nil
}
