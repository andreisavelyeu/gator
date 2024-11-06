package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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

func ParsePublishedAt(dateString string) (time.Time, error) {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		"2006-01-02",
	}

	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateString)
}
