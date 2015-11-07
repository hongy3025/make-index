package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/johnweldon/make-index/htmlindex"
)

func main() {
	templ := template.Must(template.ParseFiles("template.html"))

	for d := range noIndex(hasEntries(gen("."))) {
		target, err := htmlindex.NewTarget(d)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem: %v\n", err)
		}
		fmt.Fprintf(os.Stdout, "Target: %v\n", target)
		if err := target.Generate(templ); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}

func gen(startPath string) <-chan string {
	out := make(chan string)
	go func() {
		filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				out <- path
			}
			return nil
		})
		close(out)
	}()
	return out
}

func hasEntries(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for path := range in {
			if files, err := ioutil.ReadDir(path); err == nil {
				if len(files) > 0 {
					out <- path
				}
			}
		}
		close(out)
	}()
	return out
}

func noIndex(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for path := range in {
			_, err := os.Stat(filepath.Join(path, "index.html"))
			if os.IsNotExist(err) {
				out <- path
			}
		}
		close(out)
	}()
	return out
}
