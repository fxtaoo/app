#!/usr/bin/env bash
# 更新线上文件，CDN 刷新

if [ -z "$1" ]; then
    echo "请输入要更新文件的绝对路径"
    exit 1
fi

if [[ ! -x main ]] ;then
    chmod u+x main
fi

./main \
-storage "aliyun" \
-bucket "" \
-storagePath "" \
-cdn "qiniu" \
-url ""  \
-loadPath "$1"