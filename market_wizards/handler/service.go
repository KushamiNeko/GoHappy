package handler

import (
	"bytes"
	"fmt"
	"math"
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
	timeFormatL = "20060102150405"
	timeFormatS = "20060102"
)

type cache struct {
	symbol    string
	frequency data.Frequency
	chart     preset.ChartPreset
}

type ServiceHandler struct {
	store []*cache
	//rstore map[string][]*model.FuturesTransaction
	//nstore map[string][]*model.TradingNote

	chart   preset.ChartPreset
	records []*model.FuturesTransaction
	//notes   []*model.TradingNote
}

func (p *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.store == nil {
		p.store = make([]*cache, 0, 6)
	}

	//if p.rstore == nil {
	//p.rstore = make(map[string][]*model.FuturesTransaction)
	//}

	//if p.nstore == nil {
	//p.nstore = make(map[string][]*model.TradingNote)
	//}

	switch r.Method {

	case http.MethodGet:
		var regex *regexp.Regexp

		regex = regexp.MustCompile(`/service/plot/practice\?.+`)
		if regex.MatchString(r.RequestURI) {
			p.getPlot(w, r)
			return
		}

		regex = regexp.MustCompile(`/service/record/note/.+`)
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

	var tfmt string
	regex := regexp.MustCompile(`^\d{8}$`)
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
		p.chart.Forward()

		goto plotting

	case "backward":
		p.chart.Backward()

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
		err = p.recordsLookup(book, symbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//p.chart.ShowRecords(p.records)

		//err = p.notesLookup(book)
		//if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
		//}
	} else {
		//p.chart.ShowRecords(nil)
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

//func (p *ServiceHandler) getNote(w http.ResponseWriter, r *http.Request) {
////if p.series == nil {
////http.NotFound(w, r)
////return
////}

//const pattern = `/record/note/([a-zA-Z0-9_-]+)`
//regex := regexp.MustCompile(pattern)

//m := regex.FindAllStringSubmatch(r.RequestURI, -1)
//if m == nil {
//http.NotFound(w, r)
//return
//}

//book := m[0][1]

//p.notesLookup(book)

//snx := r.URL.Query().Get("x")
//sny := r.URL.Query().Get("y")

//x, _, err := p.inverseXY(snx, sny)
//if err != nil {
//http.Error(w, err.Error(), http.StatusBadRequest)
//return
//}

//dt := p.series.Times()[x]

//for _, n := range p.notes {
//if n.Time().Equal(dt) {
//_, err = w.Write([]byte(n.Note()))
//if err != nil {
//pretty.ColorPrintln(pretty.PaperRed400, err.Error())
//return
//}
//}
//}

//_, err = w.Write([]byte(""))
//if err != nil {
//pretty.ColorPrintln(pretty.PaperRed400, err.Error())
//return
//}
//}

func (p *ServiceHandler) recordsLookup(book, symbol string) error {
	var err error

	tradeAgent, err := agent.NewTradingAgentCompact(
		filepath.Join(
			os.Getenv("HOME"),
			"Documents/database/filedb/futures_wizards",
		),
		"aa",
		book,
	)
	if err != nil {
		return err
	}

	rs, err := tradeAgent.Transactions()
	if err != nil {
		return err
	}

	//p.rstore[book] = rs
	p.records = rs
	//} else {
	//p.records = rs
	//}

	return nil
}

//func (p *ServiceHandler) notesLookup(book string) error {
//if p.nstore == nil {
//p.nstore = make(map[string][]*model.TradingNote)
//}

//root := filepath.Join(os.Getenv("HOME"), "Documents/database/filedb/futures_wizards")

//var (
//data []byte
////err  error

//nse []map[string]string
//ok  bool
//)

//_, ok = p.nstore[book]
//if !ok {

////data, err = ioutil.ReadFile(filepath.Join(root, fmt.Sprintf("%s.yaml", book)))
////if err != nil {
////return err
////}

//tradeAgent, err := agent.NewTradingAgentCompact(
//filepath.Join(
//os.Getenv("HOME"),
//"Documents/database/filedb/futures_wizards",
//),
//"aa",
//book,
//)
//if err != nil {
//return err
//}

//b, err := tradeAgent.Reading()
//if err != nil {
//return err
//}

//data, err = ioutil.ReadFile(filepath.Join(root, fmt.Sprintf("trading_note_%s.json", b.NoteIndex())))
//if err != nil {
//return err
//}

//err = json.Unmarshal(data, &nse)
//if err != nil {
//return err
//}

//ns := make([]*model.TradingNote, 0, len(nse))

//for _, e := range nse {
//n, err := model.NewTradingNoteFromEntity(e)
//if err != nil {
//return err
//}

//ns = append(ns, n)
//}

//sort.Slice(ns, func(i, j int) bool {
//return ns[i].Time().Before(ns[j].Time())
//})

//p.nstore[book] = ns
//}

//p.notes = p.nstore[book]

//return nil
//}

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
