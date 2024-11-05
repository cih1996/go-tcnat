package utils

import (
    "encoding/json"
    "os"
)

// 定义与 JSON 结构相对应的 Go 结构体
type Config struct {
    Server ServerConfig `json:"server"`
    List   []Service    `json:"list"`
}

type ServerConfig struct {
    Host string `json:"host"`
}

type Service struct {
    Tag        string `json:"tag"`
    ServerPort int    `json:"server_port"`
    LocalAddr  string `json:"local_addr"`
}

// 读取配置文件并返回 Config 结构体
func LoadConfig(filePath string) (Config, error) {
    var config Config

    // 打开配置文件
    file, err := os.Open(filePath)
    if err != nil {
        return config, err
    }
    defer file.Close()

    // 解析 JSON 文件
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return config, err
    }

    return config, nil
}
