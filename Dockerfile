# Use bullseye image instead of alpine.
# Alpine missing some libs that allow to compile go binaries
FROM golang:1.20.3

# Creates an app directory to hold your app’s source code
WORKDIR /app

