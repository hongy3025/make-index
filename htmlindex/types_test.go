package htmlindex

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTargetGenerate(t *testing.T) {
	root, err := ioutil.TempDir("", "test-target-generate")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	defer os.RemoveAll(root)

	if err = os.Mkdir(filepath.Join(root, "sub1"), os.ModePerm); err != nil {
		t.Fatalf("error creating test dir: %v", err)
	}
	if f, err := os.Create(filepath.Join(root, "file.txt")); err != nil {
		t.Fatalf("error creating test file: %v", err)
	} else {
		f.Close()
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		t.Fatalf("error reading dir: %v", err)
	}
	target := target{Path: root, Files: files}

	templ := template.Must(template.New("test").Parse("{{ .Name }}{{range .Files}}<li><a href='{{.Name}}'>{{.Name}}</a>{{end}}"))
	err = target.Generate(templ)
	if err != nil {
		t.Fatalf("error generating: %v", err)
	}

	genfile := filepath.Join(root, "index.html")
	generated, err := ioutil.ReadFile(genfile)
	if err != nil {
		t.Fatalf("error reading in generated file: %v", err)
	}

	expected := filepath.Base(root) + "<li><a href='file.txt'>file.txt</a><li><a href='sub1'>sub1</a>"
	if string(generated) != expected {
		t.Errorf("generate produced:\n%q\n\nexpected:\n%q", string(generated), expected)
	}
}
