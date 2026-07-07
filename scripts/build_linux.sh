rm -rf dst-admin-go
GOOS=linux GOARCH=amd64 go build -o dst-admin-go cmd/server/main.go