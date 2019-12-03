package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/KushamiNeko/go_fun/chart/data"
	"github.com/KushamiNeko/go_fun/chart/indicator"
	"github.com/KushamiNeko/go_fun/chart/plot"
	"github.com/KushamiNeko/go_fun/chart/plotter"
	"github.com/KushamiNeko/go_fun/chart/utils"
	"github.com/KushamiNeko/go_fun/trading/model"
	"github.com/KushamiNeko/go_fun/utils/pretty"
	"gonum.org/v1/plot/vg"
	"gopkg.in/yaml.v2"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	plot.SmallChart()
}

const (
	timeFormatL = "20060102150405"
	timeFormatS = "20060102"
)

type cache struct {
	symbol    string
	frequency data.Frequency

	exstart time.Time
	exend   time.Time

	series *data.TimeSeries
}

type ServiceHandler struct {
	store  []*cache
	rstore map[string][]*model.FuturesTransaction
	nstore map[string][]*model.TradingNote

	series  *data.TimeSeries
	records []*model.FuturesTransaction
	notes   []*model.TradingNote

	tickWidth  vg.Length
	tickHeight vg.Length
}

func (p *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.store == nil {
		p.store = make([]*cache, 0, 6)
	}

	if p.rstore == nil {
		p.rstore = make(map[string][]*model.FuturesTransaction)
	}

	if p.nstore == nil {
		p.nstore = make(map[string][]*model.TradingNote)
	}

	switch r.Method {

	case http.MethodGet:
		var regex *regexp.Regexp

		regex = regexp.MustCompile(`/service/plot/practice/.+`)
		if regex.MatchString(r.RequestURI) {
			p.getPlot(w, r)
			return
		}

		regex = regexp.MustCompile(`/service/record/note/.+`)
		if regex.MatchString(r.RequestURI) {
			p.getNote(w, r)
			return
		}

		http.NotFound(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (p *ServiceHandler) getPlot(w http.ResponseWriter, r *http.Request) {
	const pattern = `/plot/practice/([a-zA-Z0-9]+)/(h|d|w|m)/(simple|refresh|forward|backward|info|inspect)/*(\d{8}|\d{14})*/*(records)*/*([a-zA-Z0-9_-]+)*`

	regex := regexp.MustCompile(pattern)
	match := regex.FindAllStringSubmatch(r.RequestURI, -1)
	if match == nil {
		http.Error(w, "unknown parameter", http.StatusNotFound)
		return
	}

	symbol := match[0][1]
	freq := data.ParseFrequency(match[0][2])
	function := match[0][3]
	dtime := match[0][4]
	showRecords := match[0][5] != ""

	book := match[0][6]

	var tfmt string
	regex = regexp.MustCompile(`^\d{8}$`)
	if regex.MatchString(dtime) {
		tfmt = timeFormatS
	} else {
		tfmt = timeFormatL
	}

	dt, err := time.Parse(tfmt, dtime)
	if err != nil {
		panic(err)
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
		p.series.Forward()

		goto plotting

	case "backward":
		p.series.Backward()

		goto plotting

	case "info":
		if p.series == nil {
			http.NotFound(w, r)
			return
		}

		data := struct {
			Time  string
			Open  float64
			High  float64
			Low   float64
			Close float64
			//Volume   float64
			//Interest float64
		}{
			Time:  p.series.EndTime().Format(timeFormatS),
			Open:  p.series.EndValue("open"),
			High:  p.series.EndValue("high"),
			Low:   p.series.EndValue("low"),
			Close: p.series.EndValue("close"),
			//Volume:   p.series.EndValue("volume"),
			//Interest: p.series.EndValue("openinterest"),
		}

		jd, err := json.Marshal(&data)
		if err != nil {
			pretty.ColorPrintln(pretty.PaperRed400, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jd)
		if err != nil {
			pretty.ColorPrintln(pretty.PaperRed400, err.Error())
		}

		return
	case "inspect":
		if p.series == nil {
			http.NotFound(w, r)
			return
		}

		snx := r.URL.Query().Get("x")
		sny := r.URL.Query().Get("y")

		x, y, err := p.inverseXY(snx, sny)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		msn := message.NewPrinter(message.MatchLanguage("en"))

		is := fmt.Sprintf(
			"time: %s\nprice: %s\nopen: %s\nhigh: %s\nlow: %s\nclose: %s\nvolume: %s\ninterest: %s\n",
			p.series.Times()[x].Format("2006-01-02"),
			msn.Sprintf("%.2f", y),
			msn.Sprintf("%.2f", p.series.ValueAtIndex(x, "open", 0)),
			msn.Sprintf("%.2f", p.series.ValueAtIndex(x, "high", 0)),
			msn.Sprintf("%.2f", p.series.ValueAtIndex(x, "low", 0)),
			msn.Sprintf("%.2f", p.series.ValueAtIndex(x, "close", 0)),
			msn.Sprintf("%.0f", p.series.ValueAtIndex(x, "volume", 0)),
			msn.Sprintf("%.0f", p.series.ValueAtIndex(x, "openinterest", 0)),
		)

		sanx := r.URL.Query().Get("ax")
		sany := r.URL.Query().Get("ay")

		if sanx != "" && sany != "" {
			ax, ay, err := p.inverseXY(sanx, sany)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			dt := p.series.Times()[x]
			da := p.series.Times()[ax]

			d := dt.Sub(da)
			days := int(math.Ceil(d.Hours() / 24))

			is = fmt.Sprintf(
				"%sdiff(days): %s\ndiff($): %s\ndiff(%%): %s\n",
				is,
				msn.Sprintf("%d", days),
				msn.Sprintf("%.2f", y-ay),
				msn.Sprintf("%.2f", ((y-ay)/ay)*100.0),
			)
		} else {
			if x > 0 {
				dd := p.series.ValueAtIndex(x, "close", 0) - p.series.ValueAtIndex(x-1, "close", 0)
				dp := (dd / p.series.ValueAtIndex(x-1, "close", 0)) * 100.0

				is = fmt.Sprintf(
					"%sdiff($): %s\ndiff(%%): %s\n",
					is,
					msn.Sprintf("%.2f", dd),
					msn.Sprintf("%.2f", dp),
				)
			}
		}

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
		err = p.recordsLookup(book, symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = p.notesLookup(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var buffer bytes.Buffer

	err = p.plot(&buffer, freq, showRecords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Cache-Control", "no-store")

	_, err = w.Write(buffer.Bytes())
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed400, err.Error())
		return
	}
}

func (p *ServiceHandler) getNote(w http.ResponseWriter, r *http.Request) {
	if p.series == nil {
		http.NotFound(w, r)
		return
	}

	const pattern = `/record/note/([a-zA-Z0-9_-]+)`
	regex := regexp.MustCompile(pattern)

	m := regex.FindAllStringSubmatch(r.RequestURI, -1)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	book := m[0][1]

	p.notesLookup(book)

	snx := r.URL.Query().Get("x")
	sny := r.URL.Query().Get("y")

	x, _, err := p.inverseXY(snx, sny)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dt := p.series.Times()[x]

	for _, n := range p.notes {
		if n.Time().Equal(dt) {
			_, err = w.Write([]byte(n.Note()))
			if err != nil {
				pretty.ColorPrintln(pretty.PaperRed400, err.Error())
				return
			}
		}
	}

	_, err = w.Write([]byte(""))
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed400, err.Error())
		return
	}
}

func (p *ServiceHandler) recordsLookup(book, symbol string) error {
	root := filepath.Join(os.Getenv("HOME"), "Documents/database/filedb")

	var (
		data []byte
		err  error

		rse []map[string]string
		ok  bool
	)

	_, ok = p.rstore[book]
	if !ok {
		data, err = ioutil.ReadFile(filepath.Join(root, fmt.Sprintf("%s.json", book)))
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &rse)
		if err != nil {
			return err
		}

		rs := make([]*model.FuturesTransaction, 0, len(rse))

		for _, e := range rse {
			t, err := model.NewFuturesTransactionFromEntity(e)
			if err != nil {
				return err
			}

			rs = append(rs, t)
		}

		sort.Slice(rs, func(i, j int) bool {
			if rs[i].Time().Equal(rs[j].Time()) {
				return rs[i].TimeStamp() < rs[j].TimeStamp()
			} else {
				return rs[i].Time().Before(rs[j].Time())
			}
		})

		p.rstore[book] = rs
	}

	rs := make([]*model.FuturesTransaction, 0, len(p.rstore[book]))
	for _, r := range p.rstore[book] {
		if strings.Contains(symbol, r.Symbol()) {
			rs = append(rs, r)
		}
	}

	p.records = rs

	return nil
}

func (p *ServiceHandler) notesLookup(book string) error {
	root := filepath.Join(os.Getenv("HOME"), "Documents/database/filedb")

	var (
		data []byte
		err  error

		nse []map[string]string
		ok  bool
	)

	_, ok = p.nstore[book]
	if !ok {
		data, err = ioutil.ReadFile(filepath.Join(root, fmt.Sprintf("%s.yaml", book)))
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, &nse)
		if err != nil {
			return err
		}

		ns := make([]*model.TradingNote, 0, len(nse))

		for _, e := range nse {
			n, err := model.NewTradingNoteFromEntity(e)
			if err != nil {
				return err
			}

			ns = append(ns, n)
		}

		sort.Slice(ns, func(i, j int) bool {
			return ns[i].Time().Before(ns[j].Time())
		})

		p.nstore[book] = ns
	}

	p.notes = p.nstore[book]

	return nil
}

func (p *ServiceHandler) inverseX(min, max, nx float64) float64 {
	/*

			linear scale function

		  func (LinearScale) Normalize(min, max, x float64) float64 {
				return (x - min) / (max - min)
		  }

			math of inverse linear scale

		  y = (x - min) / (max - min)
		  y * (max - min) = (x - min)
		  (y * (max - min)) + min = x

	*/

	//const wm = 0.032817628
	wm := float64(p.tickWidth / plot.ChartConfig("ChartWidth"))

	r := (max - min) / (1.0 - wm)
	return (nx - wm) * r
}

func (p *ServiceHandler) inverseY(min, max, ny float64) float64 {

	/*

		log scale function

		func (LogScale) Normalize(min, max, x float64) float64 {
			if min <= 0 || max <= 0 || x <= 0 {
				panic("Values must be greater than 0 for a log scale.")
			}

			logMin := math.Log(min)
			return (math.Log(x) - logMin) / (math.Log(max) - logMin)
		}

		math of inverse log scale

		y = (log x - log min) / (log max - log min)
		y * (log max - log min) = (log x - log min)
		(y * (log max - log min)) + log min = log x
		e ^ ((y * (log max - log min)) + log min) = x

	*/

	//const hm = 0.025
	hm := float64(p.tickHeight / plot.ChartConfig("ChartHeight"))

	r := 1.0 / (1.0 - hm)
	ly := (ny - hm) * r

	logMin := math.Log(min)
	logMax := math.Log(max)

	return math.Exp(
		(ly * (logMax - logMin)) + logMin,
	)
}

func (p *ServiceHandler) inverseXY(snx, sny string) (int, float64, error) {
	nx, err := strconv.ParseFloat(snx, 64)
	if err != nil {
		return 0, 0, err
	}

	ny, err := strconv.ParseFloat(sny, 64)
	if err != nil {
		return 0, 0, err
	}

	ymin, ymax := utils.RangeExtend(
		utils.Min(p.series.Values("low")),
		utils.Max(p.series.Values("high")),
		25.0,
	)

	x := int(math.Round(p.inverseX(-0.5, float64(len(p.series.Times()))-0.5, nx)))
	if x < 0 {
		x = 0
	} else if x > len(p.series.Times())-1 {
		x = len(p.series.Times()) - 1
	}

	y := p.inverseY(ymin, ymax, ny)
	if y < ymin {
		y = ymin
	} else if y > ymax {
		y = ymax
	}

	return x, y, nil
}

func (p *ServiceHandler) lookup(dt time.Time, symbol string, freq data.Frequency, timeSliced bool) error {
	start, end := p.chartPeriod(dt, freq)

	for _, q := range p.store {
		if q.symbol == symbol && q.frequency == freq {

			if (start.After(q.exstart) || start.Equal(q.exstart)) && (end.Before(q.exend) || end.Equal(q.exend)) {
				p.series = q.series

				if timeSliced {
					p.series.TimeSlice(start, end)
				}

				return nil
			}

		}
	}

	src := p.symbolSource(symbol)

	exstart := start.Add(-500 * 24 * time.Hour)
	exend := end.Add(500 * 24 * time.Hour)

	ts, err := src.Read(exstart, exend, symbol, freq)
	if err != nil {
		return err
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

	ts.TimeSlice(start, end)

	p.store = append(
		p.store,
		&cache{
			symbol:    symbol,
			frequency: freq,
			series:    ts,
			exstart:   exstart,
			exend:     exend,
		},
	)

	p.series = ts

	return nil
}

func (p *ServiceHandler) symbolSource(symbol string) data.DataSource {
	const pattern = `^[a-zA-Z]+$`
	regex := regexp.MustCompile(pattern)

	if regex.MatchString(symbol) {

		switch {
		case (symbol == "spx" ||
			symbol == "compq" ||
			symbol == "rut" ||
			symbol == "vix" ||
			symbol == "vxn" ||
			symbol == "ovx" ||
			symbol == "gvz"):
			return data.NewDataSource(data.Yahoo)
		case symbol == "vle" || symbol == "tyvix" || symbol == "rvx":
			return data.NewDataSource(data.StockCharts)
		default:
			return data.NewDataSource(data.AlphaVantage)
		}

	} else {
		return data.NewDataSource(data.Barchart)
	}
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

func (p *ServiceHandler) plot(out io.Writer, freq data.Frequency, showRecords bool) error {

	pt := &plot.Plot{}
	err := pt.Init()
	if err != nil {
		return err
	}

	pt.AddTickMarker(
		&plotter.TimeTicker{
			TimeSeries: p.series,
			Frequency:  freq,
		},
		&plotter.PriceTicker{
			TimeSeries: p.series,
			Step:       20,
		},
		plot.ThemeColor("ColorTick"),
		plot.ChartConfig("TickFontSize"),
	)

	pt.AddPlotter(
		plotter.GridPlotter(
			plot.ChartConfig("GridLineWidth"),
			plot.ThemeColor("ColorGrid"),
		),
		plotter.LinePlotter(
			p.series.Values("sma5"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorSMA1"),
		),
		plotter.LinePlotter(
			p.series.Values("sma20"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorSMA2"),
		),
		plotter.LinePlotter(
			p.series.Values("bb+15"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB1"),
		),
		plotter.LinePlotter(
			p.series.Values("bb-15"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB1"),
		),
		plotter.LinePlotter(
			p.series.Values("bb+20"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB2"),
		),
		plotter.LinePlotter(
			p.series.Values("bb-20"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB2"),
		),
		plotter.LinePlotter(
			p.series.Values("bb+25"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB3"),
		),
		plotter.LinePlotter(
			p.series.Values("bb-25"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB3"),
		),
		plotter.LinePlotter(
			p.series.Values("bb+30"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB4"),
		),
		plotter.LinePlotter(
			p.series.Values("bb-30"),
			plot.ChartConfig("LineWidth"),
			plot.ThemeColor("ColorBB4"),
		),
		&plotter.CandleStick{
			TimeSeries:   p.series,
			ColorUp:      plot.ThemeColor("ColorUp"),
			ColorDown:    plot.ThemeColor("ColorDown"),
			ColorNeutral: plot.ThemeColor("ColorNeutral"),
			BodyWidth:    plot.ChartConfig("CandleBodyWidth"),
			ShadowWidth:  plot.ChartConfig("CandleShadowWidth"),
		},
		&plotter.QuoteInfo{
			TimeSeries: p.series,
			XOffset:    5,
			FontSize:   plot.ChartConfig("InfoFontSize"),
			Color:      plot.ThemeColor("ColorText"),
		},
	)

	if showRecords {
		pt.AddPlotter(
			&plotter.TradesRecorder{
				TimeSeries: p.series,
				Frequency:  freq,
				Records:    p.records,
				Notes:      p.notes,
				FontSize:   plot.ChartConfig("RecordsFontSize"),
				Color:      plot.ThemeColor("ColorText"),
			},
		)
	}

	mh := utils.Max(p.series.Values("high"))
	ml := utils.Min(p.series.Values("low"))

	ymin, ymax := utils.RangeExtend(
		ml,
		mh,
		25.0,
	)

	pt.YRange(ymin, ymax)

	msg := message.NewPrinter(language.English)

	p.tickWidth = pt.TickWidth(msg.Sprintf("%.2f", pt.YMax()))
	p.tickHeight = pt.TickHeight("2006 Jan")

	err = pt.Plot(
		out,
		plot.ChartConfig("ChartWidth"),
		plot.ChartConfig("ChartHeight"),
	)
	if err != nil {
		return err
	}

	return nil
}
