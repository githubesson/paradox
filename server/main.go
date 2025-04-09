package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"paradox_server/handlers"
	"paradox_server/static"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", static.DBFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS payloads (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "build_id" TEXT UNIQUE,
        "filename" TEXT,
        "timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create payloads table: %w", err)
	}

	log.Println("Database initialized successfully.")
	return nil
}

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("Fatal: Failed to initialize database: %v\n", err)
	}
	defer db.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Paradox Server Running")
	})

	router.POST("/upload", func(c *gin.Context) {
		handlers.HandleUpload(c, db)
	})

	router.GET("/build", func(c *gin.Context) {
		log.Println("Received request to build new payload...")
		buildID, filename, err := handlers.BuildPayload(db)
		if err != nil {
			log.Printf("Error building payload: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to build payload",
				"details": err.Error(),
			})
			return
		}
		log.Printf("Successfully built new payload: BuildID=%s, Filename=%s\n", buildID, filename)
		c.JSON(http.StatusOK, gin.H{
			"message":  "Payload built successfully",
			"build_id": buildID,
			"filename": filename,
		})
	})

	router.GET("/download/build/:buildid", func(c *gin.Context) {
		buildID := c.Param("buildid")
		log.Printf("Received request to download build ID: %s\n", buildID)

		filename, err := handlers.GetFilenameByBuildID(db, buildID)
		if err != nil {
			if err.Error() == fmt.Sprintf("build ID '%s' not found in database", buildID) {
				log.Printf("Build ID '%s' not found for download.\n", buildID)
				c.JSON(http.StatusNotFound, gin.H{"error": "Build ID not found"})
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

		log.Printf("Serving file '%s' for build ID '%s'\n", filePath, buildID)
		c.FileAttachment(filePath, filename)
	})

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Printf("Fatal: Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
