BINARY_NAME=lib-smb2

build:
	GOARCH=amd64 GOOS=linux go build -buildmode=c-shared -o  bin/${BINARY_NAME} main.go
