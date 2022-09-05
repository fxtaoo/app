// 更新线上文件，CDN 刷新
package main

import (
	"flag"
	"log"
	"time"

	"github.com/fxtaoo/golib/goaliyun"
	"github.com/fxtaoo/golib/gofile"
	"github.com/fxtaoo/golib/goqiniu"
)

type AliyunOSS struct {
	Endpoint, AccessKeyID, AccessKeySecret string
}

type Qiniu struct {
	AccessKey, SecretKey string
}

type Config struct {
	AliyunOSS AliyunOSS
	Qiniu     Qiniu
}

type FileInfo struct {
	storage     string // 存储种类
	cdn         string // CDN 种类
	bucket      string // oss bucket 名称
	url         string // 访问网址
	storagePath string // 文件存储路径
	loadPath    string // 本地文件路径
}

func main() {
	// 初始化文件信息
	file := FileInfo{}
	var confPath string

	flag.StringVar(&file.storage, "storage", "", "存储种类（可选 aliyun）")
	flag.StringVar(&file.cdn, "cdn", "", "CDN 种类（可选 qiniu）")
	flag.StringVar(&file.bucket, "bucket", "", "OSS bucket 名称")
	flag.StringVar(&file.url, "url", "", "访问网址")
	flag.StringVar(&file.storagePath, "storagePath", "", "文件存储路径")
	flag.StringVar(&file.loadPath, "loadPath", "", "本地文件路径")
	flag.StringVar(&confPath, "conf", "conf.toml", "配置文件路径")
	flag.Parse()

	// 读配置
	conf := Config{}
	gofile.TomlFileRead(confPath, &conf)

	var err error
	// 上传文件
	switch file.storage {
	case "aliyun":
		err = upFileAilyunOSS(&conf, &file)
	}
	if err != nil {
		log.Fatalln(err)
	}

	// 更新文件
	switch file.cdn {
	case "qiniu":
		err = urlsRefreshQiniu(&conf, &file)
	}
	if err != nil {
		log.Fatalln(err)
	}
	// 比较文件尝试 2 次，等待 3 分钟，考虑 cdn 刷新需要时间
	var compareFileMD5 bool
	num := 0
	for !compareFileMD5 && num < 2 {
		time.Sleep(2 * time.Minute)
		compareFileMD5, err = gofile.OnlineLocalMD5Same(file.url, file.loadPath)
		if err != nil {
			log.Fatalln(err)
		}
		num++
	}

	if compareFileMD5 {
		log.Println("文件更新成功")
	} else {
		log.Println("文件更新失败")
	}
}

// 阿里云 OSS 上传文件
func upFileAilyunOSS(conf *Config, f *FileInfo) error {
	// 初始化 OSS 实例
	oss := goaliyun.OSS{Endpoint: conf.AliyunOSS.Endpoint, AccessKeyId: conf.AliyunOSS.AccessKeyID, AccessKeySecret: conf.AliyunOSS.AccessKeySecret}
	client, err := goaliyun.OSSClient(&oss)
	if err != nil {
		return err
	}

	file := goaliyun.File{BucketName: f.bucket, BucketFilePath: f.storagePath, LoadFilePath: f.loadPath}
	err = goaliyun.UploadFile(client, &file)
	if err != nil {
		return err
	}

	return nil
}

// 七牛云 CDN 刷新
func urlsRefreshQiniu(conf *Config, f *FileInfo) error {
	key := goqiniu.Key{AccessKey: conf.Qiniu.AccessKey, SecretKey: conf.Qiniu.SecretKey}
	cdn := goqiniu.CdnManager(key)
	_, err := goqiniu.UrlsRefresh(cdn, []string{f.url})
	if err != nil {
		return err
	}
	return nil
}
