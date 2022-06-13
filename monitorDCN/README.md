# monitorDCN
周期检查 cpu 内存 磁盘，超出使用率邮件通知  

## 配置

配置从同文件夹 conf.toml 读取  
配置参考示例 conf-example.toml  

## 功能

1. CPU 大于使用率持续 60 秒告警
2. 内存 大于使用率告警
3. 磁盘 大于使用率告警
4. 指定间隔时间执行
5. 指定同类型事件告警邮件间隔