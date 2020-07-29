#!/bin/bash

# 检查service进程
check_process(){
    sleep 1
    res=`ps aux |grep -v grep | grep "service/bin" | grep $1`
    if [[ $res != '' ]]; then
        echo -e "\033[32m 已启动 \033[0m" "$1"
        return 1
    else
        echo -e "\033[32m 启动失败 \033[0m" "$1"
        return 0
    fi
}

run_service(){
    nohup ./service/bin/$1 --registry=consul >> $logpath/$1.log 2>&1 &
    sleep 1
    check_process $1
}

logpath=/home/zwx/go/src/filestore-server-study

mkdir -p $logpath

cd $GOPATH/src/filestore-server-study

services="
upload
download
transfer
account
apipw
"

#执行编译service
for service in $services
do
    build_service $service
done

# 执行启动service
for service in $services
do
    run_service $service
done

echo "微服务启动完毕"