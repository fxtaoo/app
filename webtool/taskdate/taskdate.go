package taskdate

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type V struct {
	Date       string
	Chapter    float64
	Chapterday float64
	Message    string
}

func Get(ctx *gin.Context) {
	// 默认吃饭时间为明天 07:50
	V := V{
		time.Now().Format("2006-01-02"),
		0,
		1,
		"",
	}
	ctx.HTML(200, "taskdate.html", V)
}

func Post(ctx *gin.Context) {
	chapter, _ := strconv.ParseFloat(ctx.PostForm("chapter"), 64)
	chapterday, _ := strconv.ParseFloat(ctx.PostForm("chapterday"), 64)
	V := V{
		ctx.PostForm("startdate"),
		chapter,
		chapterday,
		"",
	}
	startdate, _ := time.Parse("2006-01-02", V.Date)

	manyday := int(math.Ceil(chapter / chapterday))

	// 从明天开始，周末不算
	// 统计有多少个周末
	tmpManyday := manyday
	for i := 1; i <= tmpManyday; i++ {
		if startdate.AddDate(0, 0, i).Weekday() == time.Sunday {
			manyday += 1
		}
	}

	V.Message = startdate.AddDate(0, 0, manyday).Format("2006-01-02")

	ctx.HTML(200, "taskdate.html", V)
}
