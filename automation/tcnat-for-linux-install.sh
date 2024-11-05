#!/bin/bash

EXEC_FILE="tcnat-for-linux"
DOWNLOAD_URL="http://1.12.233.218:3000/tcnat-for-linux.zip"
TEMP_FILE="/tmp/tcnat-for-linux.zip"

# 下载并检查文件
echo "正在从 GitHub 下载程序..."
curl -L "$DOWNLOAD_URL" -o "$TEMP_FILE"

# 检查下载是否成功
if [ ! -f "$TEMP_FILE" ]; then
    echo "下载失败，请检查 URL 或网络连接。"
    exit 1
fi

# 解压到当前目录
echo "正在解压文件到当前目录..."
unzip -o "$TEMP_FILE" -d .  # 解压到当前目录（-o 表示覆盖已存在文件）

# 检查解压是否成功
if [ $? -ne 0 ]; then
    echo "解压失败。"
    exit 1
fi

# 删除临时下载的 ZIP 文件
rm "$TEMP_FILE"
chmod +x "$EXEC_FILE"

# 执行解压后的程序（假设程序名为 my_program）
echo "程序解压完成，正在执行..."
# 要求用户输入以什么模式运行
echo "请输入运行模式："
echo "1. 客户端模式(client)"
echo "2. 服务端模式(server)"
read -p "请输入数字: " mode

# 检查输入是否有效
if [ "$mode" -ne 1 ] && [ "$mode" -ne 2 ]; then
    echo "输入错误，程序退出。"
    exit 1
fi

# 设置模式
if [ "$mode" -eq 1 ]; then
    mode="client"
else
    mode="server"
fi

# 客户端模式
if [ "$mode" == "client" ]; then
    # 请求用户输入服务器地址
    read -p "请输入服务器地址: " server
    # 运行 tcnat 客户端
    ./"$EXEC_FILE" "$mode" "$server" "$port"
    exit 0
else
    # 服务端模式
    echo "正在以服务端模式启动..."
    # 运行 tcnat 服务端
    ./"$EXEC_FILE" "$mode"
fi
