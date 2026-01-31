package cmd

import (
	"strconv"
)

// mustConvToFloat64 converts a string to float64, panicking on error.
func mustConvToFloat64(floatStr string) float64 {
	level, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		panic(err)
	}
	return level
}

// mustConvToInt converts a string to int, panicking on error.
func mustConvToInt(intStr string) int {
	val, err := strconv.Atoi(intStr)
	if err != nil {
		panic(err)
	}
	return val
}
