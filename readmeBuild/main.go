// 读文件按标记规则生成文本
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fxtaoo/golib/gofile"
)

type Note struct {
	content string
	url     string
}

type Question struct {
	title string
	notes []Note
	path  string
	off   bool
}

type Repository struct {
	TopPageInfo string
	URL         string
}

type Config struct {
	Store Repository
}

func main() {

	// 读配置
	var conf Config
	gofile.TomlFileRead("conf.toml", &conf)

	var questions []Question
	var wg sync.WaitGroup
	err := filepath.Walk(".", func(path string, file os.FileInfo, err error) error {
		// 遍历文件夹
		if file.IsDir() && !strings.HasPrefix(path, ".") {
			// 筛选出文件夹且排除隐藏文件夹 .. 上一个文件夹
			var tmp Question
			tmp.path = path + "/main.go"
			questions = append(questions, tmp)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for i := range questions {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			readQFile(&questions[i])
		}(i)
	}
	wg.Wait()

	content := conf.Store.TopPageInfo

	for _, q := range questions {
		if q.off || len(q.notes) == 0 {
			continue
		}
		content += "### " + q.title + "\n"
		for _, n := range q.notes {
			content += "* [" + n.content + "](" + conf.Store.URL + n.url + ")\n"
		}
	}

	// 写 README.md
	f, _ := os.OpenFile("README.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println(err)
	}
}

// 从文件搜集信息
func readQFile(q *Question) {
	data, err := ioutil.ReadFile(q.path)
	if err != nil {
		panic(err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//off") {
			q.off = true
			break
		}
		if strings.HasPrefix(line, "//t") {
			q.title = strings.TrimPrefix(line, "//t ")
			continue
		}
		if strings.HasPrefix(line, "//n") {
			var n Note
			n.content = strings.TrimPrefix(line, "//n ")
			n.url = q.path + "#L" + strconv.Itoa(i+1)
			q.notes = append(q.notes, n)
		}
	}
}
