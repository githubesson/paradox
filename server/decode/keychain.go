package decode

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"paradox_server/models"
)

func DecodeKeychain(keychainPath string, password string) (models.DecryptedKeychain, error) {
	outputDir := filepath.Dir(keychainPath)
	jsonOutputPath := filepath.Join(outputDir, "keychain_export.json")

	scriptPath := "chainbreak/chainbreaker.py"
	cmdArgs := []string{
		scriptPath,
		"--export-all",
		"--output", outputDir,
		"--password", password,
		keychainPath,
	}

	exec.Command("python3", cmdArgs...).Run()

	jsonData, err := os.ReadFile(jsonOutputPath)
	if err != nil {
		return models.DecryptedKeychain{}, fmt.Errorf("failed to read expected JSON output file %s: %v", jsonOutputPath, err)
	}

	var decodedMap models.DecryptedKeychain
	if err := json.Unmarshal(jsonData, &decodedMap); err != nil {
		return models.DecryptedKeychain{}, fmt.Errorf("failed to unmarshal JSON from file %s: %v. Content: %s", jsonOutputPath, err, string(jsonData))
	}

	return decodedMap, nil
}
