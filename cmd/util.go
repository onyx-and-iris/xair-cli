/*
LICENSE: https://github.com/onyx-and-iris/xair-cli/blob/main/LICENSE
*/
package cmd

import (
	"strconv"
)

func mustConv(levelStr string) float64 {
	level, err := strconv.ParseFloat(levelStr, 64)
	if err != nil {
		panic(err)
	}
	return level
}
