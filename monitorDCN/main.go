// 周期检查 cpu 内存 磁盘，超出使用率邮件通知
package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/fxtaoo/golib/goemail"
	"github.com/fxtaoo/golib/gofile"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// 告警
type Warn struct {
	name    string // 什么告警
	content string // 告警内容
}

// 配置
type MonitorConfig struct {
	NotifyMailList                      []string
	NotifyIntervalTime                  float64
	RepeatTime                          string
	MaxiMumCPU, MaxiMumNUM, MaxiMumDisk float64
}

type Config struct {
	Monitor MonitorConfig
}

// CPU 高负载持续 60 秒 告警
// 每 10 秒取样一次
// 使用 num 警告比例
func cpuUsage(num float64) Warn {
	var warn Warn
	var tfList [6]bool

	for i := range tfList {
		v, _ := cpu.Percent(10*time.Millisecond, false)
		if v[0] > num {
			tfList[i] = true
		} else {
			tfList[i] = false
		}
		time.Sleep(10 * time.Second)
	}

	for _, e := range tfList {
		if !e {
			return warn
		}
	}
	warn = Warn{"cpu", "使用率超过 " + strconv.FormatFloat(num, 'g', 2, 64) + "% 持续一分钟"}
	return warn
}

// 内存告警
// 使用 num 警告比例
func numUsage(num float64) Warn {
	var warn Warn
	v, _ := mem.VirtualMemory()
	used := 100 - float64(v.Available)/float64(v.Total)*100
	if used > num {
		warn = Warn{"内存", "使用率超过 " + strconv.FormatFloat(num, 'g', 2, 64) + "%"}
		return warn
	} else {
		return warn
	}
}

// 磁盘告警
// 使用 num 警告比例
func diskUsage(num float64) Warn {
	var warn Warn
	isnum := 0
	partitions, _ := disk.Partitions(false)
	for _, e := range partitions {
		info, _ := disk.Usage(e.Mountpoint)
		used := float64(info.Used) / float64(info.Total) * 100
		if used > num {
			warn = Warn{e.Device, "使用率超过 " + strconv.FormatFloat(num, 'g', 2, 64) + "%"}
			isnum++
		}
	}
	if isnum > 0 {
		return warn
	} else {
		return warn
	}
}

func checkSendMail(conf Config, sortList []string, lastSendEmailTime map[string]time.Time) {

	funcList := []Warn{cpuUsage(conf.Monitor.MaxiMumCPU), numUsage(conf.Monitor.MaxiMumNUM), diskUsage(conf.Monitor.MaxiMumDisk)}
	monitorSort := make(map[string]Warn)

	for i, e := range sortList {
		monitorSort[e] = funcList[i]
	}

	// 服务器信息
	nowTime := time.Now()
	logs := nowTime.Format("2006-01-02 15:04:05")
	hostname, _ := host.Info()
	logs += "  服务器 " + hostname.Hostname

	ip, _ := net.Interfaces()
	logs += "  IP " + ip[1].Addrs[0].Addr + "\n"

	// 只要有一项满足即发邮件
	onoff := false
	for sort, warn := range monitorSort {
		if warn.content != "" && (lastSendEmailTime[sort].IsZero() || nowTime.Sub(lastSendEmailTime[sort]).Minutes() > conf.Monitor.NotifyIntervalTime) {
			logs += "\n" + warn.name + " " + warn.content
			lastSendEmailTime[sort] = nowTime
			onoff = true
		}
	}
	if onoff {
		for _, e := range conf.Monitor.NotifyMailList {
			goemail.SendEmail(ConfigFile, e, "磁盘 | 内存 | CPU 报警", logs)
		}
	}

}

const ConfigFile = "conf.toml"

func main() {

	var conf Config
	gofile.TomlFileRead(ConfigFile, &conf)

	// 最后一次邮件发送时间
	lastSendEmailTime := make(map[string]time.Time)
	sortList := []string{"cpu", "num", "disk"}

	for _, e := range sortList {
		lastSendEmailTime[e] = time.Time{}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()
	c.AddFunc(conf.Monitor.RepeatTime, func() { checkSendMail(conf, sortList, lastSendEmailTime) })
	c.Start()

	wg.Wait()
}
