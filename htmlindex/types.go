package htmlindex

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Index interface {
	Generate(*template.Template) error
}

func NewTarget(path string) (Index, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(abspath)
	if err != nil {
		return nil, err
	}
	return target{Path: abspath, Files: files}, nil
}

type target struct {
	Path  string
	Files []os.FileInfo
}

func (t target) Name() string {
	return filepath.Base(t.Path)
}

func (t target) Generate(templ *template.Template) error {
	index := filepath.Join(t.Path, "index.html")
	f, err := os.Create(index)
	if err != nil {
		return err
	}
	defer f.Close()

	return templ.Execute(f, t)
}
