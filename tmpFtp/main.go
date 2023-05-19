// 临时 FTP
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
	gopsutilNet "github.com/shirou/gopsutil/net"
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

	address := "访问地址："
	if *ip != "" {
		address += fmt.Sprintf("\nhttp://%s:%s", *ip, strconv.Itoa(*port))
	} else {
		host, _ := host.Info()
		ip, _ := gopsutilNet.Interfaces()

		// 默认为 linux
		ipIndex := 0
		switch host.OS {
		case "windows", "darwin":
			ipIndex = 1
		}

		for _, e := range ip {
			if len(e.Addrs) > 1 {
				address += fmt.Sprintf("\nhttp://%s:%d %s",
					strings.Split(e.Addrs[ipIndex].Addr, "/")[0],
					*port,
					e.Name)
			}
		}
	}

	fmt.Printf("\n%s\n\n", address)

	r.Run(fmt.Sprintf("%s:%d", *ip, *port))
}
