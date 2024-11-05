只需要一个低廉的云服务器作为桥梁，让个人高配置电脑秒变云服务器，特别适合需要高运算需求（如AI服务，高并发运算）而预算较低的朋友
基于GO语言原生开发，极大的提高多线程及内存管理效率，无延迟无卡顿，体验感极高！

# 云服务器部署流程
打开config/server.json文件，按如下配置说明
```
{
    "server": {
        "main_port": 控制端口,
        "transfer_port":临时交互端口
    }
}
```

[示例](config/server.json)

启动服务


```
#通过发布程序运行
./main server
#通过源码运行
go run main.go server
```
> 也可以通过修改.env的MODE来控制选择client/server，则无需传递参数，直接执行

# 家庭电脑部署流程
打开config/client.json文件，按如下配置说明
```
{
    "server": {
        "host": "云服务器地址"
    },
    "list": [{
        "tag": "bt-login",
        "server_port": 云服务器端口,
        "local_addr": "127.0.0.1:本机服务端口"
    }
}
```
[示例](config/client.json)


启动服务
```
#通过发布程序运行
./main client
#通过源码运行
go run main.go client
```
> 也可以通过修改.env的MODE来控制选择client/server，则无需传递参数，直接执行


这样本地电脑将会连接到云服务器，云服务器将会监听 3333 端口，当有用户连接 3333 端口，那么云服务器将全部转发给家庭电脑，家庭电脑在连接到本机的 127.0.0.1:3389 端口，实现外网访问内网服务



# 编译发布
安装库之前切换到国内源，否则容易安装失败
- windows
```
export GOPROXY=https://goproxy.cn,direct

```
- linux
```
set GOPROXY=https://goproxy.cn,direct
```
安装成功后执行编译
```
go get github.com/joho/godotenv
go build -o tcnat.exe main.go
```

