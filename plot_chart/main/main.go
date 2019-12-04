package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KushamiNeko/go_fun/chart/data"
	"github.com/KushamiNeko/go_fun/chart/indicator"
	"github.com/KushamiNeko/go_fun/chart/plot"
	"github.com/KushamiNeko/go_fun/chart/plotter"
	"github.com/KushamiNeko/go_fun/chart/utils"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

func validInput(symbols, years, freqs, outdir, version string) error {
	var regex *regexp.Regexp
	var err error

	regex = regexp.MustCompile(`^[a-zA-Z]+\d*(?:,\w+\d*)*$`)
	if !regex.MatchString(symbols) {
		return fmt.Errorf("invalid symbols: %s", symbols)
	}

	regex = regexp.MustCompile(`^(?:h|d|w|m)(?:,(?:h|d|w|m))*$`)
	if !regex.MatchString(freqs) {
		return fmt.Errorf("invalid freqs: %s", freqs)
	}

	regex = regexp.MustCompile(`^(\d{4})(?:(?:\-|\~)(\d{4}))*$`)
	if !regex.MatchString(years) {
		return fmt.Errorf("invalid years: %s", years)
	}

	if outdir != "" {
		_, err := os.Stat(outdir)
		if err != nil {
			return fmt.Errorf("invalid outdir: %s", outdir)
		}
	}

	regex = regexp.MustCompile(`^[a-zA-Z]\d$`)
	if err != nil {
		return fmt.Errorf("invalid version: %s", version)
	}

	return nil
}

func main() {
	symbols := flag.String("symbols", "", "symbols to plot")
	years := flag.String("years", "", "years to plot")
	freqs := flag.String("freqs", "d,w", "frequency to plot")
	outdir := flag.String("outdir", "", "output directory")
	version := flag.String("version", "", "records version")
	records := flag.Bool("records", false, "show records")

	flag.Parse()

	err := validInput(*symbols, *years, *freqs, *outdir, *version)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed500, err.Error())
		return
	}

	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("symbols: %s", *symbols))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("years: %s", *years))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("frequency: %s", *freqs))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("outdir: %s", *outdir))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("version: %s", *version))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("records: %v", *records))

	sy, ey := parseYears(*years)

	var wg sync.WaitGroup

	for _, symbol := range strings.Split(*symbols, ",") {
		for _, freq := range strings.Split(*freqs, ",") {
			for y := sy; y < ey; y++ {

				s := symbol
				regex := regexp.MustCompile(`^([a-zA-Z]+)([fghjkmnquvxz][0-9]+)?$`)
				m := regex.FindAllStringSubmatch(symbol, -1)
				if m != nil {
					s = m[0][1]
				}

				name := fmt.Sprintf("%d_%s_%s.png", y, s, freq)
				path := filepath.Clean(filepath.Join(*outdir, name))

				wg.Add(1)
				go func(path string, symbol string, year int, freq data.Frequency, version string, records bool) {

					f, err := os.Create(path)
					if err != nil {
						panic(err)
					}

					buffer, err := plotChart(symbol, year, freq, version, records)
					if err != nil {
						panic(err)
					}

					_, err = io.Copy(f, buffer)
					if err != nil {
						panic(err)
					}

					wg.Done()
				}(path, symbol, y, data.ParseFrequency(freq), *version, *records)

			}
		}
	}

	wg.Wait()
}

func parseYears(years string) (int, int) {
	var err error

	regex := regexp.MustCompile(`^(\d{4})(?:(?:\-|\~)(\d{4}))*$`)
	m := regex.FindAllStringSubmatch(years, -1)

	s := m[0][1]
	e := m[0][2]

	var sy, ey int64

	sy, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid year: %s", s))
	}

	if e == "" {
		ey = sy + 1
	} else {
		ey, err = strconv.ParseInt(e, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid year: %s", e))
		}

		ey += 1
	}

	if ey <= sy {
		panic(fmt.Sprintf("invalid years: %d - %d", sy, ey))
	}

	return int(sy), int(ey)
}

func plotChart(symbol string, year int, freq data.Frequency, version string, records bool) (io.Reader, error) {

	var src data.DataSource

	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if regex.MatchString(symbol) {
		src = data.NewDataSource(data.Yahoo)
	} else {
		src = data.NewDataSource(data.Barchart)
	}

	var st, et time.Time
	loc, err := time.LoadLocation("")
	if err != nil {
		panic(err)
	}

	et = time.Date(
		year+1,
		time.January,
		1,
		0,
		0,
		0,
		0,
		loc,
	)

	switch freq {
	case data.Hourly:
		panic("unimplemented")
	case data.Daily:
		st = time.Date(
			year,
			time.January,
			1,
			0,
			0,
			0,
			0,
			loc,
		)
	case data.Weekly:
		st = time.Date(
			year-3,
			time.January,
			1,
			0,
			0,
			0,
			0,
			loc,
		)
	case data.Monthly:
		panic("unimplemented")
	default:
		panic("unknown frequency")
	}

	ts, err := src.Read(st, et, symbol, freq)
	if err != nil {
		return nil, err
	}

	ts.SetColumn("sma5", indicator.SimpleMovingAverge(ts.FullValues("close"), 5))
	ts.SetColumn("sma20", indicator.SimpleMovingAverge(ts.FullValues("close"), 20))

	ts.SetColumn("bb+15", indicator.BollingerBand(ts.FullValues("close"), 20, 1.5))
	ts.SetColumn("bb-15", indicator.BollingerBand(ts.FullValues("close"), 20, -1.5))
	ts.SetColumn("bb+20", indicator.BollingerBand(ts.FullValues("close"), 20, 2.0))
	ts.SetColumn("bb-20", indicator.BollingerBand(ts.FullValues("close"), 20, -2.0))
	ts.SetColumn("bb+25", indicator.BollingerBand(ts.FullValues("close"), 20, 2.5))
	ts.SetColumn("bb-25", indicator.BollingerBand(ts.FullValues("close"), 20, -2.5))
	ts.SetColumn("bb+30", indicator.BollingerBand(ts.FullValues("close"), 20, 3.0))
	ts.SetColumn("bb-30", indicator.BollingerBand(ts.FullValues("close"), 20, -3.0))

	p := &plot.Plot{}
	err = p.Init()
	if err != nil {
		return nil, err
	}

	p.AddTickMarker(
		&plotter.TimeTicker{
			TimeSeries: ts,
			Frequency:  freq,
		},
		&plotter.PriceTicker{
			TimeSeries: ts,
			Step:       20,
		},
		plot.ThemeColor("ColorTick"),
		plot.ChartConfig("TickFontSize"),
	)

	p.AddPlotter(
		plotter.GridPlotter(
			plot.ChartConfig("GridLineWidth"),
			plot.ThemeColor("ColorGrid"),
		),
		plotter.LinePlotter(
			ts.Values("sma5"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorSMA1"),
		),
		plotter.LinePlotter(
			ts.Values("sma20"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorSMA2"),
		),
		plotter.LinePlotter(
			ts.Values("bb+15"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB1"),
		),
		plotter.LinePlotter(
			ts.Values("bb-15"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB1"),
		),
		plotter.LinePlotter(
			ts.Values("bb+20"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB2"),
		),
		plotter.LinePlotter(
			ts.Values("bb-20"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB2"),
		),
		plotter.LinePlotter(
			ts.Values("bb+25"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB3"),
		),
		plotter.LinePlotter(
			ts.Values("bb-25"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB3"),
		),
		plotter.LinePlotter(
			ts.Values("bb+30"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB4"),
		),
		plotter.LinePlotter(
			ts.Values("bb-30"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB4"),
		),
		&plotter.CandleStick{
			TimeSeries:   ts,
			ColorUp:      plot.ThemeColor("ColorUp"),
			ColorDown:    plot.ThemeColor("ColorDown"),
			ColorNeutral: plot.ThemeColor("ColorNeutral"),
			BodyWidth:    plot.ChartConfig("CandleBodyWidth"),
			ShadowWidth:  plot.ChartConfig("CandleShadowWidth"),
		},
	)

	if records {
		//p.AddPlotter(
		//&plotter.LeverageRecorder{
		//TimeSeries: ts,
		//Records:    rs,
		//FontSize:   plot.ChartConfig("RecordsFontSize"),
		//Color:      plot.ThemeColor("ColorText"),
		//},
		//)
	}

	ymn, ymx := utils.RangeExtend(utils.Min(ts.Values("low")), utils.Max(ts.Values("high")), 25.0)
	p.YRange(ymn, ymx)

	buffer := new(bytes.Buffer)

	err = p.Plot(
		buffer,
		plot.ChartConfig("ChartWidth"),
		plot.ChartConfig("ChartHeight"),
	)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
