package conf

import (
	"github.com/fxtaoo/golib/file"
	"github.com/fxtaoo/golib/mail"
)

type Conf struct {
	Smtp  mail.Smtp `json:"smtp"`
	Alert struct {
		Mails   []string `json:"mails"`
		CronVOO string   `json:"cronvoo"`
	} `json:"alert"`
}

func (c *Conf) Init() error {
	return file.JsonInitValue("conf.json", c)
}
