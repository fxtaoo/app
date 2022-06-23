// Actions 程序，生成 [fxtaoo/cmd](https://github.com/fxtaoo/cmd) 该仓库 README.md 文件
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fxtaoo/golib/gofile"
)

type InfoConf struct {
	TopPageInfo, CodeURL, GithubRaw, FxtaooRaw string
	DirInfo                                    [][]string
}

type Config struct {
	Info InfoConf
}

type script struct {
	path      string
	name      string
	info      string
	codeUrl   string
	githubRaw string
	fxtaooRaw string
}

type scriptDir struct {
	path    string
	info    string
	content string
}

// app 从文件读 info
func (s *script) readInfo() {
	f, err := os.Open(s.path)
	if err != nil {
		log.Fatalln(s.path + " 打开失败！")
	}
	defer f.Close()
	content := bufio.NewScanner(f)

	// 第一行 #!/usr/bin/env bash 不取
	content.Scan()

	// 第二行 信息
	content.Scan()
	s.info = content.Text()[2:]
}

// app 相关信息添加 markdown 标签
func (s script) appMarkdown() string {
	// return "[" + s.name + "]" + "(" + s.codeUrl + ")" + "　[githubRaw]" + "(" + s.githubRaw + ")" + " [fxtaooRaw]" + "(" + s.fxtaooRaw + ")" + "  \n" + s.info + "  "
	return fmt.Sprintf("| [%v](%v) | [%v](%v) |", s.name, s.codeUrl, s.info, s.githubRaw)

}

// 读整个文件夹下 app 信息，并聚合
func (d *scriptDir) readApp() {
	err := filepath.Walk(d.path, func(path string, file os.FileInfo, err error) error {
		// 遍历文件
		if file.Name() != d.path && !strings.HasPrefix(path, ".") {
			// 排除隐藏文件 . ..
			var tmpApp script
			tmpApp.name = file.Name()
			tmpApp.path = path
			tmpApp.codeUrl = conf.Info.CodeURL + d.path + "/" + tmpApp.name
			tmpApp.githubRaw = conf.Info.GithubRaw + d.path + "/" + tmpApp.name
			tmpApp.fxtaooRaw = conf.Info.FxtaooRaw + d.path + "/" + tmpApp.name

			tmpApp.readInfo()
			d.content += tmpApp.appMarkdown() + "\n"
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func initDir(d *scriptDir, path, info string, wg *sync.WaitGroup) {
	d.path = path
	d.info = info
	d.content = fmt.Sprintf("## %v\n| | |\n| :---- | :---- |\n", info)
	d.readApp()
	wg.Done()
}

var conf Config

func main() {

	gofile.TomlFileRead("conf.toml", &conf)

	appDirList := make([]scriptDir, len(conf.Info.DirInfo))

	var wg sync.WaitGroup
	for i := range appDirList {
		wg.Add(1)
		go initDir(&appDirList[i], conf.Info.DirInfo[i][0], conf.Info.DirInfo[i][1], &wg)
	}

	wg.Wait()

	content := conf.Info.TopPageInfo + "\n\n"

	for i := range appDirList {
		content += appDirList[i].content
	}

	// 写 README.md
	f, _ := os.OpenFile("README.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	_, err := f.WriteString(content)
	if err != nil {
		fmt.Println(err)
	}
}
