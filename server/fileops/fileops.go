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

func ZipDirectory(sourceDir, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file %s: %w", zipFilePath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	sourceDirAbs, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of %s: %w", sourceDir, err)
	}

	err = filepath.Walk(sourceDirAbs, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking %s: %w", filePath, err)
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return fmt.Errorf("failed to create header for %s: %w", filePath, err)
		}

		relPath, err := filepath.Rel(sourceDirAbs, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", filePath, err)
		}

		if fileInfo.IsDir() {
			header.Name = relPath + "/"
		} else {
			header.Name = relPath
		}

		header.Method = zip.Deflate

		if fileInfo.IsDir() {
			_, err = zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("failed to create directory header for %s: %w", relPath, err)
			}
		} else {
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("failed to create file header for %s: %w", relPath, err)
			}

			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", filePath, err)
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("failed to write file %s to zip: %w", filePath, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating zip archive: %w", err)
	}

	return nil
}
