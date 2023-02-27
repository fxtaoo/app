package lib

import (
	"fmt"
	"strconv"
)

func FormatFloat(f float64) string {
	tmpF, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	return fmt.Sprintf("%g", tmpF)
}
