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
	"github.com/KushamiNeko/go_fun/chart/preset"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

func validInput(symbols, years, freqs, outdir string) error {
	var regex *regexp.Regexp

	regex = regexp.MustCompile(`^\w+(?:,\w+)*$`)
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

	return nil
}

func main() {
	symbols := flag.String("symbols", "", "symbols to plot")
	years := flag.String("years", "", "years to plot")
	freqs := flag.String("freqs", "d,w", "frequency to plot")
	outdir := flag.String("outdir", "", "output directory")
	records := flag.Bool("records", false, "show records")

	flag.Parse()

	err := validInput(*symbols, *years, *freqs, *outdir)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed500, err.Error())
		return
	}

	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("symbols: %s", *symbols))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("years: %s", *years))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("frequency: %s", *freqs))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("outdir: %s", *outdir))
	pretty.ColorPrintln(pretty.PaperPink400, fmt.Sprintf("records: %v", *records))

	sy, ey := parseYears(*years)

	var wg sync.WaitGroup

	for _, symbol := range strings.Split(*symbols, ",") {
		for _, freq := range strings.Split(*freqs, ",") {
			for y := sy; y < ey; y++ {

				name := fmt.Sprintf("%d_%s_%s.png", y, symbol, freq)
				path := filepath.Clean(filepath.Join(*outdir, name))

				wg.Add(1)
				go func(path string, symbol string, year int, freq data.Frequency, records bool) {

					f, err := os.Create(path)
					if err != nil {
						panic(err)
					}

					buffer, err := plotChart(symbol, year, freq, records)
					if err != nil {
						panic(err)
					}

					_, err = io.Copy(f, buffer)
					if err != nil {
						panic(err)
					}

					wg.Done()
				}(path, symbol, y, data.ParseFrequency(freq), *records)

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

func plotChart(symbol string, year int, freq data.Frequency, records bool) (io.Reader, error) {

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

	c, err := preset.NewGeneral(symbol, st.Add(-500*24*time.Hour), et.Add(500*24*time.Hour), freq)
	if err != nil {
		return nil, err
	}

	c.TimeSlice(st, et)

	c.QuoteInfo(false)

	if records {
		err = c.ShowRecordsInBook(fmt.Sprintf("%s_%d", "es", year))
		if err != nil {
			return nil, err
		}
	}

	buffer := new(bytes.Buffer)

	err = c.Plot(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
