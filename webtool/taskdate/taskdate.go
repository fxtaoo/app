package taskdate

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type V struct {
	StarDate   string
	EndDate    string
	Chapter    int
	Chapterday int
	Message    string
}

func Get(ctx *gin.Context) {
	// 默认吃饭时间为明天 07:50
	V := V{
		time.Now().Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
		0,
		0,
		"",
	}
	ctx.HTML(200, "taskdate.html", V)
}

func Post(ctx *gin.Context) {
	chapter, _ := strconv.Atoi(ctx.PostForm("chapter"))
	chapterday, _ := strconv.Atoi(ctx.PostForm("chapterday"))
	V := V{
		ctx.PostForm("startdate"),
		ctx.PostForm("enddate"),
		chapter,
		chapterday,
		"",
	}
	today := time.Now().Format("2006-01-02")
	startdate, _ := time.Parse("2006-01-02", V.StarDate)
	enddate, _ := time.Parse("2006-01-02", V.EndDate)

	if today != V.EndDate {
		// 输入 开始日期、结束日期、章节数，计算每天章节数
		manyday := int(enddate.Sub(startdate).Hours() / 24)
		tmpManyday := manyday
		for i := 1; i <= tmpManyday; i++ {
			if startdate.AddDate(0, 0, i).Weekday() == time.Sunday {
				manyday -= 1
			}
		}
		V.Chapterday = V.Chapter / manyday
		remainder := V.Chapter % manyday

		if remainder != 0 {
			V.Message = fmt.Sprintf("最后一天 %d", V.Chapterday+remainder)
		}
	} else {
		// 输入 开始日期、章节数、每天章节数，计算计算任务最后一天
		manyday := int(math.Ceil(float64(chapter) / float64(chapterday)))

		// 从明天开始，周末不算
		// 统计有多少个周末
		tmpManyday := manyday
		for i := 1; i <= tmpManyday; i++ {
			if startdate.AddDate(0, 0, i).Weekday() == time.Sunday {
				manyday += 1
			}
		}

		V.Message = fmt.Sprintf("%s 为计划任务最后一天", startdate.AddDate(0, 0, manyday).Format("2006-01-02"))
	}

	ctx.HTML(200, "taskdate.html", V)
}
