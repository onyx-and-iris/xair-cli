package xair

import "math"

func linGet(min float64, max float64, value float64) float64 {
	return min + (max-min)*value
}

func linSet(min float64, max float64, value float64) float64 {
	return (value - min) / (max - min)
}

func logGet(min float64, max float64, value float64) float64 {
	return min * math.Exp(math.Log(max/min)*value)
}

func logSet(min float64, max float64, value float64) float64 {
	return math.Log(value/min) / math.Log(max/min)
}

func mustDbInto(db float64) float64 {
	switch {
	case db >= 10:
		return 1
	case db >= -10:
		return float64((db + 30) / 40)
	case db >= -30:
		return float64((db + 50) / 80)
	case db >= -60:
		return float64((db + 70) / 160)
	case db >= -90:
		return float64((db + 90) / 480)
	default:
		return 0
	}
}

func mustDbFrom(level float64) float64 {
	switch {
	case level >= 1:
		return 10
	case level >= 0.5:
		return toFixed(float64(level*40)-30, 1)
	case level >= 0.25:
		return toFixed(float64(level*80)-50, 1)
	case level >= 0.0625:
		return toFixed(float64(level*160)-70, 1)
	case level >= 0:
		return toFixed(float64(level*480)-90, 1)
	default:
		return -90
	}
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(math.Round(num*output)) / output
}
