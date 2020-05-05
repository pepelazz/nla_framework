# powershell.exe -executionpolicy bypass -file .\generate.ps1
$ErrorActionPreference = "Stop"

# генерация кода
cd projectTemplate
echo "start generate"
$StartTime = (Get-Date).Second
go run .
$EndTime = (Get-Date).Second
echo "time elapsed $($EndTime - $StartTime) sec"


