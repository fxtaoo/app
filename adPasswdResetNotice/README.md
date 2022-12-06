# adPasswdResetNotice
AD 域用户密码到期邮件提醒  

## 配置
配置文件名为 conf.json,参考 conf-example.json

## 功能
1. 域控密码到期发送邮件提醒
2. 统计过期用户，发送通知失败、成功用户数据到指定邮箱

## 实现
运行在域控机器  
1. 执行 finduser.ps1 将筛选用户数据生成 userdata.csv
2. main.go 从 userdata.csv 获取分析数据，发送通知
3. 发送邮件失败或成功记录在 adPasswdResetNotice.log 中

finduser.ps1  
查找并生成用户信息列表数据，该列表用户符合以下条件：  
1. 用户开启  
2. 用户不是密码永不过期
3. 用户更改过密码
4. 距离密码过期天数小于 14 天






