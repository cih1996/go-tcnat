package utils

import (
    "fmt"
    "io"
    "log"
    "net"
)

type User struct {
    controlConn      net.Conn
    localServerAddr  string
    serverAddr       string
	tagID			 string
}

// NewUser creates a new User instance
func NewUser(serverAddr,tagID string, localServerAddr string) (*User, error) {
    log.Printf("【%s】【通信】正在接入中转服务器:%s", tagID, serverAddr)
    conn, err := net.Dial("tcp", serverAddr+":8080")
    if err != nil {
        return nil, fmt.Errorf("连接中转服务器失败: %v", err)
    }

    return &User{
        controlConn:     conn,
        localServerAddr: localServerAddr,
        serverAddr:      serverAddr,
		tagID:			 tagID,
    }, nil
}

// RequestPort requests a specific port for data transfer
func (u *User) RequestPort(port string) {
    fmt.Fprintln(u.controlConn, "request", port, u.tagID)
}

// StartListening starts handling incoming connections
func (u *User) StartListening() {
    u.handleIncomingConnections()
}

// 接收云服务器的指令
func (u *User) handleIncomingConnections() {
    for {
        var command string
        if _, err := fmt.Fscan(u.controlConn, &command); err != nil {
            log.Printf("【%s】【通信】读取数据失败: %v", u.tagID, err)
            return
        }

        // 有新的连接请求
        if command == "client-connect" {
            log.Printf("【%s】【通信】有新的连接请求", u.tagID)
            var clientTag string
            if _, err := fmt.Fscan(u.controlConn, &clientTag); err != nil {
                log.Printf("【%s】【通信】无法读取标记: %v", u.tagID, err)
                break
            }
            go u.handleTransfer(clientTag)
        }

        if command == "request-tag-repeat" {
            log.Printf("【%s】【通信】建立失败，原因：自定义标识符重复，请更换!", u.tagID)
        }

        if command == "pong" {
            log.Printf("【%s】【通信】收到pong指令", u.tagID)
            fmt.Fprintln(u.controlConn, "pong")
        }

        if command == "request-success" {
            var port string
            if _, err := fmt.Fscan(u.controlConn, &port); err != nil {
                log.Printf("【%s】【通信】无法读取外网端口: %v", u.tagID, err)
                break
            }
            log.Printf("【%s】【通信】接入中转服务器成功，外网地址为: %s:%s", u.tagID, u.serverAddr, port)
        }

        if command == "request-failed" {
            log.Printf("【%s】【通信】无法接入中转服务器，可能云服务器端口已被占用！", u.tagID)
        }
    }
}

//建立中转临时连接
func (u *User) handleTransfer(clientTag string) {
    log.Printf("【%s】【中转】建立临时中转连接，标记:%s", u.tagID, clientTag)
    conn, err := net.Dial("tcp", u.serverAddr+":7077")
    if err != nil {
        log.Printf("【%s】【中转】连接到中转服务器失败: %v", u.tagID, err)
        return
    }
    fmt.Fprintln(conn, clientTag)
    
    connServer, err := net.Dial("tcp", u.localServerAddr)
    if err != nil {
        log.Printf("【%s】【中转】连接到站点服务器失败: %v", u.tagID, err)
        return
    }
    
    log.Printf("【%s】【中转】已经连接站点服务器:%s", u.tagID, u.localServerAddr)

    // 接收数据
    go func() {
        buffer := make([]byte, 10240) // 缓冲区
        for {
            n, err := conn.Read(buffer)
            if err != nil {
                if err != io.EOF {
                    log.Printf("【%s】读取连接数据失败: %v\n", u.tagID, err)
                }
                log.Printf("【%s】接收中转数据结束: %v\n", u.tagID, err)
                break
            }
            log.Printf("【%s】【中转】(%s)中转(%d)-->本地", u.tagID, clientTag, n)
            connServer.Write(buffer[:n])
        }
        conn.Close()
    }()

    // 接收本地服务器的数据，并转发到中转服务器
    go func() {
        buffer := make([]byte, 10240) // 缓冲区
        for {
            n, err := connServer.Read(buffer)
            if err != nil {
                if err != io.EOF {
                    log.Printf("【%s】读取本地服务器数据失败: %v\n", u.tagID, err)
                }
                log.Printf("【%s】接收本地服务器数据结束: %v\n", u.tagID, err)
                break
            }
            log.Printf("【%s】【中转】(%s)本地(%d)-->中转", u.tagID, clientTag, n)
            conn.Write(buffer[:n])
        }
        connServer.Close()
    }()
}
