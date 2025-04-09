package decode

import (
	"database/sql"
	"paradox_server/models"
	"paradox_server/static"
)

// not implemented

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
