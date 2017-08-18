// Package proc provides repository processing.
package proc // import "docc.io/source/proc"

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"docc.io/source"
	"docc.io/source/parse"
)

var (
	// DirSizeMax is the upper limit for the number of files.
	// BUG(pascaldekloe): DirSizeMax not enforced.
	DirSizeMax = 99

	// FileSizeMax is the upper limit for the number of bytes.
	FileSizeMax int64 = 1024 * 1024
)

// Repo parses all supported files from a directory tree.
// The error is for filesystem issues. All I/O runs in calling routine.
// BUG(pascaldekloe): Should skip hidden files.
func Repo(root string) (<-chan *source.File, error) {
	files := make(chan *source.File, 9)
	var wg sync.WaitGroup

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// filter C files
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".c" && ext != ".h" {
			return nil
		}

		var f source.File
		f.Path, err = filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if info.Size() > FileSizeMax {
			f.Issues = "File too big."

			// don't block
			wg.Add(1)
			go func(entry *source.File) {
				files <- entry
				wg.Done()
			}(&f)

			return nil
		}

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(entry *source.File, content []byte) {
			parseC(entry, content)

			files <- entry
			wg.Done()
		}(&f, bytes)

		return nil
	})

	// close on last
	go func() {
		wg.Wait()
		close(files)
	}()

	if err != nil {
		// flush to unblock pending routines
		go func() {
			for range files {
			}
		}()

		return nil, err
	}

	return files, nil
}

func parseC(f *source.File, content []byte) {
	decls := make(chan *source.Decl)
	go parse.C(content, decls)

	for d := range decls {
		// BUG(pascaldekloe): Highly incomplete declaration recognition.
		if strings.IndexByte(d.Source, '(') >= 0 {
			f.Funcs = append(f.Funcs, d)
		} else {
			f.Vars = append(f.Vars, d)
		}
	}
}
