// 发送邮件
package main

import (
	"flag"

	"github.com/fxtaoo/golib/file"
	"github.com/fxtaoo/golib/mail"
)

type Config struct {
	Smtp mail.Smtp
}

func main() {
	configFile := flag.String("conf", "conf.toml", "配置文件名（当前目录）或绝对路径")
	to := flag.String("to", "", "接收邮箱")
	subject := flag.String("sub", "", "邮件主题")
	body := flag.String("body", "", "邮件内容")
	attachPath := flag.String("attach", "", "附件路径")
	flag.Parse()

	if *to == "" {
		panic("接收邮箱不能为空")
	}

	var conf Config
	file.TomlInitValue(*configFile, &conf)

	var mail mail.Mail
	mail.To[0] = *to
	mail.Subject = *subject
	mail.Body = *body
	mail.AttachPath = *attachPath

	if err := mail.Send(&conf.Smtp); err != nil {
		panic(err)
	}
}
