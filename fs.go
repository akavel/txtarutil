// Package txtarutil provides utility functions for the txtar trivial
// text-based file archive format.
//
// The functions provided by the package allow easy reading and writing of
// [txtar] archive contents from/to disk.
package txtarutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/tools/txtar"
)

// ToDir writes the files specified in the archive to directory at root path.
// Directories are created if needed. Existing files may get overwritten.
// Checks file names using [filepath.IsLocal], returning an error if not passed.
func ToDir(root string, archive *txtar.Archive) error {
	for _, f := range archive.Files {
		if !filepath.IsLocal(f.Name) {
			return fmt.Errorf("txtarutil.ToDir: file name is not a valid local path on current OS: %s", f.Name)
		}
		path := filepath.Join(root, f.Name)
		// Create file's parent dir
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return fmt.Errorf("txtarutil.ToDir: %w", err)
		}
		// Write file contents
		err = os.WriteFile(path, f.Data, 0666)
		if err != nil {
			return fmt.Errorf("txtarutil.ToDir: %w", err)
		}
	}
	return nil
}

// FromFS walks the filesystem tree at root, adding any files and their contents
// into a [txtar.Archive].
func FromFS(root fs.FS) (*txtar.Archive, error) {
	a := &txtar.Archive{}
	err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(root, path)
		if err != nil {
			return err
		}
		a.Files = append(a.Files, txtar.File{
			Name: path,
			Data: data,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("txtarutil.FromFS: %w", err)
	}
	return a, nil
}
