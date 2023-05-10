# Specifies a parent image
FROM golang:1.20.3

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
# COPY . .
# RUN go mod download
# CMD GOARCH=amd64 GOOS=linux go build -buildmode=c-shared -o  ./bin/lib-smb2 main.go
