package txtarutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/tools/txtar"
)

func ToDir(root string, a *txtar.Archive) error {
	for _, f := range a.Files {
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
			Name: path, // TODO: verify if we don't have to call filepath.ToSlash
			Data: data,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("txtarutil.FromFS: %w", err)
	}
	return a, nil
}
