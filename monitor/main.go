// 一些监控
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
type CPUNumDisk struct {
	NotifyMailList                      []string // 邮件通知列表
	NotifyIntervalTime                  float64  // 通知间隔时间
	RepeatTime                          string   // 执行周期
	MaxiMumCPU, MaxiMumNUM, MaxiMumDisk float64  // 最大比例
	LogPath                             string   // 日志绝对路径
}

type RestartStopContainer struct {
	LogPath        string   // 日志绝对路径
	RepeatTime     string   // 执行周期
	NotifyMailList []string // 邮件通知列表
	ONOFF          bool     // 是否开启
}

type Config struct {
	CPUNumDisk           CPUNumDisk
	Smtp                 gomail.Smtp
	RestartStopContainer RestartStopContainer
}

func main() {
	configFile := flag.String("conf", "conf.toml", "配置文件名（当前目录）或绝对路径")
	flag.Parse()

	// 读配置
	var conf Config
	gofile.TomlFileRead(*configFile, &conf)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go checkCPUNumDisk(&conf)

	if conf.RestartStopContainer.ONOFF {
		go checkRestartStopContainer(&conf)
	}

	wg.Wait()
}

// 发送邮件写日志
func sendEmailWriteLog(smtp *gomail.Smtp, subject, content string, mailList []string, logPath *string) {
	// 日志缺省路径
	if *logPath == "" {
		*logPath = filepath.Join(filepath.Dir(os.Args[0]), "monitor.log")
	}

	// 服务器信息
	hostname, _ := host.Info()
	ip, _ := net.Interfaces()
	outPut := fmt.Sprintf("%s<br>服务器 %s IP %s<br>%s", time.Now().Format("2006-01-02 15:04:05"), hostname.Hostname, ip[1].Addrs[0].Addr, content)

	mail := gomail.Mail{Subject: subject, Body: outPut}

	gomail.SendEmailMP(smtp, &mail, mailList)

	// 写日志
	gofile.AppendFile(*logPath, fmt.Sprintf("%s\n%s", time.Now().Format("2006-01-02 15:04:05"), strings.ReplaceAll(content, "<br>", "\n")))
}

// 检查 CPU、磁盘、内存,超出设置比例发送邮件
func checkCPUNumDisk(conf *Config) {
	// 顺序为 cpu 内存 磁盘，有顺序对应关系
	warnList := []monitor.Warn{{Time: time.Now()}, {Time: time.Now()}, {Time: time.Now()}}
	checkList := []func(float64) (*monitor.Warn, error){monitor.CpuUsage, monitor.NumUsage, monitor.DiskUsage}
	maxNumList := []float64{conf.CPUNumDisk.MaxiMumCPU, conf.CPUNumDisk.MaxiMumNUM, conf.CPUNumDisk.MaxiMumDisk}

	c := cron.New()
	c.AddFunc(conf.CPUNumDisk.RepeatTime, func() {
		checkSendMail(conf, warnList, checkList, maxNumList)
	})
	c.Start()
}

func checkSendMail(conf *Config, warnList []monitor.Warn, checkList []func(float64) (*monitor.Warn, error), maxNumList []float64) {

	var warnContent string

	for i := range warnList {
		update, err := warnList[i].Check(checkList[i], maxNumList[i], conf.CPUNumDisk.NotifyIntervalTime)
		if err != nil {
			gofile.AppendFile(conf.CPUNumDisk.LogPath, err.Error())
			continue
		}
		if update {
			warnContent += warnList[i].Content + "<br>"
		}
	}

	if warnContent != "" {
		sendEmailWriteLog(&conf.Smtp, "磁盘 | 内存 | CPU 告警", warnContent, conf.CPUNumDisk.NotifyMailList, &conf.CPUNumDisk.LogPath)
	}
}

// 检查并重启停止容器
func checkRestartStopContainer(conf *Config) {
	c := cron.New()
	c.AddFunc(conf.RestartStopContainer.RepeatTime, func() {
		fmt.Println("日志路径=", conf.RestartStopContainer.LogPath)
		result, err := monitor.RestartStopContainer()
		if err != nil {
			gofile.AppendFile(conf.RestartStopContainer.LogPath, err.Error())
			return
		}
		if result != "" {
			sendEmailWriteLog(&conf.Smtp, "容器重启告警", strings.ReplaceAll(result, "\n", "<br>"), conf.RestartStopContainer.NotifyMailList, &conf.RestartStopContainer.LogPath)
		}
	})
	c.Start()
}
