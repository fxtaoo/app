// 周期检查 cpu 内存 磁盘，超出使用率邮件通知
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fxtaoo/golib/gofile"
	"github.com/fxtaoo/golib/gomail"
	"github.com/fxtaoo/golib/monitor"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/net"
)

// 配置
type MonitorConfig struct {
	NotifyMailList                      []string
	NotifyIntervalTime                  float64
	RepeatTime                          string
	MaxiMumCPU, MaxiMumNUM, MaxiMumDisk float64
}

type Config struct {
	Monitor MonitorConfig
	Smtp    gomail.Smtp
}

func checkSendMail(conf *Config, warnList []monitor.Warn, checkList []func(float64) (*monitor.Warn, error), maxNumList []float64) {

	var warnContent string

	for i := range warnList {
		update, err := warnList[i].Check(checkList[i], maxNumList[i], conf.Monitor.NotifyIntervalTime)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if update {
			warnContent += warnList[i].Content + "<br>"
		}
	}

	// 发送邮件
	if warnContent != "" {
		// 服务器信息
		hostname, _ := host.Info()
		ip, _ := net.Interfaces()
		outPut := fmt.Sprintf("%s<br>服务器 %s IP %s<br>%s", time.Now().Format("2006-01-02 15:04:05"), hostname.Hostname, ip[1].Addrs[0].Addr, warnContent)

		mail := gomail.Mail{Subject: "磁盘 | 内存 | CPU 报警", Body: outPut}

		gomail.SendEmailMP(&conf.Smtp, &mail, conf.Monitor.NotifyMailList)
	}

}

const ConfigFile = "conf.toml"

func main() {
	var conf Config
	gofile.TomlFileRead(ConfigFile, &conf)

	// 顺序为 cpu 内存 磁盘，有对应关系
	warnList := []monitor.Warn{{Time: time.Now()}, {Time: time.Now()}, {Time: time.Now()}}
	checkList := []func(float64) (*monitor.Warn, error){monitor.CpuUsage, monitor.NumUsage, monitor.DiskUsage}
	maxNumList := []float64{conf.Monitor.MaxiMumCPU, conf.Monitor.MaxiMumNUM, conf.Monitor.MaxiMumDisk}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()
	c.AddFunc(conf.Monitor.RepeatTime, func() { checkSendMail(&conf, warnList, checkList, maxNumList) })
	c.Start()

	wg.Wait()
}
