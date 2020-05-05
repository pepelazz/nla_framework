# powershell.exe -executionpolicy bypass -file .\generateAndRun.ps1
$ErrorActionPreference = "Stop"

# генерация кода
cd projectTemplate
echo "start generate"
go run .

# рестарт проекта
cd ../src
echo "start project"
go run . -dev
