package main

import "strconv"

// mustConvToFloat64 converts a string to float64, panicking on error.
func mustConvToFloat64(floatStr string) float64 {
	level, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		panic(err)
	}
	return level
}
