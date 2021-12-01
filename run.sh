#!/bin/bash
# 參數說明：
# - start 啟動程序
# - stop 停止程序
# - 重啟程序

cd ~/apps/boat || exit

COMMAND_GO_SERVER=./boat
go_server_pid=$(pidof ${COMMAND_GO_SERVER})

stop() {
  kill -9 "$go_server_pid"
  date +"%Y/%m/%d %H:%M:%S-Stop completed"
}

start() {
  # 清理和解壓資源
  rm -rf static/
  unzip static.zip
  nohup $COMMAND_GO_SERVER >./log/boat.log 2>&1 &
  date +"%Y/%m/%d %H:%M:%S-Start completed"
}

if [ "$1" == 'stop' ]; then
  stop
elif [ "$1" == 'start' ]; then
  start
else
  if [ "$go_server_pid" ]; then
    stop
    start
  else
    start
  fi
fi
