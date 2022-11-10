// gin 练习，访问页面，执行命令
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Sort   string
	Name   string
	TagNew string
	TagOld string
	Result []string
}

type Servers struct {
	ServicesTagPath string
	ServicesTag     map[string]map[string]string
	UpdateServer    Server
}

func (s *Servers) read() {

	f, err := os.ReadFile(s.ServicesTagPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(f, &s.ServicesTag)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Servers) save(tagPath string) {
	data, err := json.Marshal(s.ServicesTag)
	if err != nil {
		log.Fatalln(err)
		return
	}
	os.WriteFile(tagPath, data, 0666)
}

func main() {
	ip := flag.String("ip", "", "IP")
	port := flag.Int("port", 9158, "端口")
	user := flag.String("user", "", "用户名，不配置无 HTTP 基本认证")
	passwd := flag.String("passwd", "", "密码，不配置无 HTTP 基本认证")
	flag.Parse()

	r := gin.Default()
	if *user != "" && *passwd != "" {
		r.Use(gin.BasicAuth(gin.Accounts{
			*user: *passwd,
		}))
	}

	r.LoadHTMLFiles("template/update.html")

	r.GET("/test/update", func(ctx *gin.Context) {
		s := Servers{ServicesTagPath: "services_tag.json"}
		s.read()

		ctx.HTML(200, "update.html", s)
	})

	r.POST("/test/update", func(ctx *gin.Context) {
		s := Servers{ServicesTagPath: "services_tag.json"}
		s.read()

		server := Server{
			Sort:   ctx.Query("sort"),
			Name:   ctx.PostForm("name"),
			TagNew: strings.TrimSpace(ctx.PostForm("tag")),
		}

		server.TagOld = s.ServicesTag[server.Sort][server.Name]

		cmd := exec.Command("fab", "update", "--ij", server.Sort, "-n", server.Name, "-t", server.TagNew)
		// cmd := exec.Command("echo", "fab", "update", "--ij", server.Sort, "-n", server.Name, "-t", server.TagNew)
		out, err := cmd.CombinedOutput()
		if err == nil {
			s.save(s.ServicesTagPath + ".b")
			s.ServicesTag[server.Sort][server.Name] = server.TagNew
			s.save(s.ServicesTagPath)
		}

		server.Result = strings.Split(string(out), "\n")
		s.UpdateServer = server
		ctx.HTML(200, "update.html", s)
		s.UpdateServer = Server{}
	})

	r.Run(fmt.Sprintf("%s:%d", *ip, *port))
}
