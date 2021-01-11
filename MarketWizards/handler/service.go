package handler

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/KushamiNeko/go_fun/chart/data"
	"github.com/KushamiNeko/go_fun/chart/preset"
	"github.com/KushamiNeko/go_fun/trading/agent"
	"github.com/KushamiNeko/go_fun/trading/model"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

func init() {
	preset.SmallChart()
}

const (
	//timeFormatL = "20060102150405"
	//timeFormatS = "20060102"
	timeFormat = "20060102"
)

type cache struct {
	symbol    string
	frequency data.Frequency
	chart     preset.ChartPreset
}

type ServiceHandler struct {
	store []*cache
	chart preset.ChartPreset
}

func (p *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.store == nil {
		p.store = make([]*cache, 0, 6)
	}

	switch r.Method {

	case http.MethodGet:
		var regex *regexp.Regexp

		regex = regexp.MustCompile(`/service/plot/practice\?.+`)
		if regex.MatchString(r.RequestURI) {
			p.getPlot(w, r)
			return
		}

		regex = regexp.MustCompile(`/service/record/note\?.+`)
		if regex.MatchString(r.RequestURI) {
			//p.getNote(w, r)
			return
		}

		http.NotFound(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
func (p *ServiceHandler) getPlot(w http.ResponseWriter, r *http.Request) {

	symbol := r.URL.Query().Get("symbol")
	freq := data.ParseFrequency(r.URL.Query().Get("frequency"))
	function := r.URL.Query().Get("function")
	dtime := r.URL.Query().Get("time")
	showRecords := r.URL.Query().Get("records") == "true"

	book := r.URL.Query().Get("book")

	regex := regexp.MustCompile(`^\d{4}$`)
	if regex.MatchString(dtime) {
		dtime = fmt.Sprintf("%s1231", dtime)
	} else {
		regex = regexp.MustCompile(`^\d{8}$`)
		if !regex.MatchString(dtime) {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}
	}

	dt, err := time.Parse(timeFormat, dtime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch function {
	case "simple":
		err = p.lookup(dt, symbol, freq, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		goto plotting

	case "refresh":
		err = p.lookup(dt, symbol, freq, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		goto plotting

	case "forward":
		p.chart.Forward()

		goto plotting

	case "backward":
		p.chart.Backward()

		goto plotting

	case "randomTrade":
		ta, err := agent.NewTradingAgentCompact(
			filepath.Join(
				os.Getenv("HOME"),
				"Documents/database/filedb/futures_wizards",
			),
			"aa",
			"",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bs, err := ta.Books()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tbs := make([]*model.TradingBook, 0, len(bs))
		bre := regexp.MustCompile(`^\w{2}_\d{4}$`)
		for _, b := range bs {
			if bre.MatchString(b.Title()) {
				tbs = append(tbs, b)
			}
		}

		rand.Seed(time.Now().Unix())

		//ri := rand.Intn(len(bs))
		ri := rand.Intn(len(tbs))

		//ta.SetReading(bs[ri])
		ta.SetReading(tbs[ri])

		ts, err := ta.Trades()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		times := make([]time.Time, 0, len(ts)+1)

		for _, t := range ts {
			times = append(times, t.OpenTime())
		}

		ps, _ := ta.Positions()
		if len(ps) != 0 {
			times = append(times, ps[0].Time())
		}

		ri = rand.Intn(len(times))

		st := times[ri]

		dt = st.Add(-1 * time.Duration(rand.Intn(9)+1) * 7 * 24 * time.Hour)

		err = p.lookup(dt, symbol, freq, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		goto plotting

	case "randomDate":

		rand.Seed(time.Now().Unix())

		var (
			y int
			m time.Month
			d int
		)

		y = 1999 + rand.Intn(time.Now().Year()-1999)

		switch y {
		case 1999:
			m = time.Month(rand.Intn(9) + 4)
		default:
			m = time.Month(rand.Intn(12) + 1)
		}

		switch m {
		case time.January:
			d = rand.Intn(31) + 1
		case time.February:
			d = rand.Intn(28) + 1
		case time.March:
			d = rand.Intn(31) + 1
		case time.April:
			d = rand.Intn(30) + 1
		case time.May:
			d = rand.Intn(31) + 1
		case time.June:
			d = rand.Intn(30) + 1
		case time.July:
			d = rand.Intn(31) + 1
		case time.August:
			d = rand.Intn(31) + 1
		case time.September:
			d = rand.Intn(30) + 1
		case time.October:
			d = rand.Intn(31) + 1
		case time.November:
			d = rand.Intn(30) + 1
		case time.December:
			d = rand.Intn(31) + 1
		}

		dt = time.Date(
			y,
			m,
			d,
			0,
			0,
			0,
			0,
			time.Now().Location(),
		)

		err = p.lookup(dt, symbol, freq, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		goto plotting

	case "info":
		if p.chart == nil {
			http.NotFound(w, r)
			return
		}

		d, err := p.chart.LatestInfo()
		if err != nil {
			pretty.ColorPrintln(pretty.PaperRed400, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(d.Bytes())
		if err != nil {
			pretty.ColorPrintln(pretty.PaperRed400, err.Error())
			return
		}

		return
	case "inspect":
		if p.chart == nil {
			http.NotFound(w, r)
			return
		}

		snx := r.URL.Query().Get("x")
		sny := r.URL.Query().Get("y")

		nx, err := strconv.ParseFloat(snx, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ny, err := strconv.ParseFloat(sny, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sanx := r.URL.Query().Get("ax")
		sany := r.URL.Query().Get("ay")

		var anx, any float64
		if sanx == "" || sany == "" {
			anx = math.NaN()
			any = math.NaN()
		} else {
			anx, err = strconv.ParseFloat(sanx, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			any, err = strconv.ParseFloat(sany, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		is := p.chart.Inspect(nx, ny, anx, any)

		_, err = w.Write([]byte(is))
		if err != nil {
			pretty.ColorPrintln(pretty.PaperRed500, err.Error())
			return
		}

		return
	default:
		http.Error(w, fmt.Sprintf("unknown function: %s", function), http.StatusNotFound)
		return
	}

plotting:

	if showRecords {
		err = p.chart.ShowRecordsInBook(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		p.chart.ShowRecords(nil)
	}

	var buffer bytes.Buffer

	p.chart.Plot(&buffer)

	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Cache-Control", "no-store")

	_, err = w.Write(buffer.Bytes())
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed400, err.Error())
		return
	}
}

func (p *ServiceHandler) lookup(dt time.Time, symbol string, freq data.Frequency, timeSliced bool) error {

	start, end := p.chartPeriod(dt, freq)

	for _, q := range p.store {
		if q.symbol == symbol && q.frequency == freq {
			p.chart = q.chart

			if timeSliced {
				p.chart.TimeSlice(start, end)
			}

			return nil
		}
	}

	c, err := preset.NewChartPreset(symbol, start, end, freq)
	if err != nil {
		return err
	}

	p.store = append(
		p.store,
		&cache{
			symbol:    symbol,
			frequency: freq,
			chart:     c,
		},
	)

	p.chart = c

	return nil
}

func (p *ServiceHandler) chartPeriod(end time.Time, freq data.Frequency) (time.Time, time.Time) {
	var s time.Time

	switch freq {
	case data.Hourly:
		ne := end.Add(-1 * 15 * 24 * time.Hour)
		s = time.Date(
			ne.Year(),
			ne.Month(),
			ne.Day(),
			ne.Hour(),
			ne.Minute(),
			ne.Second(),
			ne.Nanosecond(),
			ne.Location(),
		)
	case data.Daily:
		s = time.Date(
			end.Year()-1,
			end.Month(),
			end.Day(),
			end.Hour(),
			end.Minute(),
			end.Second(),
			end.Nanosecond(),
			end.Location(),
		)
	case data.Weekly:
		s = time.Date(
			end.Year()-4,
			end.Month(),
			end.Day(),
			end.Hour(),
			end.Minute(),
			end.Second(),
			end.Nanosecond(),
			end.Location(),
		)
	case data.Monthly:
		s = time.Date(
			end.Year()-18,
			end.Month(),
			end.Day(),
			end.Hour(),
			end.Minute(),
			end.Second(),
			end.Nanosecond(),
			end.Location(),
		)
	default:
		panic(fmt.Sprintf("unknown frequency: %s", freq))
	}

	return s, end
}
