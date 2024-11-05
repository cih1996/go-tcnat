#!/bin/bash

# 设置 GitHub 仓库地址
REPO_OWNER="cih1996"
REPO_NAME="go-tcnat"
VERSION="latest"  # 可以指定某个特定版本，比如 v1.0.0

# 构造 GitHub 上文件的 URL
DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/tcnat20241106/tcnat-for-linux.zip"

# 临时下载文件并保存
TEMP_FILE="/tmp/tcnat-for-linux.zip"

# 下载并检查文件
echo "正在从 GitHub 下载程序..."
curl -L "$DOWNLOAD_URL" -o "$TEMP_FILE"

# 检查下载是否成功
if [ ! -f "$TEMP_FILE" ]; then
    echo "下载失败，请检查 URL 或网络连接。"
    exit 1
fi

# 赋予执行权限
chmod +x "$TEMP_FILE"

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

# 执行解压后的程序（假设程序名为 my_program）
echo "程序解压完成，正在执行..."
./tcnat server

echo "程序执行完毕。"