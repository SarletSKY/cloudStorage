#!/bin/bash

# 检查service进程
check_process(){
    # ps 命令用于显示当前进程 (process) 的状态。
    sleep 1
    res=`ps aux | grep -v grep | grep "service/bin" | grep $1`
    if [[ $res != '' ]]; then
        echo -e "\033[32m 已启动 \033[0m" "$1"
        return 1
    else
        echo -e "\033[32m 启动失败 \033[0m" "$1"
        return 0
    fi
}

# 编译service可执行文件
build_service(){
    go build -o service/bin/$1 service/$1/main.go
    resbin=`ls service/bin/ | grep $1`
    echo -e "\033[32m 编译完成： \033[0m service/bin/$resbin"
}

# 启动service
run_service(){
    # nohup 英文全称 no hang up（不挂起），用于在系统后台不挂断地运行命令，退出终端不会影响程序的运行。
    # &：让命令在后台执行，终端退出后命令仍旧执行。
    nohup ./service/bin/$1 --registry=consul >> $logpath/$1.log 2>&1 &
    sleep 1
    check_process $1
}


# 创建运行日志目录
logpath=/home/zwx/go/src/filestore-server-study/log

mkdir -p $logpath

# 切换到工程根目录
cd $GOPATH/src/filestore-server-study

# 微服务可以用supervisor进程管理工具;
# 或者也可以通过docker/k8s进行部署

services="
upload
download
transfer
account
apigw
"

# 执行编译service
for service in $services
do
    build_service $service
done