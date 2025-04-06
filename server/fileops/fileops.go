package fileops

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file %s: %w", src, err)
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create base destination directory %s: %w", dest, err)
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", fpath, err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to open file %s for writing: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file %s inside zip: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)

		closeErr1 := outFile.Close()
		closeErr2 := rc.Close()

		if err != nil {
			return fmt.Errorf("failed to copy content to file %s: %w", fpath, err)
		}
		if closeErr1 != nil {
			return fmt.Errorf("failed to close output file %s: %w", fpath, closeErr1)
		}
		if closeErr2 != nil {
			return fmt.Errorf("failed to close zip entry reader for %s: %w", f.Name, closeErr2)
		}
	}

	return nil
}

func FindFile(root, fileName string) (string, error) {
	var foundPath string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == fileName {
			foundPath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if foundPath == "" {
		return "", fmt.Errorf("file '%s' not found in '%s'", fileName, root)
	}
	return foundPath, nil
}

func FindAllFiles(root, fileName string) ([]string, error) {
	var foundPaths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info == nil {
			fmt.Printf("Warning: Received nil FileInfo for path %q\n", path)
			return nil
		}
		if !info.IsDir() && info.Name() == fileName {
			foundPaths = append(foundPaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %q: %w", root, err)
	}

	return foundPaths, nil
}
