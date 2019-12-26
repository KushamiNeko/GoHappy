package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/KushamiNeko/go_fun/chart/data"
	"github.com/KushamiNeko/go_fun/trading/agent"
	"github.com/KushamiNeko/go_fun/trading/model"
	"github.com/KushamiNeko/go_fun/trading/utils"
	"github.com/KushamiNeko/go_fun/utils/pretty"
	"gonum.org/v1/gonum/stat"
)

func processInput(symbol, period, op, book string) (*agent.TradingAgent, time.Time, time.Time, error) {
	var regex *regexp.Regexp

	regex = regexp.MustCompile(`^es|nq|qr|zn$`)
	if !regex.MatchString(symbol) {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("invalid symbol: %s", symbol)
	}

	regex = regexp.MustCompile(`^(\d{4}|\d{8})(?:\s*[-~]\s*(\d{4}|\d{8}))?$`)
	if !regex.MatchString(period) {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("invalid period: %s", period)
	}

	var (
		f, t time.Time
		err  error
	)
	m := regex.FindAllStringSubmatch(period, -1)
	if len(m) != 1 {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("invalid period: %s", period)
	} else {
		if m[0][2] == "" {
			f, err = time.Parse("20060102", period)
			if err != nil {
				return nil, time.Time{}, time.Time{}, err
			}

			t = f
		} else {
			f, err = time.Parse("20060102", m[0][1])
			if err != nil {
				return nil, time.Time{}, time.Time{}, err
			}

			t, err = time.Parse("20060102", m[0][2])
			if err != nil {
				return nil, time.Time{}, time.Time{}, err
			}
		}
	}

	regex = regexp.MustCompile(`^+|-$`)
	if !regex.MatchString(op) {
		return nil, time.Time{}, time.Time{}, fmt.Errorf("invalid op: %s", op)
	}

	a, err := agent.NewTradingAgentCompact(
		filepath.Join(
			os.Getenv("HOME"),
			"Documents/database/filedb/futures_wizards",
		),
		"aa",
		book,
	)
	if err != nil {
		return nil, time.Time{}, time.Time{}, err
	}

	return a, f, t, nil
}

func timeExtend(f, t time.Time) (time.Time, time.Time) {

	nf := time.Date(
		f.Year(),
		time.January,
		1,
		0,
		0,
		0,
		0,
		f.Location(),
	)

	nt := time.Date(
		t.Year(),
		time.December,
		31,
		0,
		0,
		0,
		0,
		t.Location(),
	)

	return nf, nt
}

func main() {
	symbol := flag.String("symbol", "", "symbol to calculate")
	period := flag.String("period", "", "time period")
	op := flag.String("op", "", "operation")
	book := flag.String("book", "", "records book")

	flag.Parse()

	ta, f, t, err := processInput(*symbol, *period, *op, *book)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed500, err.Error())
		return
	}

	records, err := ta.Transactions()
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed500, err.Error())
		return
	}

	nf, nt := timeExtend(f, t)

	src := data.NewDataSource(data.StockCharts)
	ysrc := data.NewDataSource(data.Yahoo)

	series, err := src.Read(nf, nt, *symbol, data.Daily)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed500, err.Error())
		return
	}

	vix, err := ysrc.Read(nf, nt, "vix", data.Daily)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed500, err.Error())
		return
	}

	vxn, err := ysrc.Read(nf, nt, "vxn", data.Daily)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed500, err.Error())
		return
	}

	calculateRisk(
		series,
		records,
		vix,
		vxn,
		f,
		t,
		*op,
	)

	fmt.Println()

	volatilityAverage(vix, vxn, f, t)
}

func calculateRisk(
	series *data.TimeSeries,
	records []*model.FuturesTransaction,
	vix *data.TimeSeries,
	vxn *data.TimeSeries,
	f, t time.Time,
	op string,
) {

	srs := utils.CalculateRisk(
		series,
		records,
		f,
		t,
		op,
		false,
	)

	ars := utils.CalculateRisk(
		series,
		records,
		f,
		t,
		op,
		true,
	)

	if len(srs) != len(ars) {
		panic(fmt.Sprintf("length should be the same"))
	}

	msr := math.Inf(1)
	mar := math.Inf(1)

	cmsr := math.Inf(1)
	cmar := math.Inf(1)

	var ct time.Time
	for i := 0; i < len(srs); i++ {
		if !srs[i].Combined() {
			msr = math.Min(srs[i].Risk(), msr)
		} else {
			cmsr = math.Min(srs[i].Risk(), cmsr)
		}

		if !ars[i].Combined() {
			mar = math.Min(ars[i].Risk(), mar)
		} else {
			cmar = math.Min(ars[i].Risk(), cmar)
		}

		if ct.IsZero() || !srs[i].Time().Equal(ct) {
			if !ct.IsZero() {
				fmt.Println()
			}

			ct = srs[i].Time()
			pretty.ColorPrintln(pretty.PaperPink300, srs[i].Time().Format("20060102"))

			pretty.ColorPrintln(pretty.PaperBlue300, fmt.Sprintf("VIX: %.2f", vix.ValueInTimes(ct, "close", 0)))
			pretty.ColorPrintln(pretty.PaperPurple200, fmt.Sprintf("VXN: %.2f", vxn.ValueInTimes(ct, "close", 0)))
		}

		pretty.ColorPrintln(pretty.PaperCyan300, fmt.Sprintf("%s: %.4f%%", srs[i].Label(), srs[i].Risk()))
		pretty.ColorPrintln(pretty.PaperYellow300, fmt.Sprintf("%s: %.4f%%", ars[i].Label(), ars[i].Risk()))

	}

	fmt.Println()
	pretty.ColorPrintln(pretty.PaperPink300, "Maximum Risks")
	pretty.ColorPrintln(pretty.PaperCyan500, fmt.Sprintf("%s: %.4f%%", "Simple", msr))
	pretty.ColorPrintln(pretty.PaperYellow500, fmt.Sprintf("%s: %.4f%%", "Adjusted", mar))
	pretty.ColorPrintln(pretty.PaperCyan500, fmt.Sprintf("%s: %.4f%%", "Combined Simple", cmsr))
	pretty.ColorPrintln(pretty.PaperYellow500, fmt.Sprintf("%s: %.4f%%", "Combined Adjusted", cmar))
}

func volatilityAverage(
	vix *data.TimeSeries,
	vxn *data.TimeSeries,
	f, t time.Time,
) {

	vix.TimeSlice(f, t)
	vxn.TimeSlice(f, t)

	pretty.ColorPrintln(pretty.PaperBlue300, fmt.Sprintf("Average VIX: %.2f", stat.Mean(vix.Values("close"), nil)))
	pretty.ColorPrintln(pretty.PaperPurple200, fmt.Sprintf("Average VXN: %.2f", stat.Mean(vxn.Values("close"), nil)))
}
