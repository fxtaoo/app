package etf

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type ETF struct {
	ETF   string
	PCT   float64
	Money int
}

type V struct {
	ETFS []ETF
	URL  string
}

var Value V

func Get(ctx *gin.Context) {
	Value.ETFS = Value.ETFS[:0]
	etfArray := ctx.QueryArray("etf")
	pctArray := ctx.QueryArray("pct")

	for i := range etfArray {
		pct, _ := strconv.ParseFloat(pctArray[i], 64)
		Value.ETFS = append(Value.ETFS, ETF{etfArray[i], pct, 0})
	}
	Value.URL = ctx.Request.URL.String()

	ctx.HTML(200, "etf.html", Value)
}

func POST(ctx *gin.Context) {
	money, _ := strconv.ParseFloat(ctx.PostForm("money"), 64)

	lo.ForEach(Value.ETFS, func(_ ETF, i int) {
		Value.ETFS[i].Money = int(Value.ETFS[i].PCT * money)
	})
	ctx.HTML(200, "etf.html", Value)
}
