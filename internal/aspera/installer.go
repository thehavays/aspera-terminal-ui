package aspera

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed assets/ascp
var ascpBinary []byte

//go:embed assets/asperaweb_id_dsa.openssh
var ascpKey []byte

//go:embed assets/aspera-license
var ascpLicense []byte

// GetAscpPaths ensures that the embedded ascp binary and ssh key are present
// in the ~/.atui/bin directory and returns their absolute paths.
func GetAscpPaths() (string, string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", err
	}

	atuiBinDir := filepath.Join(configDir, "atui", "bin")
	atuiEtcDir := filepath.Join(configDir, "atui", "etc")
	if err := os.MkdirAll(atuiBinDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create atui bin directory: %w", err)
	}
	if err := os.MkdirAll(atuiEtcDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create atui etc directory: %w", err)
	}

	ascpPath := filepath.Join(atuiBinDir, "ascp")
	keyPath := filepath.Join(atuiEtcDir, "asperaweb_id_dsa.openssh")
	licensePath := filepath.Join(atuiEtcDir, "aspera-license")

	// Check if ascp is already extracted and has the right size/permissions
	// For simplicity, we check if they exist. We could also check hash if needed.
	ascpStat, errAscp := os.Stat(ascpPath)
	keyStat, errKey := os.Stat(keyPath)
	licenseStat, errLicense := os.Stat(licensePath)

	if errAscp != nil || ascpStat.Size() != int64(len(ascpBinary)) {
		fmt.Println("Extracting embedded ascp binary to", ascpPath)
		if err := os.WriteFile(ascpPath, ascpBinary, 0755); err != nil {
			return "", "", fmt.Errorf("failed to write ascp binary: %w", err)
		}
	}

	if errKey != nil || keyStat.Size() != int64(len(ascpKey)) {
		fmt.Println("Extracting embedded ascp ssh key to", keyPath)
		if err := os.WriteFile(keyPath, ascpKey, 0600); err != nil {
			return "", "", fmt.Errorf("failed to write ascp ssh key: %w", err)
		}
	}

	if errLicense != nil || licenseStat.Size() != int64(len(ascpLicense)) {
		fmt.Println("Extracting embedded aspera license to", licensePath)
		if err := os.WriteFile(licensePath, ascpLicense, 0644); err != nil {
			return "", "", fmt.Errorf("failed to write aspera license: %w", err)
		}
	}

	return ascpPath, keyPath, nil
}
