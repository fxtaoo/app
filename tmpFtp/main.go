// 临时 FTP
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	dir := flag.String("dir", ".", "指定文件夹路径")
	list := flag.Bool("list", true, "文件列表")
	ip := flag.String("ip", "", "IP")
	port := flag.Int("port", 9527, "端口")
	user := flag.String("user", "", "用户名，不配置无 HTTP 基本认证")
	passwd := flag.String("passwd", "", "密码，不配置无 HTTP 基本认证")

	flag.Parse()

	// 文件夹是否可读
	_, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	// 验证
	if *user != "" && *passwd != "" {
		r.Use(gin.BasicAuth(gin.Accounts{
			*user: *passwd,
		}))
	}

	r.StaticFS("/", gin.Dir(*dir, *list))

	r.Run(fmt.Sprintf("%s:%d", *ip, *port))
}
