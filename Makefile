BINARY_NAME=lib-smb2

build:
	GOARCH=amd64 GOOS=linux go build -buildmode=c-shared -o  bin/${BINARY_NAME} main.go;

dockerbuild:
	docker build -t go-compiler .
	docker run -v .:/app go-compiler /bin/sh -c "go build -buildmode=c-shared -o bin/lib-smb2 main.go"
