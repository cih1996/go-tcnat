让个人电脑秒变云服务器

准备材料：
1台云服务器
1台家庭电脑


# 云服务器部署流程：

#执行服务端，默认端口是（8080和7077）
```
./server
```

# 家庭电脑部署流程：
打开client.json文件，按如下配置说明
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
示例数据：
```
{
    "server": {
        "host": "1.12.233.218"
    },
    "list": [{
        "tag": "bt-login",
        "server_port": 3333,
        "local_addr": "127.0.0.1:3389"
    }
}
```

运行 main.exe 

这样本地电脑将会连接到云服务器，云服务器将会监听 3333 端口，当有用户连接 3333 端口，那么云服务器将全部转发给家庭电脑，家庭电脑在连接到本机的 127.0.0.1:3389 端口，实现外网访问内网服务



# 通过源码运行
如果你的云服务器也是windows，可以自行编译go源码，通过以下命令打包
go build -o server.exe server.go

如果希望直接使用源码运行，首先先安装好go的运行环境，可前往官网下载安装：
https://go.dev/doc/install

云服务器运行命令：
go run server.go

家庭电脑运行命令：
go run client.go