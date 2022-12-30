#!/usr/bin/env bash
# 打包镜像

CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o main .

docker build -t fxtaoo/tool:latest .

docker push fxtaoo/tool:latest
