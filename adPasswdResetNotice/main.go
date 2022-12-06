// AD 域用户密码到期邮件提醒gomail

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fxtaoo/golib/file"
	"github.com/fxtaoo/golib/mail"
	"github.com/go-cmd/cmd"
	"github.com/robfig/cron/v3"
)

const (
	ADLogFileName       = "adPasswdResetNotice.log"
	ConfigFileName      = "conf.json"
	UserDataFileName    = "userdata.csv"
	ExcludeUserFileName = "excludeuser.csv"
	PowerShellFileName  = "finduser.ps1"
	MailContentFileName = "mailContent"
)

var dirPath string

type User struct {
	email             string // 邮件地址
	passwdListSetDate string // 最后修改密码时间
	name              string // 姓名
}

type ADConfig struct {
	MxPasswordAge int      `json:"mxPasswordAge"`
	DayStart      int      `json:"dayStart"`
	MailDomain    string   `json:"mailDomain"`
	MailContent   string   `json:"-"`
	AdminMail     []string `json:"adminMail"`
}

type Config struct {
	AD   ADConfig  `json:"ad"`
	Stmp mail.Smtp `json:"smtp"`
}

func init() {
	dirPath = filepath.Dir(os.Args[0])
}

func main() {
	// 日志
	logFile, err := os.OpenFile(filePath(ADLogFileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// 配置
	var conf Config
	file.JsonInitValue(filePath(ConfigFileName), &conf)

	email := mail.Mail{
		To:         make([]string, 1),
		Subject:    "AD 域用户密码到期邮件提醒",
		AttachPath: filePath("nopush-example.png"),
	}

	// mailContent
	tmpByte, _ := os.ReadFile(filePath(MailContentFileName))
	conf.AD.MailContent = string(tmpByte)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()
	c.AddFunc("0 9 * * *", func() { readDataSendMail(&conf, logFile, &email) })
	c.Start()

	wg.Wait()
}

func filePath(name string) string {
	return filepath.Join(dirPath, name)
}

func sendAdmailCount(userList *[]string, content *string, sort string) {
	*content += sort + "\n"
	for _, e := range *userList {
		*content += e + "\n"
	}
}

// 读数据发送邮件
func readDataSendMail(conf *Config, logFile *os.File, email *mail.Mail) {
	mw := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(mw)
	log.Println("\n" + time.Now().Format("2006-01-02 15:04:05"))

	// 从域控生成数据
	cmd := cmd.NewCmd("powershell.exe", filePath(PowerShellFileName))
	status := <-cmd.Start()

	for _, line := range status.Stdout {
		log.Println(line)
	}

	// 从域控生成数据生成数据需要时间
	time.Sleep(time.Minute * 3)

	// 用户数据
	var userList []User
	csvdata, _ := file.CSVRead(filePath(UserDataFileName))
	for _, row := range csvdata {
		userList = append(userList, User{row[0], strings.Split(row[1], " ")[0], row[2]})
	}

	// 排除用户
	var excludeUser []string
	csvdata, _ = file.CSVRead(filePath(ExcludeUserFileName))
	for _, row := range csvdata {
		excludeUser = append(excludeUser, row[0])
	}

	var expiredUser, trueSendMail, falseSendMail []string
	admailContent := ""

	for _, user := range userList[1:] {
		// 排除用户
		var userInExcludeSwitch bool
		for _, e := range excludeUser {
			if e == user.email {
				userInExcludeSwitch = true
				break
			}
		}

		if userInExcludeSwitch {
			continue
		}

		a := time.Now()
		b, _ := time.Parse("2006-1-2", strings.ReplaceAll(user.passwdListSetDate, "/", "-"))
		d := a.Sub(b)
		dateInterval := (conf.AD.MxPasswordAge - int(d.Hours()/24))

		if dateInterval < 0 {
			// 密码以彻底过期
			expiredUser = append(expiredUser, fmt.Sprintf("%s %s 以过期 %d 天 <br>", user.email, user.name, dateInterval*-1))
		} else if dateInterval <= conf.AD.DayStart {
			email.To[0] = user.email + conf.AD.MailDomain
			email.Body = fmt.Sprintf("<strong>人事 VPN 密码还有 %d 天到期！请尽快按提示重置密码<br>%s", dateInterval, conf.AD.MailContent)
			if error := email.SendAlone(&conf.Stmp); error != nil {
				falseSendMail = append(falseSendMail, user.email+" "+user.name+" 还有 "+strconv.Itoa(dateInterval)+" 天到期 <br>")
				log.Println(user.email + conf.AD.MailDomain + " 通知邮件发送失败！")
			} else {
				trueSendMail = append(trueSendMail, fmt.Sprintf("%s %s 还有 %s 天到期 <br>", user.email, user.name, strconv.Itoa(dateInterval)))
				log.Println(user.email + conf.AD.MailDomain + " 通知邮件发送成功！")
			}
		}
		time.Sleep(3 * time.Second)
	}

	sendAdmailCount(&falseSendMail, &admailContent, "<br><strong>提醒邮件发送失败用户名：</strong><br>")
	sendAdmailCount(&trueSendMail, &admailContent, "<br><strong>提醒邮件发送成功用户名:</strong><br>")
	sendAdmailCount(&expiredUser, &admailContent, "<br><strong>过期用户名：</strong><br>")

	// 域控邮件提醒统计
	email.Subject = "今日域控邮件提醒统计"
	email.Body = admailContent
	email.To = conf.AD.AdminMail
	if err := email.SendAlone(&conf.Stmp); err != nil {
		log.Println("域控邮件提醒统计发送失败！")
	} else {
		for i := range err {
			log.Println(email.To[i] + "域控邮件提醒统计发送失败！")
		}
	}
}
