GOOS=linux GOARCH=amd64 go build -o sublink_amd64 main.go
GOOS=linux GOARCH=arm64 go build -o sublink_arm64 main.go
GOOS=windows GOARCH=amd64 go build -o sublink_amd64.exe main.go