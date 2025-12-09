go build -o "../build/updater_console.exe"

#Без консоли
go build -ldflags="-H=windowsgui" -o "../build/updater.exe"
