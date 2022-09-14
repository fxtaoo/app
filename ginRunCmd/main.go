// gin 练习，访问页面，执行命令
package main

import (
	"encoding/json"
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
	ServicesTag     map[string]interface{}
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

// func (s *Servers) Save() {
// 	data, err := json.Marshal(s.ServicesTag)
// 	if err != nil {
// 		log.Fatalln(err)
// 		return
// 	}
// 	os.WriteFile(s.ServicesTagPath, data, 0666)
// }

// func (s *Servers) update(server *Servers) {
// 	s.ServicesTag[Server.Sort].(map[string]interface{})[server.Name] = server.tag
// }

func main() {
	s := Servers{ServicesTagPath: "services_tag.json"}
	s.read()
	r := gin.Default()
	r.LoadHTMLFiles("template/update.html")
	r.GET("/test/update", func(ctx *gin.Context) {
		ctx.HTML(200, "update.html", s)
	})
	r.POST("/test/update", func(ctx *gin.Context) {
		server := Server{
			Sort:   ctx.Query("sort"),
			Name:   ctx.PostForm("name"),
			TagNew: strings.TrimSpace(ctx.PostForm("tag")),
		}

		tmp := s.ServicesTag[server.Sort]
		server.TagOld = tmp.(map[string]interface{})[server.Name].(string)

		cmd := exec.Command("fab", "update", "--ij", server.Sort, "-n", server.Name, "-t", server.TagNew)
		// cmd := exec.Command("cat","t")
		out, _ := cmd.CombinedOutput()
		server.Result = strings.Split(string(out), "\n")
		s.UpdateServer = server
		ctx.HTML(200, "update.html", s)
		s.UpdateServer = Server{}
	})
	r.Run(":8173")
}
