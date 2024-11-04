package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetConfigFilePath() (string, error) {
	var configBasePath = ".gatorconfig.json"

	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	filePath := filepath.Join(homedir, configBasePath)

	return filePath, nil
}
