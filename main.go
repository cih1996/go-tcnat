package main

import (
	"fmt"
	"log"
	"network_trans/utils"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	args := os.Args
	var mode string
	if len(args) < 2 {
		mode = os.Getenv("MODE")
	} else {
		mode = args[1]
	}

	if mode != "client" && mode != "server" {
		log.Fatal("请指定运行模式: client/server")
		return
	}

	if mode == "client" {
		cfg, err := utils.LoadClientConfig("config/client.json")
		if err != nil {
			log.Fatal(err)
			return
		}
		// 输出读取到的配置
		fmt.Printf("Server Host: %s\n", cfg.Server.Host)
		for _, service := range cfg.List {

			fmt.Printf("Service Tag: %s, Server Port: %d, Local Addr: %s\n",
				service.Tag, service.ServerPort, service.LocalAddr)

			user, err := utils.NewUser(cfg.Server.Host, cfg.Server.MainPort, cfg.Server.TransferPort, service.Tag, service.LocalAddr)
			if err != nil {
				log.Fatalf("无法创建用户实例: %v", err)
			}
			go user.RequestPort(strconv.Itoa(service.ServerPort))
			go user.StartListening()
		}
	} else if mode == "server" {
		server := utils.NewServer()
		//负责接受家庭电脑的控制指令连接
		go server.StartController("8080")
		//负责接受家庭电脑的临时中转连接
		server.StartTransfer("7077")
	}

	select {}

}
