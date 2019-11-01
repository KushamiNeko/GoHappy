package model

import (
	"fmt"
	"math"

	"github.com/KushamiNeko/futures_wizards/config"
)

type Statistic struct {
	trades     []*FuturesTrade
	categories *categories
}

type categories struct {
	winners []*FuturesTrade
	losers  []*FuturesTrade

	long  []*FuturesTrade
	short []*FuturesTrade

	longWinners []*FuturesTrade
	longLosers  []*FuturesTrade

	shortWinners []*FuturesTrade
	shortLosers  []*FuturesTrade
}

func newCategories(trades []*FuturesTrade) *categories {
	c := new(categories)

	c.winners = make([]*FuturesTrade, 0, len(trades))
	c.losers = make([]*FuturesTrade, 0, len(trades))

	c.long = make([]*FuturesTrade, 0, len(trades))
	c.short = make([]*FuturesTrade, 0, len(trades))

	c.longWinners = make([]*FuturesTrade, 0, len(trades))
	c.longLosers = make([]*FuturesTrade, 0, len(trades))

	c.shortWinners = make([]*FuturesTrade, 0, len(trades))
	c.shortLosers = make([]*FuturesTrade, 0, len(trades))

	for _, t := range trades {
		if t.GL() > 0 {
			c.winners = append(c.winners, t)
		}

		if t.GL() < 0 {
			c.losers = append(c.losers, t)
		}

		if t.Operation() == "+" {
			c.long = append(c.long, t)
		}

		if t.Operation() == "-" {
			c.short = append(c.short, t)
		}

		if t.GL() > 0 && t.Operation() == "+" {
			c.longWinners = append(c.longWinners, t)
		}

		if t.GL() < 0 && t.Operation() == "+" {
			c.longLosers = append(c.longLosers, t)
		}

		if t.GL() > 0 && t.Operation() == "-" {
			c.shortWinners = append(c.shortWinners, t)
		}

		if t.GL() < 0 && t.Operation() == "-" {
			c.shortLosers = append(c.shortLosers, t)
		}
	}

	return c
}

func NewStatistic(trades []*FuturesTrade) (*Statistic, error) {
	if len(trades) == 0 {
		return nil, fmt.Errorf("empty trade")
	}

	s := new(Statistic)
	s.trades = trades
	s.categories = newCategories(trades)

	return s, nil
}

func (s *Statistic) NumberOfTrades() int {
	return len(s.trades)
}

func (s *Statistic) NumberOfWinners() int {
	return len(s.categories.winners)
}

func (s *Statistic) NumberOfLosers() int {
	return len(s.categories.losers)
}

func (s *Statistic) BattingAverage() float64 {
	return float64(len(s.categories.winners)) / float64(len(s.trades))
}

func (s *Statistic) BattingAverageL() float64 {
	return float64(len(s.categories.longWinners)) / float64(len(s.categories.long))
}

func (s *Statistic) BattingAverageS() float64 {
	return float64(len(s.categories.shortWinners)) / float64(len(s.categories.short))
}

func (s *Statistic) WinLossRatio() float64 {
	return s.WGLMean() / math.Abs(s.LGLMean())
}

func (s *Statistic) WinLossRatioL() float64 {
	return s.WGLMeanL() / math.Abs(s.LGLMeanL())
}

func (s *Statistic) WinLossRatioS() float64 {
	return s.WGLMeanS() / math.Abs(s.LGLMeanS())
}

func (s *Statistic) AdjWinLossRatio() float64 {
	return (s.WGLMean() * s.BattingAverage()) / (math.Abs(s.LGLMean()) * (1.0 - s.BattingAverage()))
}

func (s *Statistic) AdjWinLossRatioL() float64 {
	return (s.WGLMeanL() * s.BattingAverageL()) / (math.Abs(s.LGLMeanL()) * (1.0 - s.BattingAverageL()))
}

func (s *Statistic) AdjWinLossRatioS() float64 {
	return (s.WGLMeanS() * s.BattingAverageS()) / (math.Abs(s.LGLMeanS()) * (1.0 - s.BattingAverageS()))
}

func (s *Statistic) WGLMax() float64 {
	var max float64 = 0

	for _, t := range s.categories.winners {
		if t.GL() > max {
			max = t.GL()
		}
	}

	return max
}

func (s *Statistic) WGLMean() float64 {
	var wgl float64 = 0

	for _, t := range s.categories.winners {
		wgl += t.GL()
	}

	return wgl / float64(len(s.categories.winners))
}

func (s *Statistic) WGLMeanL() float64 {
	var wgl float64 = 0

	for _, t := range s.categories.longWinners {
		wgl += t.GL()
	}

	return wgl / float64(len(s.categories.longWinners))
}

func (s *Statistic) WGLMeanS() float64 {
	var wgl float64 = 0

	for _, t := range s.categories.shortWinners {
		wgl += t.GL()
	}

	return wgl / float64(len(s.categories.shortWinners))
}

func (s *Statistic) LGLMax() float64 {
	var max float64 = 0

	for _, t := range s.categories.losers {
		if t.GL() < max {
			max = t.GL()
		}
	}

	return max
}

func (s *Statistic) LGLMean() float64 {
	var lgl float64 = 0

	for _, t := range s.categories.losers {
		lgl += t.GL()
	}

	return lgl / float64(len(s.categories.losers))
}

func (s *Statistic) LGLMeanL() float64 {
	var lgl float64 = 0

	for _, t := range s.categories.longLosers {
		lgl += t.GL()
	}

	return lgl / float64(len(s.categories.longLosers))
}

func (s *Statistic) LGLMeanS() float64 {
	var lgl float64 = 0

	for _, t := range s.categories.shortLosers {
		lgl += t.GL()
	}

	return lgl / float64(len(s.categories.shortLosers))
}

func holdingDays(t *FuturesTrade) int {
	tc := t.CloseDate()
	to := t.OpenDate()

	d := int(math.Ceil(tc.Sub(to).Hours() / 24.0))
	return d
}

func (s *Statistic) HoldingDaysMax() int {
	h := 0

	for _, t := range s.trades {
		d := holdingDays(t)

		if d > h {
			h = d
		}
	}

	return h
}

func (s *Statistic) HoldingDaysMaxL() int {
	h := 0

	for _, t := range s.categories.long {
		d := holdingDays(t)

		if d > h {
			h = d
		}
	}

	return h
}

func (s *Statistic) HoldingDaysMaxS() int {
	h := 0

	for _, t := range s.categories.short {
		d := holdingDays(t)

		if d > h {
			h = d
		}
	}

	return h
}

func (s *Statistic) HoldingDaysMean() float64 {
	d := 0

	for _, t := range s.trades {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.trades))
}

func (s *Statistic) HoldingDaysMeanL() float64 {
	d := 0

	for _, t := range s.categories.long {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.long))
}

func (s *Statistic) HoldingDaysMeanS() float64 {
	d := 0

	for _, t := range s.categories.short {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.short))
}

func (s *Statistic) WHoldingDaysMean() float64 {
	d := 0

	for _, t := range s.categories.winners {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.winners))
}

func (s *Statistic) WHoldingDaysMeanL() float64 {
	d := 0

	for _, t := range s.categories.longWinners {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.longWinners))
}

func (s *Statistic) WHoldingDaysMeanS() float64 {
	d := 0

	for _, t := range s.categories.shortWinners {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.shortWinners))
}

func (s *Statistic) LHoldingDaysMean() float64 {
	d := 0

	for _, t := range s.categories.losers {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.losers))
}

func (s *Statistic) LHoldingDaysMeanL() float64 {
	d := 0

	for _, t := range s.categories.longLosers {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.longLosers))
}

func (s *Statistic) LHoldingDaysMeanS() float64 {
	d := 0

	for _, t := range s.categories.shortLosers {
		d += holdingDays(t)
	}

	return float64(d) / float64(len(s.categories.shortLosers))
}

func (s *Statistic) ExpectedValue() float64 {
	return (s.WGLMean() * s.BattingAverage()) + (s.LGLMean() * (1.0 - s.BattingAverage()))
}

func (s *Statistic) ExpectedValueL() float64 {
	return (s.WGLMeanL() * s.BattingAverageL()) + (s.LGLMeanL() * (1.0 - s.BattingAverageL()))
}

func (s *Statistic) ExpectedValueS() float64 {
	return (s.WGLMeanS() * s.BattingAverageS()) + (s.LGLMeanS() * (1.0 - s.BattingAverageS()))
}

func (s *Statistic) KellyCriterion() float64 {
	return s.BattingAverage() - ((1.0 - s.BattingAverage()) / s.WinLossRatio())
}

func (s *Statistic) Entity() map[string]string {
	return map[string]string{
		"number_of_trades":        fmt.Sprintf("%d", s.NumberOfTrades()),
		"number_of_winners":       fmt.Sprintf("%d", s.NumberOfWinners()),
		"number_of_losers":        fmt.Sprintf("%d", s.NumberOfLosers()),
		"batting_average":         fmt.Sprintf("%.2f", s.BattingAverage()),
		"win_loss_ratio":          fmt.Sprintf("%.2f", s.WinLossRatio()),
		"adjusted_win_loss_ratio": fmt.Sprintf("%.2f", s.AdjWinLossRatio()),
		"expected_value":          fmt.Sprintf("%.2f", s.ExpectedValue()),
		//"kelly_criterion":         fmt.Sprintf("%.2f", s.KellyCriterion()),
	}
}

const (
	//statisticFmtString = "%-[3]*[4]s%-[3]*[5]s%-[3]*[6]s%-[3]*[7]s%-[3]*[8]s%-[3]*[9]s"
	statisticFmtString = "%-[3]*[4]s%-[3]*[5]s%-[3]*[6]s%-[3]*[7]s%-[3]*[8]s"
)

func (s *Statistic) Fmt() string {
	return fmt.Sprintf(
		statisticFmtString,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		fmt.Sprintf("%d", s.NumberOfTrades()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.BattingAverage()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.WinLossRatio()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.AdjWinLossRatio()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.ExpectedValue()),
		//fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.KellyCriterion()),
	)
}

func StatisticFmtLabels() string {
	return fmt.Sprintf(
		statisticFmtString,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		"Number of Trades",
		"Batting Average",
		"Win Loss Ratio",
		"Adj WL Ratio",
		"Expected Value",
		//"Kelly Criterion",
	)
}

const (
	statisticFmtStringL = "%-[3]*[4]s%-[3]*[5]s%-[3]*[6]s%-[3]*[7]s%-[3]*[8]s%-[3]*[9]s%-[3]*[10]s"
)

func (s *Statistic) FmtL() string {
	return fmt.Sprintf(
		statisticFmtStringL,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		fmt.Sprintf("%s", "Long"),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.BattingAverageL()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.WinLossRatioL()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.AdjWinLossRatioL()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.WHoldingDaysMeanL()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.LHoldingDaysMeanL()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.ExpectedValueL()),
	)
}

func StatisticFmtLabelsL() string {
	return fmt.Sprintf(
		statisticFmtStringL,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		"Operation",
		"Batting Average",
		"Win Loss Ratio",
		"Adj WL Ratio",
		"W Holding Mean",
		"L Holding Mean",
		"Expected Value",
	)
}

const (
	statisticFmtStringS = "%-[3]*[4]s%-[3]*[5]s%-[3]*[6]s%-[3]*[7]s%-[3]*[8]s%-[3]*[9]s%-[3]*[10]s"
)

func (s *Statistic) FmtS() string {
	return fmt.Sprintf(
		statisticFmtStringS,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		fmt.Sprintf("%s", "Short"),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.BattingAverageS()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.WinLossRatioS()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.AdjWinLossRatioS()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.WHoldingDaysMeanS()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.LHoldingDaysMeanS()),
		fmt.Sprintf("%.[1]*f", config.FloatDecimals, s.ExpectedValueS()),
	)
}

func StatisticFmtLabelsS() string {
	return fmt.Sprintf(
		statisticFmtStringS,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		"Operation",
		"Batting Average",
		"Win Loss Ratio",
		"Adj WL Ratio",
		"W Holding Mean",
		"L Holding Mean",
		"Expected Value",
	)
}
