# PowerShell �ű��汾��ʹ�� curl��

$EXEC_FILE = "tcnat-for-windows.exe"
$DOWNLOAD_URL="http://1.12.233.218:3000/tcnat-for-windows.zip"
$TEMP_FILE = "$env:TEMP\tcnat-for-windows.zip"

# ���ز�����ļ�
Write-Host "���ڴ� GitHub ���س���..."
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TEMP_FILE

# ��������Ƿ�ɹ�
if (-Not (Test-Path $TEMP_FILE)) {
    Write-Host "����ʧ�ܣ����� URL ���������ӡ�"
    exit 1
}

# ��ѹ����ǰĿ¼
Write-Host "���ڽ�ѹ�ļ�����ǰĿ¼..."
Expand-Archive -Path $TEMP_FILE -DestinationPath . -Force

# ����ѹ�Ƿ�ɹ�
if (-Not (Test-Path $EXEC_FILE)) {
    Write-Host "��ѹʧ�ܣ�δ�ҵ���ִ���ļ���"
    exit 1
}

# ɾ����ʱ���ص� ZIP �ļ�
Remove-Item -Path $TEMP_FILE

# ���г���
Write-Host "�����ѹ��ɣ�����ִ��..."

# ��ʾ�û�ѡ������ģʽ
Write-Host "����������ģʽ��"
Write-Host "1. �ͻ���ģʽ(client)"
Write-Host "2. �����ģʽ(server)"
$mode = Read-Host "����������"

if ($mode -ne "1" -and $mode -ne "2") {
    Write-Host "��Ч�����룬������ 1 �� 2��"
    exit 1
}

if ($mode -eq "1") {
    $mode = "client"
    $server_addr = Read-Host "�������Ʒ�������ַ"
} else {
    $mode = "server"
}


# ʹ�� $PWD.Path ȷ��·������ȷ��
$execPath = Join-Path -Path $PWD.Path -ChildPath $EXEC_FILE


if ($mode -eq "server") {
    Start-Process -FilePath "$execPath" -ArgumentList $mode -NoNewWindow -Wait
} else {
    $arguments = @($mode, $server_addr)
    Start-Process -FilePath "$execPath" -ArgumentList $arguments -NoNewWindow -Wait

}

