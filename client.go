package main

import (
	"network_trans/utils"
	"log"
	"fmt"
	"strconv"
)


func main() {
	
	// 读取配置文件
	cfg, err :=  utils.LoadConfig("client.json")
	if err != nil {
		log.Fatal(err)
	}

	// 输出读取到的配置
	fmt.Printf("Server Host: %s\n", cfg.Server.Host)
	for _, service := range cfg.List {

		fmt.Printf("Service Tag: %s, Server Port: %d, Local Addr: %s\n",
			service.Tag, service.ServerPort, service.LocalAddr)

		user, err := utils.NewUser(cfg.Server.Host, service.Tag, service.LocalAddr)
		if err != nil {
			log.Fatalf("无法创建用户实例: %v", err)
		}
		go user.RequestPort(strconv.Itoa(service.ServerPort))
		go user.StartListening()
	}

	select {}


}
