// AD 域用户密码到期邮件提醒

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fxtaoo/golib/goemail"
	"github.com/fxtaoo/golib/gofile"
	"github.com/go-cmd/cmd"
	"github.com/robfig/cron/v3"
)

type User struct {
	email             string // 邮件地址
	passwdListSetDate string // 最后修改密码时间
	name              string // 姓名
}

type ADConfig struct {
	MxPasswordAge           int
	MailDomain, MailContent string
	AdminMail               []string
}

type Config struct {
	AD ADConfig
}

func sendAdmailCount(userList *[]string, content *string, sort string) {
	*content += sort + "\n"
	for _, e := range *userList {
		*content += e + "\n"
	}
}

// 读数据发送邮件
func readDataSendMail(logFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(mw)
	log.Println("\n" + time.Now().Format("2006-01-02 15:04:05"))

	// 配置
	var conf Config
	gofile.TomlFileRead(ConfigFileName, &conf)

	// 从域控生成数据
	cmd := cmd.NewCmd("powershell.exe", filepath.Join(filepath.Dir(os.Args[0]), PowerShellFileName))
	status := <-cmd.Start()

	for _, line := range status.Stdout {
		log.Println(line)
	}

	// 从域控生成数据生成数据需要时间
	time.Sleep(time.Minute * 3)

	// 用户数据
	var userList []User
	csvdata := gofile.CSVFileRead(UserDataFileName)
	for _, row := range csvdata {
		userList = append(userList, User{row[0], strings.Split(row[1], " ")[0], row[2]})
	}

	// 排除用户
	var excludeUser []string
	csvdata = gofile.CSVFileRead(ExcludeUserFileName)
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
			// 密码以过期
			expiredUser = append(expiredUser, user.email+" "+user.name+" 以过期 "+strconv.Itoa(dateInterval)[1:]+" 天 <br>")
		} else {
			if error := goemail.SendEmail(ConfigFileName, user.email+conf.AD.MailDomain, "人事 VPN 密码到期提醒", "<strong>人事 VPN 密码还有 "+strconv.Itoa(dateInterval)+" 天到期！请尽快按提示重置密码！</strong>"+conf.AD.MailContent, filepath.Join(filepath.Dir(os.Args[0]), "nopush-example.png")); error != nil {
				falseSendMail = append(falseSendMail, user.email+" "+user.name+" 还有 "+strconv.Itoa(dateInterval)+" 天到期 <br>")
				log.Println(user.email + conf.AD.MailDomain + " 通知邮件发送失败！")
			} else {
				trueSendMail = append(trueSendMail, user.email+" "+user.name+" 还有 "+strconv.Itoa(dateInterval)+" 天到期 <br>")
				log.Println(user.email + conf.AD.MailDomain + " 通知邮件发送成功！")
			}
		}
		time.Sleep(3 * time.Second)
	}

	sendAdmailCount(&expiredUser, &admailContent, "<br><strong>过期用户名：</strong><br>")
	sendAdmailCount(&falseSendMail, &admailContent, "<br><strong>提醒邮件发送失败用户名：</strong><br>")
	sendAdmailCount(&trueSendMail, &admailContent, "<br><strong>提醒邮件发送成功用户名:</strong><br>")

	// 域控邮件提醒统计
	for _, e := range conf.AD.AdminMail {
		goemail.SendEmail(ConfigFileName, e, "今日域控邮件提醒统计", admailContent)
		time.Sleep(3 * time.Second)
		log.Println(e + " 今日域控邮件提醒统计邮件发送成功！")
	}
}

const (
	ADLogFileName       = "adPasswdResetNotice.log"
	ConfigFileName      = "conf.toml"
	UserDataFileName    = "userdata.csv"
	ExcludeUserFileName = "excludeuser.csv"
	PowerShellFileName  = "finduser.ps1"
)

func main() {
	// 日志
	logFile, err := os.OpenFile(filepath.Join(filepath.Dir(os.Args[0]), ADLogFileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	c := cron.New()
	c.AddFunc("0 9 * * *", func() { readDataSendMail(logFile) })
	c.Start()

	wg.Wait()
}
