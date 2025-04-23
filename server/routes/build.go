package routes

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"paradox_server/handlers"
	"paradox_server/static"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SetupBuildRoutes(router *gin.RouterGroup, db *sql.DB) {
	router.GET("/build", func(c *gin.Context) {
		handleBuildRequest(c, db)
	})

	router.GET("/builds", func(c *gin.Context) {
		handleListBuilds(c, db)
	})

	router.GET("/download/build/:buildid", func(c *gin.Context) {
		handleDownloadBuild(c, db)
	})
}

func handleBuildRequest(c *gin.Context, db *sql.DB) {
	log.Println("Received request to build new payload...")

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

	buildID, filename, err := handlers.BuildPayload(db)
	if err != nil {
		log.Printf("Error building payload: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to build payload",
			"details": err.Error(),
		})
		return
	}

	_, err = db.Exec("INSERT INTO payloads (build_id, filename, user_id) VALUES (?, ?, ?)",
		buildID, filename, userID)
	if err != nil {
		log.Printf("Error saving build to database: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save build",
		})
		return
	}

	log.Printf("Successfully built new payload: BuildID=%s, Filename=%s, UserID=%d\n", buildID, filename, userID)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Payload built successfully",
		"build_id": buildID,
		"filename": filename,
	})
}

func handleListBuilds(c *gin.Context, db *sql.DB) {
	log.Println("Received request to list all builds")

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
		SELECT build_id, filename, timestamp 
		FROM payloads 
		WHERE user_id = ?
		ORDER BY timestamp DESC
	`, userID)
	if err != nil {
		log.Printf("Database error fetching builds: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve builds"})
		return
	}
	defer rows.Close()

	type BuildEntry struct {
		BuildID   string `json:"build_id"`
		Filename  string `json:"filename"`
		Timestamp string `json:"timestamp"`
	}

	var builds []BuildEntry

	for rows.Next() {
		var entry BuildEntry
		err := rows.Scan(
			&entry.BuildID,
			&entry.Filename,
			&entry.Timestamp,
		)
		if err != nil {
			log.Printf("Error scanning build row: %v\n", err)
			continue
		}
		builds = append(builds, entry)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating build rows: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while retrieving builds"})
		return
	}

	log.Printf("Found %d builds for user %s\n", len(builds), username)
	c.JSON(http.StatusOK, gin.H{
		"builds": builds,
		"count":  len(builds),
	})
}

func handleDownloadBuild(c *gin.Context, db *sql.DB) {
	buildID := c.Param("buildid")
	username := c.GetString("username")
	log.Printf("Received request to download build ID: %s from user: %s\n", buildID, username)

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		log.Printf("Error getting user ID: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user ID",
		})
		return
	}

	var filename string
	err = db.QueryRow("SELECT filename FROM payloads WHERE build_id = ? AND user_id = ?", buildID, userID).Scan(&filename)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Build ID '%s' not found or not owned by user '%s'\n", buildID, username)
			c.JSON(http.StatusNotFound, gin.H{"error": "Build not found or access denied"})
		} else {
			log.Printf("Database error fetching filename for build ID '%s': %v\n", buildID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	filePath := filepath.Join(static.PayloadOutputDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found for build ID '%s' at path '%s'\n", buildID, filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Payload file not found on server"})
		return
	}

	log.Printf("Serving file '%s' for build ID '%s' to user '%s'\n", filePath, buildID, username)
	c.FileAttachment(filePath, filename)
}
