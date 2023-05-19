// tool合集
package main

import (
	"app/webTool/dfb"
	"app/webTool/lib"
	"bytes"
	"flag"
	"os"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	port := flag.String("port", "20231", "监听端口")
	flag.Parse()

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatFloat": lib.FormatFloat,
		"formatTime": func(t time.Time) string {
			return t.Format("2006/01/02")
		},
	})
	r.LoadHTMLGlob("templates/*.html")
	rG := r.Group("/tool")
	rG.GET("/", func(ctx *gin.Context) {
		file, _ := os.ReadFile("README.md")
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
			),
		)
		var buf bytes.Buffer
		if err := md.Convert(file, &buf); err != nil {
			panic(err)
		}

		ctx.Data(200, "text/html; charset=utf-8", buf.Bytes())
	})
	rG.GET("/dfb", dfb.Get)
	rG.POST("/dfb", dfb.Post)

	r.Run(":" + *port)
}
