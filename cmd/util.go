/*
LICENSE: https://github.com/onyx-and-iris/xair-cli/blob/main/LICENSE
*/
package cmd

import (
	"strconv"
)

func mustConvToFloat64(floatStr string) float64 {
	level, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		panic(err)
	}
	return level
}

func mustConvToInt(intStr string) int {
	val, err := strconv.Atoi(intStr)
	if err != nil {
		panic(err)
	}
	return val
}
