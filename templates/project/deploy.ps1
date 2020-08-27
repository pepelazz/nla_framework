# powershell.exe -executionpolicy bypass -file .\deploy.ps1
$ErrorActionPreference = "Stop"


function git_push {
    git add .
    git commit -m "m"
    git push origin master
}

# обновление из git
echo "full project git pull..."
git pull

# сборка бинарника
cd src
Remove-Item 'app'
$env:GOOS = "linux"
$env:GOARCH = "amd64"
echo "start build"
go build -o app 2>&1 # redirect error stream (2) to success stream (1)

# копирование бинарника на сервер
echo "transfer file to server..."
pscp -r app  {{.Config.WebServer.Username}}@{{.Config.WebServer.Ip}}:/{{.Config.WebServer.Path}}/src

cd ./webClient
echo "start quasar build..."
npx quasar build

# коммит в git
cd ../..
git_push
