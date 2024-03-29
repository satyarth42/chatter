package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Host string    `json:"host"`
	Port string    `json:"port"`
	DB   MySQLConf `json:"user_db"`
}

type MySQLConf struct {
	User     string `json:"user"`
	Password string `json:"pass"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

const configFilePath = "CONF_PATH"

var conf Config

func LoadConfig() {
	filePath := os.Getenv(configFilePath)
	if filePath == "" {
		filePath, _ = filepath.Abs("config.json")
	}

	bytes, readErr := os.ReadFile(filePath)
	if readErr != nil {
		log.Fatalf("failed to read file: %s, err:%+v", filePath, readErr)
	}

	config := Config{}
	unmarshalErr := json.Unmarshal(bytes, &config)
	if unmarshalErr != nil {
		log.Fatalf("failed to unmarshal config file bytes: %s, err: %+v", string(bytes), unmarshalErr)
	}

	conf = config
}

func GetConfig() Config {
	return conf
}
