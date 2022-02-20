#!/bin/bash

name=tools.wechat.interface

if [ ! -n "$1" ];then
    echo '完整的发布命令格式：sh release.sh <版本号>'
    exit
fi

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ../../cmd/interface/interface.go

image=registry.cn-hangzhou.aliyuncs.com/consultech-saas/$name:$1

docker build -t $image .

rm -f ../../cmd/interface/interface.go

docker push $image

docker rmi $image