// 检查 Mfs 有没有出现错误
package main

import (
	"app/mfsalarm/conf"
	"app/mfsalarm/disk"
	"fmt"
	"log"
	"time"

	"github.com/fxtaoo/golib/mail"
	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

func main() {
	// 初始化配置
	conf := conf.Conf{}
	err := conf.Init()
	if err != nil {
		log.Fatal(err)
	}

	errorDisks := make(map[string]*disk.Disk)

	// 定时运行
	cron := cron.New()
	// CheckIntervalTime 分钟运行
	cron.AddFunc(conf.DiskAlarm.CheckIntervalTime, func() {
		c := colly.NewCollector()
		c.OnHTML("table.acid_tab tr.C1,table.acid_tab tr.C2", func(e *colly.HTMLElement) {
			disk := &disk.Disk{}
			err := e.Unmarshal(disk)
			if err != nil {
				log.Fatal(err)
			}
			disk.Format()
			disk.Check(errorDisks, conf.DiskAlarm.AlarmIntervalTime)
		})
		c.Visit(conf.DiskAlarm.Url)

		disksAlertMail(errorDisks, &conf)
	})

	cron.Start()

	select {}
}

func disksAlertMail(errorDisks map[string]*disk.Disk, conf *conf.Conf) {
	content := ""
	timeNow := time.Now()
	for key, errDisk := range errorDisks {
		if errDisk.AlertMail {
			// 发送过邮件，且大于间隔时间，删除
			if timeNow.Sub(errDisk.RecordTime).Minutes() > float64(conf.DiskAlarm.AlarmIntervalTime) {
				delete(errorDisks, key)
			}
		} else {
			// 没有发送过邮件
			content += fmt.Sprintf("%s %s %s %s <br>",
				errDisk.ID,
				errDisk.IpPath,
				errDisk.LastErrorTime,
				errDisk.Status,
			)
			errDisk.AlertMail = true
		}
	}

	if content != "" {
		// 邮件
		mail := mail.Mail{
			To:      conf.DiskAlarm.Mails,
			Subject: "MFS 告警",
			Body:    content,
		}

		mail.SendAlone(&conf.Smtp)
	}
}
