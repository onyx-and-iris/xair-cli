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
