# mfsalarm
检查 Mfs 并邮件告警

## 配置

配置文件为执行文件同目录下 `conf.json` 文件
具体配置参考 `conf-example.json`

## 功能

检查 Mfs Disks，`status` 不为 `ok`，发送邮件告警。
注意 `alarmIntervalTime` 为告警间隔时间（分钟）：
例如配置为 `30`,是指同一错误，邮件告警后 30 分钟内，不会重复邮件告警。  
该功能目的是留有解决问题时间，避免邮箱塞满无意义的告警邮件。
