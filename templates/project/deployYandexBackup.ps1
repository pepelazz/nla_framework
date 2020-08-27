# powershell.exe -executionpolicy bypass -file .\deployYandexBackup.ps1
$ErrorActionPreference = "Stop"

# сборка бинарника для сервиса yandexBackup
cd src/yandexDiskBackup
Remove-Item 'yandexDiskBackup'
$env:GOOS = "linux"
$env:GOARCH = "amd64"
echo "start build yandexDiskBackup"
go build -o yandexDiskBackup 2>&1 # redirect error stream (2) to success stream (1)

# копирование бинарника на сервер
echo "transfer file to server..."
pscp -r yandexDiskBackup  [[.Config.WebServer.Username]]@[[.Config.WebServer.Ip]]:/[[.Config.WebServer.Path]]
cd ../..
