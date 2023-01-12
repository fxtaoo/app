package dfb

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type V struct {
	DateTime string
	Cooktime int
	Message  string
}

func Get(ctx *gin.Context) {
	// 默认吃饭时间为明天 07:50
	V := V{
		time.Now().AddDate(0, 0, 1).Format("2006-01-02") + "T07:50",
		0,
		"",
	}

	if ctx.Query("cooktime") != "" {
		V.Cooktime, _ = strconv.Atoi(ctx.Query("cooktime"))
	}

	if ctx.Query("mealtime") != "" {
		V.DateTime = fmt.Sprintf("%sT%s",
			time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			ctx.Query("mealtime"),
		)
	}

	ctx.HTML(200, "dfb.html", V)
}

func Post(ctx *gin.Context) {
	V := V{
		ctx.PostForm("mealtime"),
		0,
		""}
	V.Cooktime, _ = strconv.Atoi(ctx.PostForm("cooktime"))

	mealtime, _ := time.Parse("2006-01-02T15:04", V.DateTime)
	// cst 转 utc
	mealtime = mealtime.Add(-time.Hour * 8)

	sinceMinute := int(mealtime.Sub(time.Now().UTC()).Minutes())

	if sinceMinute < V.Cooktime {
		V.Message = "煮饭时间不够！"
	} else {
		sinceMinute -= V.Cooktime
		V.Message = fmt.Sprintf("定时 %d 小时 %d 分钟后饭煮好", sinceMinute/60, sinceMinute%60)
	}

	ctx.HTML(200, "dfb.html", V)
}
