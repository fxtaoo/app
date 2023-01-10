package disk

import (
	"fmt"
	"strings"
	"time"
)

type Disk struct {
	ID            string    `selector:":nth-child(1)"`
	IpPath        string    `selector:":nth-child(2)"`
	LastErrorTime string    `selector:":nth-child(4)"`
	Status        string    `selector:":nth-child(5)"`
	RecordTime    time.Time `selector:"-"`
	AlertMail     bool      `selector:"-"`
}

// 格式化信息
func (d *Disk) Format() {
	d.IpPath = strings.Split(d.IpPath, " ")[0]
	// 分割并去除第一个字符，去除多余 0
	d.IpPath = fmt.Sprintf("%s:%s",
		strings.Split(d.IpPath, ":")[0][1:],
		strings.Split(d.IpPath, ":")[1][1:],
	)

	d.LastErrorTime = strings.Join(strings.Split(d.LastErrorTime, " ")[1:], " ")
}

// 检查，添加新的 errDisk
func (d *Disk) Check(errorDisks map[string]*Disk, alarmIntervalTime int) {
	if d.Status != "ok" {
		if errorDisk, ok := errorDisks[d.ID]; ok {
			// 发送过邮件
			// 大于间隔时间，重置，再次告警
			if time.Since(errorDisk.RecordTime).Minutes() > float64(alarmIntervalTime) {
				errorDisks[d.ID].RecordTime = time.Now()
				errorDisks[d.ID].AlertMail = false
			}
		} else {
			d.RecordTime = time.Now()
			errorDisks[d.ID] = d
		}
	}
}
