package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"paradox_server/static"
	"path/filepath"
)

func StoreBuildInfo(db *sql.DB, buildID, filename string) error {
	insertSQL := `INSERT INTO payloads (build_id, filename) VALUES (?, ?)`
	_, err := db.Exec(insertSQL, buildID, filename)
	if err != nil {
		return fmt.Errorf("failed to insert build info (BuildID: %s, Filename: %s): %w", buildID, filename, err)
	}
	log.Printf("Stored build info: BuildID=%s, Filename=%s\n", buildID, filename)
	return nil
}

func GenerateRandomName(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func BuildPayload(db *sql.DB) (string, string, error) {
	buildID, err := GenerateRandomName(16)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate build ID: %w", err)
	}

	payloadFilename, err := GenerateRandomName(8)
	if err != nil {
		return buildID, "", fmt.Errorf("failed to generate payload filename: %w", err)
	}

	if err := os.MkdirAll(static.PayloadOutputDir, 0755); err != nil {
		return buildID, payloadFilename, fmt.Errorf("failed to create payload output directory '%s': %w", static.PayloadOutputDir, err)
	}

	outputPath := filepath.Join(static.PayloadOutputDir, payloadFilename)

	ldflags := fmt.Sprintf("-X main.BuildID=%s -w -s -buildid=", buildID)

	cmd := exec.Command("go", "build", "-ldflags="+ldflags, "-trimpath", "-o", outputPath)
	cmd.Dir = static.PayloadSourceDir

	log.Printf("Building payload with BuildID: %s -> %s\n", buildID, outputPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return buildID, payloadFilename, fmt.Errorf("payload build failed: %w", err)
	}

	log.Printf("Payload built successfully: %s\n", outputPath)

	if err := StoreBuildInfo(db, buildID, payloadFilename); err != nil {
		log.Printf("Error storing build info after successful build: %v\n", err)
	}

	return buildID, payloadFilename, nil
}

func GetFilenameByBuildID(db *sql.DB, buildID string) (string, error) {
	querySQL := `SELECT filename FROM payloads WHERE build_id = ? LIMIT 1`
	var filename string
	err := db.QueryRow(querySQL, buildID).Scan(&filename)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("build ID '%s' not found in database", buildID)
		}

		return "", fmt.Errorf("database error querying for build ID '%s': %w", buildID, err)
	}
	return filename, nil
}
