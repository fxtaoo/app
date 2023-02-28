package voo

import (
	"app/webTool/conf"
	"app/webTool/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fxtaoo/golib/mail"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type V struct {
	PE            float64
	ROE           float64
	Value         float64
	Multiple      float64
	Money         float64
	MultipleMoney float64
	URL           string
	Date          time.Time
	Chan          chan struct{} `json:"-"`
	First         bool
}

func (v *V) GetWebData() error {
	// 避免数据冲突
	v.Chan <- struct{}{}
	defer func() { <-v.Chan }()

	VOO := struct {
		EquityCharacteristic struct {
			Fund struct {
				PE      string `json:"priceEarningsRatio"`
				ROE     string `json:"returnOnEquity"`
				ROEDate string `json:"returnOnEquityDate"`
			} `json:"fund"`
		} `json:"equityCharacteristic"`
	}{}

	resp, err := http.Get("https://investor.vanguard.com/investment-products/etfs/profile/api/VOO/characteristic")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&VOO)
	if err != nil {
		return err
	}

	v.PE, err = strconv.ParseFloat(VOO.EquityCharacteristic.Fund.PE[:len(VOO.EquityCharacteristic.Fund.PE)-1], 64)
	if err != nil {
		return err
	}
	v.ROE, err = strconv.ParseFloat(VOO.EquityCharacteristic.Fund.ROE, 64)
	if err != nil {
		return err
	}

	v.Date, err = time.Parse("2006-01-02", strings.Split(VOO.EquityCharacteristic.Fund.ROEDate, "T")[0])
	if err != nil {
		return err
	}

	v.Value = v.PE / v.ROE
	v.Multiple = 2 - v.Value
	if v.Multiple < 0 {
		v.Multiple = 0
	}
	v.MultipleMoney = v.Money * v.Multiple
	return nil
}

func (v *V) Get(ctx *gin.Context) {
	if ctx.Query("money") == "" {
		ctx.String(200, "没有设置金额，请 URL 参数方式 ?money=xxx 添加，限整数。")
	} else {
		if v.Date.IsZero() || v.URL != ctx.Request.URL.String() || v.First {
			var err error
			v.Money, err = strconv.ParseFloat(ctx.Query("money"), 64)
			if err != nil {
				ctx.String(200, err.Error())
			} else {
				err := v.GetWebData()
				if err != nil {
					ctx.String(200, err.Error())
				} else {
					v.URL = ctx.Request.URL.String()
					v.First = false
					ctx.HTML(200, "voo.html", v)
				}
			}
		} else {
			ctx.HTML(200, "voo.html", v)
		}

	}
}

func (v *V) Timing(conf *conf.Conf) {
	// 定时运行
	cron := cron.New()
	cron.AddFunc(conf.Alert.CronVOO, func() {
		v.First = true
		oldDate := v.Date
		v.GetWebData()
		if v.Date != oldDate {
			mail := mail.Mail{
				To:      conf.Alert.Mails,
				Subject: "VOO 数据以更新",
				Body: fmt.Sprintf("<html><head><style type=\"text/css\">table{border-collapse:collapse;border:2px solid rgb(200,200,200);letter-spacing:1px;font-size:0.8rem}td,th{border:1px solid rgb(190,190,190);padding:10px 20px}th{background-color:rgb(235,235,235)}td{text-align:center}caption{padding:10px}.small-font{font-size:0.8rem}</style></head><body><span class=\"small-font\">数据更新时间：%s</span><table><tr><td>名称</td><td>PE</td><td>ROE</td><td>价值</td><td>倍数</td><td>金额</td></tr><tr><td>VOO</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr></table></body></html>",
					v.Date.Format("2006-01-02"),
					lib.FormatFloat(v.PE),
					lib.FormatFloat(v.ROE),
					lib.FormatFloat(v.PE/v.ROE),
					lib.FormatFloat(v.Multiple),
					lib.FormatFloat(v.Money),
				),
			}
			mail.SendAlone(&conf.Smtp)
		}
	})
	cron.Start()
}
