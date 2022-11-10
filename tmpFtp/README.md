# tmpFtp

临时 ftp

命令参数

```bash
  -dir string
        指定文件夹路径 (default ".")
  -ip string
        IP
  -list
        文件列表 (default true)
  -passwd string
        密码，不配置无 HTTP 基本认证
  -port int
        端口 (default 9527)
  -user string
        用户名，不配置无 HTTP 基本认证
```

docker 方式运行

```bash
# 默认端口 9527
# 默认文件夹 /dir
docker run --rm -it -p 9527:9527  -v /dir:/dir fxtaoo/tmpftp

# 使用命名参数
docker run --rm -it -p 9527:9527  -v /dir:/dir fxtaoo/tmpftp -user=test -passwd=test -dir=/dir
# 注意！这种运行方式 默认端口 9527，默认文件夹为 . ，需要使用 -dir=/dir 指定文件夹
```
