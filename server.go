package main

import (
    "fmt"
    "log"
    "math/rand"
    "time"
    "net"
    "io"
    "sync"
)

// 存储家庭电脑的信息
type Client struct {
    tag_id      string
    conn        net.Conn
    controlCh chan net.Conn
}

type Trans struct {
    tag_id      string
    conn1       net.Conn
    conn2       net.Conn
    cond    *sync.Cond
    mu      sync.Mutex
}

// 创建结构体
type Server struct {
    mu    sync.Mutex
    clients map[string]*Client
    trans map[string]*Trans
}

// 创建服务器
func NewServer() *Server {
    return &Server{
        mu: sync.Mutex{},
        clients: make(map[string]*Client),
        trans: make(map[string]*Trans),
    }
}



func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    b := make([]byte, length)
    rand.Seed(time.Now().UnixNano())
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }

    return string(b)
}


// 开始监听端口,接受家庭电脑的临时中转连接
func (s *Server) StartTransfer(port string) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("【中转】开启中转端口失败 %s: %v", port, err)
    }
    defer listener.Close()
    log.Printf("【中转】端口监听成功 %s", port)
    //接受来自家庭电脑的临时中转连接
    for {
        log.Printf("【中转】等待中转临时用户")
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        log.Printf("【中转】有新的连接请求到端口:%s 客户:%s", port, conn.RemoteAddr().String())
        go s.handleTransfer(conn)
    }
    log.Printf("【中转】端口监听结束 %s", port)
}


func (s *Server) handleTransfer(conn net.Conn) {
    var clientAddr string
    //获取家庭电脑的地址和端口
    clientAddr = conn.RemoteAddr().String()
    log.Printf("【中转】来自家庭电脑的中转连接: %s", clientAddr)
    var client_tag string
    if _, err := fmt.Fscan(conn, &client_tag); err != nil {
        log.Printf("【中转】读取数据失败: %v", err)
        return
    }
    log.Printf("【中转】对应客户端的标记符: %s", client_tag)
    s.mu.Lock()
    trans:= s.trans[client_tag]
    s.mu.Unlock()
    if(trans != nil){
        s.mu.Lock()
        s.trans[client_tag].conn1 = conn
        s.mu.Unlock()
        log.Printf("【中转】与外部客户绑定成功，进入交互阶段，标记符: %s", client_tag)
       
        go func(){
            var wg sync.WaitGroup

            wg.Add(2) // 等待两个 goroutine
            go s.handleUserToClientTransData(client_tag, &wg)
            go s.handleClientToUserTransData(client_tag, &wg)
          
            wg.Wait() // 等待两个 goroutine 完成

            
            //log.Printf("【交互】(%s:%s-%s)结束连接", tag_id, port,client_tag)
            s.mu.Lock()
            delete(s.trans, client_tag)
            s.mu.Unlock()
        }()
    }

}


// 开始监听端口,给家庭电脑建立通信连接
func (s *Server) StartController(port string) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Fatalf("【控制】监听控制端口失败 %s: %v", port, err)
    }
    defer listener.Close()

    log.Printf("【控制】监听控制端口成功 %s", port)

    //接受来自家庭电脑的连接
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("【控制】服务失败: %v", err)
            continue
        }

        go s.handleControlConnection(conn)
    }
}

// 处理家庭电脑的连接
func (s *Server) handleControlConnection(conn net.Conn) {
    defer conn.Close()
    var clientAddr string
    var tag_id string
    //对外端口
    var port string
    clientAddr = conn.RemoteAddr().String()
    //对外tcp
    var listener net.Listener

    //定时发送心跳包
    go func(){
        for {
            time.Sleep(30 * time.Second)
            //判断conn是否已经关闭
            if conn == nil {
                break
            }
            
            fmt.Fprintln(conn, "pong")
       
        }
    }()

    for {
        var command string
        var err error
        if conn == nil {
            break
        }

        if _, err := fmt.Fscan(conn, &command); err != nil {
            break
        }

        if command == "request" {
            if _, err := fmt.Fscan(conn, &port); err != nil {
                break
            }
            if _, err := fmt.Fscan(conn, &tag_id); err != nil {
                break
            }

            //判断tag_id的s.clients是否存在
            if s.clients[tag_id] != nil {
                fmt.Fprintln(conn,"request-tag-repeat")
            }else{
                s.clients[tag_id] = &Client{
                    tag_id:      tag_id,
                    conn:        conn,
                    controlCh: make(chan net.Conn),
                }
                listener, err = net.Listen("tcp", ":"+port)
                if err != nil {
                    log.Printf("【控制】(%s)内网设备对外端口监听失败: %s: %v", tag_id, port,err)
                    fmt.Fprintln(conn, "request-failed")
                }else{
                    fmt.Fprintln(conn, "request-success",port)
                    log.Printf("【控制】(%s)内网设备对外端口监听: %s:", tag_id, port)
                    
                    go s.handleNewConnection(listener,tag_id, port)
                    log.Printf("【控制】(%s)接入新内网设备: %s", tag_id, clientAddr)
                }
                
            }

            
            
            
           
        }

    }
    if conn != nil {
        conn.Close()
        conn = nil
    }
   

    delete(s.clients, tag_id)
    log.Printf("【控制】内网设备已端口，销毁对应端口%s : %s",port, clientAddr)
    //关闭监听
    if listener != nil {
        listener.Close()
    }
    
}

// 开始监听对外端口
func (s *Server) handleNewConnection(listener net.Listener,tag_id string, port string) {
    defer listener.Close()
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
         // 生成随机字母
         client_tag := generateRandomString(5)
         s.trans[client_tag] = &Trans{
            tag_id:      tag_id,
            conn1:        nil,
            conn2:      conn,
        }
        s.trans[client_tag].cond = sync.NewCond(&sync.Mutex{})

        log.Printf("【通信】(%s:%s-%s)收到新的客户请求，分配关系tag，通知内网设备建立中转连接", tag_id, port,client_tag)
        fmt.Fprintln(s.clients[tag_id].conn, "client-connect", client_tag)
   

    }
}

//接受外部客户端的数据，然后转发给家庭电脑
func (s *Server) handleClientToUserTransData(client_tag string, wg *sync.WaitGroup){
    defer wg.Done() // 完成时调用 Done
    conn1 := s.trans[client_tag].conn1
    conn2 := s.trans[client_tag].conn2
    if conn1 == nil {
        log.Printf("【错误】(%s) conn1 is nil", client_tag)
        return
    }
    if conn2 == nil {
        log.Printf("【错误】(%s) conn2 is nil", client_tag)
        return
    }
    buffer := make([]byte, 10240) // 缓冲区
    log.Printf("【中转】(%s)开始接收客户端数据", client_tag)
    for {
        n, err := conn2.Read(buffer)
        if err != nil {
            if err != io.EOF {
            }
            break
        }
        log.Printf("【中转】(%s)客户端(%s)-->内网设备(%d)", client_tag, n,conn1.RemoteAddr().String())
        conn1.Write(buffer[:n])
    }
    if conn1 != nil {
        conn1.Close()
    }
    if conn2 != nil {
        conn2.Close()
    }
}

//接收家庭电脑的数据，转发给客户端
func (s *Server) handleUserToClientTransData(client_tag string, wg *sync.WaitGroup){
    defer wg.Done() // 完成时调用 Done
    conn1 := s.trans[client_tag].conn1
    conn2 := s.trans[client_tag].conn2
    if conn1 == nil {
        log.Printf("【错误】(%s) conn1 is nil", client_tag)
        return
    }
    if conn2 == nil {
        log.Printf("【错误】(%s) conn2 is nil", client_tag)
        return
    }
    buffer := make([]byte, 10240) // 缓冲区
    for {
        n, err := conn1.Read(buffer)
        if err != nil {
            if err != io.EOF {
            }
            break
        }
        log.Printf("【中转】(%s)内网设备(%d)-->客户端", client_tag, n)
        conn2.Write(buffer[:n])
    }
    if conn1 != nil {
        conn1.Close()
    }
    if conn2 != nil {
        conn2.Close()
    }
}


func main() {
    server := NewServer()
    //负责接受家庭电脑的控制指令连接
    go server.StartController("8080")
    //负责接受家庭电脑的临时中转连接
    server.StartTransfer("7077")
}
