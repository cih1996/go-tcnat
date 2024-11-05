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
#通过发布程序
./main server
#通过源码
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

```
启动服务（通过发布程序）
```
./main client
```
启动服务（通过源码）
```
go run main.go client
```
> 也可以通过修改.env的MODE来控制选择client/server，则无需传递参数，直接执行
```
#源码
go run main.go
#发布版
./main
```

这样本地电脑将会连接到云服务器，云服务器将会监听 3333 端口，当有用户连接 3333 端口，那么云服务器将全部转发给家庭电脑，家庭电脑在连接到本机的 127.0.0.1:3389 端口，实现外网访问内网服务



# 通过源码运行
如果你的云服务器也是windows，可以自行编译go源码，通过以下命令打包
```
go build -o server.exe server.go
```
如果希望直接使用源码运行，首先先安装好go的运行环境，可前往官网下载安装：
https://go.dev/doc/install

云服务器运行命令：
```
go run mian.go server
```

家庭电脑运行命令：
```
go run main.go client
```


# 升级记录
- 2024.11.6
1. 在操作map时增加锁防止死锁保持
2. 优化结构，将客户端和服务端合并
3. 新增配置文件.env用于控制启动模式(client/server)
4. 新增可通过命令行来切换模式，在无参数的情况下默认使用.env
```
./main server   #以客户端模式运行
./main client   #以服务端模式运行
```

- 2024.11.5
1. 初次发布，构建最基本的转发功能

