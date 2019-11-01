package model

import (
	"fmt"
	"testing"
)

func TestNewStatisticSimple(t *testing.T) {
	tol1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190704",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	tcl1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190712",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10500",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tsl1 := []*FuturesTransaction{
		tol1,
		tcl1,
	}

	tradel1, err := NewFuturesTrade(tsl1)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tol2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190705",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	tcl2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190707",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "9950",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tsl2 := []*FuturesTransaction{
		tol2,
		tcl2,
	}

	tradel2, err := NewFuturesTrade(tsl2)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tos1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190704",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10500",
		},
	)

	tcs1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190708",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tss1 := []*FuturesTransaction{
		tos1,
		tcs1,
	}

	trades1, err := NewFuturesTrade(tss1)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tos2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190705",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "9950",
		},
	)

	tcs2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190706",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tss2 := []*FuturesTransaction{
		tos2,
		tcs2,
	}

	trades2, err := NewFuturesTrade(tss2)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	trades := []*FuturesTrade{
		tradel1,
		tradel2,
		trades1,
		trades2,
	}

	s, err := NewStatistic(trades)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	intEqual(t, s.NumberOfTrades(), 4)
	intEqual(t, s.NumberOfWinners(), 2)
	intEqual(t, s.NumberOfLosers(), 2)
	float64Equal(t, s.BattingAverage(), 0.5)
	float64Equal(t, s.BattingAverageL(), 0.5)
	float64Equal(t, s.BattingAverageS(), 0.5)
	float64Equal(t, s.WinLossRatio(), (500*5-1.5*2)/(50*5+1.5*2))
	float64Equal(t, s.WinLossRatioL(), (500*5-1.5*2)/(50*5+1.5*2))
	float64Equal(t, s.WinLossRatioS(), (500*5-1.5*2)/(50*5+1.5*2))
	float64Equal(t, s.AdjWinLossRatio(), (500*5-1.5*2)/(50*5+1.5*2))
	float64Equal(t, s.AdjWinLossRatioL(), (500*5-1.5*2)/(50*5+1.5*2))
	float64Equal(t, s.AdjWinLossRatioS(), (500*5-1.5*2)/(50*5+1.5*2))
	float64Equal(t, s.WGLMax(), 500*5-1.5*2)
	float64Equal(t, s.WGLMean(), 500*5-1.5*2)
	float64Equal(t, s.LGLMax(), -50*5-1.5*2)
	float64Equal(t, s.LGLMean(), -50*5-1.5*2)
	float64Equal(t, s.WGLMeanL(), 500*5-1.5*2)
	float64Equal(t, s.LGLMeanL(), -50*5-1.5*2)
	float64Equal(t, s.WGLMeanS(), 500*5-1.5*2)
	float64Equal(t, s.LGLMeanS(), -50*5-1.5*2)
	intEqual(t, s.HoldingDaysMax(), 8)
	intEqual(t, s.HoldingDaysMaxL(), 8)
	intEqual(t, s.HoldingDaysMaxS(), 4)
	float64Equal(t, s.HoldingDaysMean(), 3.75)
	float64Equal(t, s.HoldingDaysMeanL(), 5)
	float64Equal(t, s.HoldingDaysMeanS(), 2.5)
	float64Equal(t, s.WHoldingDaysMean(), 6)
	float64Equal(t, s.WHoldingDaysMeanL(), 8)
	float64Equal(t, s.WHoldingDaysMeanS(), 4)
	float64Equal(t, s.LHoldingDaysMean(), 1.5)
	float64Equal(t, s.LHoldingDaysMeanL(), 2)
	float64Equal(t, s.LHoldingDaysMeanS(), 1)
	float64Equal(t, s.ExpectedValue(), 2497*0.5-253*0.5)
	float64Equal(t, s.ExpectedValueL(), 2497*0.5-253*0.5)
	float64Equal(t, s.ExpectedValueS(), 2497*0.5-253*0.5)
	float64Equal(t, s.KellyCriterion(), 0.5-(0.5/9.869565))
}

func TestNewStatisticComplex(t *testing.T) {
	tol1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190704",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	tcl1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190712",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10500",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tsl1 := []*FuturesTransaction{
		tol1,
		tcl1,
	}

	tradel1, err := NewFuturesTrade(tsl1)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tol2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190705",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	tcl2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190707",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "9950",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tsl2 := []*FuturesTransaction{
		tol2,
		tcl2,
	}

	tradel2, err := NewFuturesTrade(tsl2)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tol3, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190705",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	tcl3, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190707",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "11000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tsl3 := []*FuturesTransaction{
		tol3,
		tcl3,
	}

	tradel3, err := NewFuturesTrade(tsl3)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tos1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190704",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10500",
		},
	)

	tcs1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190708",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tss1 := []*FuturesTransaction{
		tos1,
		tcs1,
	}

	trades1, err := NewFuturesTrade(tss1)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tos2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190705",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "9950",
		},
	)

	tcs2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190706",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tss2 := []*FuturesTransaction{
		tos2,
		tcs2,
	}

	trades2, err := NewFuturesTrade(tss2)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tos3, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190705",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "9960",
		},
	)

	tcs3, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"date":      "20190706",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tss3 := []*FuturesTransaction{
		tos3,
		tcs3,
	}

	trades3, err := NewFuturesTrade(tss3)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	trades := []*FuturesTrade{
		tradel1,
		tradel2,
		tradel3,
		trades1,
		trades2,
		trades3,
	}

	s, err := NewStatistic(trades)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}
	intEqual(t, s.NumberOfTrades(), 6)
	intEqual(t, s.NumberOfWinners(), 3)
	intEqual(t, s.NumberOfLosers(), 3)
	float64Equal(t, s.BattingAverage(), 0.5)
	float64Equal(t, s.BattingAverageL(), 2.0/3.0)
	float64Equal(t, s.BattingAverageS(), 1.0/3.0)
	float64Equal(t, s.WinLossRatio(), 14.09167842)
	float64Equal(t, s.WinLossRatioL(), 14.81027668)
	float64Equal(t, s.WinLossRatioS(), 10.951754386)
	float64Equal(t, s.AdjWinLossRatio(), 14.09167842)
	float64Equal(t, s.AdjWinLossRatioL(), 29.62055336)
	float64Equal(t, s.AdjWinLossRatioS(), 5.475877193)
	float64Equal(t, s.WGLMax(), 1000*5-1.5*2)
	float64Equal(t, s.LGLMax(), -50*5-1.5*2)
	float64Equal(t, s.WGLMean(), ((500*5-1.5*2)+(1000*5-1.5*2)+(500*5-1.5*2))/3)
	float64Equal(t, s.WGLMeanL(), ((500*5-1.5*2)+(1000*5-1.5*2))/2)
	float64Equal(t, s.WGLMeanS(), 500*5-1.5*2)
	float64Equal(t, s.LGLMean(), ((-50*5-1.5*2)+(-50*5-1.5*2)+(-40*5-1.5*2))/3)
	float64Equal(t, s.LGLMeanL(), -50*5-1.5*2)
	float64Equal(t, s.LGLMeanS(), ((-50*5-1.5*2)+(-40*5-1.5*2))/2)
	intEqual(t, s.HoldingDaysMax(), 8)
	intEqual(t, s.HoldingDaysMaxL(), 8)
	intEqual(t, s.HoldingDaysMaxS(), 4)
	float64Equal(t, s.HoldingDaysMean(), 3)
	float64Equal(t, s.HoldingDaysMeanL(), 4)
	float64Equal(t, s.HoldingDaysMeanS(), 2)
	float64Equal(t, s.WHoldingDaysMean(), 14.0/3.0)
	float64Equal(t, s.WHoldingDaysMeanL(), 5)
	float64Equal(t, s.WHoldingDaysMeanS(), 4)
	float64Equal(t, s.LHoldingDaysMean(), 4.0/3.0)
	float64Equal(t, s.LHoldingDaysMeanL(), 2)
	float64Equal(t, s.LHoldingDaysMeanS(), 1)
	float64Equal(t, s.ExpectedValue(), (2497+2497+4997)/3*0.5-(253+253+203)/3*0.5)
	float64Equal(t, s.ExpectedValueL(), (2497+4997)/2*(2.0/3.0)-253*(1.0/3.0))
	float64Equal(t, s.ExpectedValueS(), 2497*(1.0/3.0)-(253+203)/2*(2.0/3.0))
	float64Equal(t, s.KellyCriterion(), 0.5-(0.5/14.09167842))
}

func intEqual(t *testing.T, a, b int) {
	if fmt.Sprintf("%d", a) != fmt.Sprintf("%d", b) {
		t.Errorf("%d and %d are not equal", a, b)
	}
}

func float64Equal(t *testing.T, a, b float64) {
	if fmt.Sprintf("%.3f", a) != fmt.Sprintf("%.3f", b) {
		t.Errorf("%.3f and %.3f are not equal", a, b)
	}
}
