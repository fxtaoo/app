# updateOnlineFile
更新线上文件，CDN 刷新

## 配置
配置文件参考 conf-example.toml，缺省从当前目录 conf.toml 读取，否则使用参数 `-conf` 配置绝对路径指定 
参数使用参考 run.sh

## 功能
 - [x] 验证线上本地文件是否一致
 - [x] 阿里云 OSS 上传文件
 - [x] 七牛云 CDN 刷新
