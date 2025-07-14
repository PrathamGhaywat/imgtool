package utils

import (
	"os"
	"path/filepath"
	"strings"
)

var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// FindImageFiles searches for image files in a directory.
func FindImageFiles(root string, recursive bool) ([]string, error) {
	var files []string

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if imageExtensions[strings.ToLower(filepath.Ext(path))] {
				files = append(files, path)
			}
		} else if !recursive && path != root {
			return filepath.SkipDir
		}

		return nil
	}

	err := filepath.Walk(root, walkFn)
	return files, err
}
