# PowerShell 脚本版本（使用 curl）

$EXEC_FILE = "tcnat-for-windows.exe"
$DOWNLOAD_URL="http://1.12.233.218:3000/tcnat-for-windows.zip"
$TEMP_FILE = "$env:TEMP\tcnat-for-windows.zip"

# 下载并检查文件
Write-Host "正在从 GitHub 下载程序..."
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TEMP_FILE

# 检查下载是否成功
if (-Not (Test-Path $TEMP_FILE)) {
    Write-Host "下载失败，请检查 URL 或网络连接。"
    exit 1
}

# 解压到当前目录
Write-Host "正在解压文件到当前目录..."
Expand-Archive -Path $TEMP_FILE -DestinationPath . -Force

# 检查解压是否成功
if (-Not (Test-Path $EXEC_FILE)) {
    Write-Host "解压失败，未找到可执行文件。"
    exit 1
}

# 删除临时下载的 ZIP 文件
Remove-Item -Path $TEMP_FILE

# 运行程序
Write-Host "程序解压完成，正在执行..."

# 提示用户选择运行模式
Write-Host "请输入运行模式："
Write-Host "1. 客户端模式(client)"
Write-Host "2. 服务端模式(server)"
$mode = Read-Host "请输入数字"

if ($mode -ne "1" -and $mode -ne "2") {
    Write-Host "无效的输入，请输入 1 或 2。"
    exit 1
}

if ($mode -eq "1") {
    $mode = "client"
    $server_addr = Read-Host "请输入云服务器地址"
} else {
    $mode = "server"
}


# 使用 $PWD.Path 确保路径是正确的
$execPath = Join-Path -Path $PWD.Path -ChildPath $EXEC_FILE


if ($mode -eq "server") {
    Start-Process -FilePath "$execPath" -ArgumentList $mode -NoNewWindow -Wait
} else {
    $arguments = @($mode, $server_addr)
    Start-Process -FilePath "$execPath" -ArgumentList $arguments -NoNewWindow -Wait

}

