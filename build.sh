#!/bin/bash -e
# 構建服務器上可用的應用安裝包
export GO111MODULE=on
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -o boat
