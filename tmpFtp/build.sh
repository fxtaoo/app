#!/usr/bin/env bash
# 打包镜像

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o main .

docker build -t fxtaoo/tmpftp:latest .

docker push fxtaoo/tmpftp:latest
