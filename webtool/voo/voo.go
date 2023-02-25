package voo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type V struct {
	PE       float64
	ROE      float64
	Value    float64
	Multiple float64
	Money    float64
	URL      string
	Time     time.Time
}

func (v *V) Init() error {
	VOO := struct {
		EquityCharacteristic struct {
			Fund struct {
				PE  string `json:"priceEarningsRatio"`
				ROE string `json:"returnOnEquity"`
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

	v.Value = v.PE / v.ROE
	v.Multiple = 2 - v.Value
	if v.Multiple < 0 {
		v.Multiple = 0
	}
	v.Money = v.Money * v.Multiple
	return nil
}

func (v *V) Get(ctx *gin.Context) {
	if v.Time.IsZero() || time.Since(v.Time).Hours() > 24 || v.URL != ctx.Request.URL.String() {
		if ctx.Query("money") == "" {
			ctx.String(200, "没有设置金额，请 URL 参数方式 ?money=xxx 添加，限整数。")
		} else {
			var err error
			v.Money, err = strconv.ParseFloat(ctx.Query("money"), 64)
			if err != nil {
				ctx.String(200, err.Error())
			} else {
				err := v.Init()
				if err != nil {
					ctx.String(200, err.Error())
				} else {
					v.Time = time.Now()
					v.URL = ctx.Request.URL.String()
					ctx.HTML(200, "voo.html", v)
				}
			}

		}
	} else {
		ctx.HTML(200, "voo.html", v)
	}

}
