#!/usr/bin/env bash
# 创建新 app

github_repository_url="https://github.com/fxtaoo/app/tree/master"

while true ; do
  read -rp "app 名称：" app_name
  if [[ ! -d $app_name ]] ; then
    break;
  fi
  echo -e "重新输入 $app_name 该 app 以存在!\n"
done


read -rp "app 信息：" intro

content="//t $intro
package main

"

# 创建目录与文件
(mkdir ./$app_name && cd $app_name && \
go mod init app/${app_name}&& \
echo -e "$content" > ./main.go
echo -e "$content" > ./main_test.go
echo -e "# ${app_name} \n ${intro}" > ./README.md)

echo -e "| [${app_name}](${github_repository_url}/${app_name}) | ${intro} |">> ./README.md




