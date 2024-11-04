package config

import (
	"encoding/json"
	"fmt"
	"gator/internal/utils"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) {
	c.Current_user_name = username

	filePath, err := utils.GetConfigFilePath()

	if err != nil {
		return
	}

	config, err := json.Marshal(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(filePath, config, os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}
}

func Read() Config {
	var config Config
	filePath, err := utils.GetConfigFilePath()

	if err != nil {
		fmt.Println(err)
		return config
	}

	fileData, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println(err)
		return config
	}

	err = json.Unmarshal(fileData, &config)

	if err != nil {
		fmt.Println(err)
		return config
	}

	return config
}
