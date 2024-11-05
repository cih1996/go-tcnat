package utils

import (
	"encoding/json"
	"os"
)

type ClientServerConfig struct {
	Host         string `json:"host"`
	MainPort     int    `json:"main_port"`
	TransferPort int    `json:"transfer_port"`
}

type ClientList struct {
	Tag        string `json:"tag"`
	ServerPort int    `json:"server_port"`
	LocalAddr  string `json:"local_addr"`
}

type ClientConfig struct {
	Server ClientServerConfig `json:"server"`
	List   []ClientList       `json:"list"`
}

// 读取client配置
func LoadClientConfig(filePath string) (ClientConfig, error) {
	var config ClientConfig

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

// 读取server配置
