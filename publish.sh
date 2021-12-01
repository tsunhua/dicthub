#!/bin/bash

# 停止遠程服務
ssh tencent 'bash -s' < run.sh stop

# 發送應用到服務器
## 1) build & package
sh build.sh
zip -r static.zip static/

## 2) send
scp config_prod.yaml tencent:~/apps/boat/config.yaml
scp static.zip run.sh backup.sh boat tencent:~/apps/boat/

## 3) clean
rm -f static.zip boat

# 啟動遠程服務
ssh tencent 'bash -s' < run.sh start
