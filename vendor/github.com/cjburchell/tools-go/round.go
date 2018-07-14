package tools

import (
	"fmt"
	"strconv"
)

func Round(value float64, digits int) float64 {
	format := fmt.Sprintf("%%.%df", digits)
	formatted, _ := strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return formatted
}
