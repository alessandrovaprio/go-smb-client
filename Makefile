BINARY_NAME=lib-smb2

build:
	GOARCH=amd64 GOOS=linux go build -buildmode=c-shared -o  bin/${BINARY_NAME} main.go;
	#GOARCH=arm64 GOOS=linux go build main.go;

dockerbuild:
	docker build -t go-compiler .
	docker run -v ${PWD}:/app go-compiler /bin/sh -c "go build -buildmode=c-shared -o bin/lib-smb2 main.go"

testwithdocker:
	docker rm -f samba
	docker run -d --network host -v ${PWD}/tests/share:/shared --name samba pwntr/samba-alpine;
