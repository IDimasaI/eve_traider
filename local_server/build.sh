
# Build the project
go build -o ./../build/local_server.exe

go build -ldflags="-H=windowsgui" -o "./../build/local_server_gui.exe"
