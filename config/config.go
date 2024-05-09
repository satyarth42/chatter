package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Host  string      `json:"host"`
	Port  string      `json:"port"`
	DB    MySQLConf   `json:"user_db"`
	Redis RedisConfig `json:"redis"`
	Etcd  EtcdConf    `json:"etcd"`
}

type MySQLConf struct {
	User     string `json:"user"`
	Password string `json:"pass"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

type RedisConfig struct {
	Addresses            []string `json:"addresses"`
	ReadOnly             bool     `json:"readOnly"`
	DialTimeoutInSec     int      `json:"dialTimeoutInSec"`
	ReadTimeoutInSec     int      `json:"readTimeoutInSec"`
	WriteTimeoutInSec    int      `json:"writeTimeoutInSex"`
	PoolSize             int      `json:"poolSize"`
	MinIdleConns         int      `json:"minIdleConns"`
	ConnMaxIdleTimeInSec int      `json:"connMaxIdleTimeInSec"`
}

type EtcdConf struct {
	Endpoints        []string `json:"endpoints"`
	DialTimeoutInSec int      `json:"dialTimeoutInSec"`
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
