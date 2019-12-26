package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/KushamiNeko/go_fun/chart/data"
	"github.com/KushamiNeko/go_fun/trading/model"
)

var (
	series  *data.TimeSeries
	records []*model.FuturesTransaction
)

const (
	year = 2018
)

func init() {
	var err error

	src := data.NewDataSource(data.StockCharts)

	from, err := time.Parse("20060102", fmt.Sprintf("%d0101", year))
	if err != nil {
		panic(err)
	}

	from = from.Add(-500 * 24 * time.Hour)

	to, err := time.Parse("20060102", fmt.Sprintf("%d0101", year+1))
	if err != nil {
		panic(err)
	}

	to = to.Add(500 * 24 * time.Hour)

	series, err = src.Read(from, to, "es", data.Daily)
	if err != nil {
		panic(err)
	}

	agent, err := newTradingAgent(fmt.Sprintf("es_%d", year))
	if err != nil {
		panic(err)
	}

	records, err = agent.Transactions()
	if err != nil {
		panic(err)
	}
}

func transactionSlice(f, t time.Time) []*model.FuturesTransaction {
	ts := make([]*model.FuturesTransaction, 0)

	sliced := false
	for _, r := range records {
		if (r.Time().Equal(f) || r.Time().After(f)) && !sliced {
			sliced = true
		}

		if (r.Time().After(t)) && sliced {
			sliced = false
		}

		if sliced {
			ts = append(ts, r)
		}
	}

	if len(ts) > 1 && ts[0].Time().Equal(ts[1].Time()) {
		ts = ts[1:]
	}

	if len(ts) > 1 && ts[len(ts)-1].Time().Equal(ts[len(ts)-2].Time()) {
		ts = ts[:len(ts)-1]
	}

	return ts
}

func formatRisk(label string, t time.Time, risk float64) string {
	return fmt.Sprintf("%-15s @ %s: %.4f%%", label, t.Format("20060102"), risk)
}

func risk(from, to, op string, adjusted bool) string {
	f, err := time.Parse("20060102", from)
	if err != nil {
		panic(err)
	}

	t, err := time.Parse("20060102", to)
	if err != nil {
		panic(err)
	}

	ts := transactionSlice(f, t)

	s := make([]string, 0)
	positions := make([]float64, 0)
	sizes := make([]float64, 0)

	for _, record := range ts {

		var risk float64

		index := series.TimesIndex(record.Time())
		if index == -1 {
			panic(fmt.Errorf("unknown index"))
		}

		nl := series.ValueAtTimesIndex(index+1, "low", 0)
		nh := series.ValueAtTimesIndex(index+1, "high", 0)
		c := series.ValueAtTimesIndex(index, "close", 0)

		//fmt.Println(c)
		//fmt.Println(nl)
		//fmt.Println(nh)

		var q float64
		var label string
		if adjusted {
			q = float64(record.Quantity())
			label = "Risk(adj)"
		} else {
			q = 1
			label = "Risk(sim)"
		}

		switch op {
		case "+":
			r := ((nl - c) / c) * 100.0
			risk += r * q

			s = append(s, formatRisk(label, record.Time(), risk))

			if len(positions) > 0 {
				for i, p := range positions {
					r := ((nl - p) / p) * 100.0
					risk += r * sizes[i]
				}

				s = append(s, formatRisk(fmt.Sprintf("Total %s", label), record.Time(), risk))
			}

		case "-":
			r := ((c - nh) / c) * 100.0
			risk += r * q

			s = append(s, formatRisk(label, record.Time(), risk))

			if len(positions) > 0 {
				for i, p := range positions {
					r := ((p - nh) / p) * 100.0
					risk += r * sizes[i]
				}

				s = append(s, formatRisk(fmt.Sprintf("Total %s", label), record.Time(), risk))
			}

		default:
			panic(fmt.Errorf("unknown op"))
		}

		positions = append(positions, c)
		sizes = append(sizes, q)
	}

	return strings.TrimSpace(strings.Join(s, "\n"))
}
