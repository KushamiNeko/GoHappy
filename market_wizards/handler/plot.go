package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	//"github.com/KushamiNeko/go_chart/data"

	"github.com/KushamiNeko/go_fun/chart/data"
)

//const (
//timeFormatLong  = "20060102150405"
//timeFormatShort = "20060102"
//)

type cache struct {
	symbol    string
	frequency data.Frequency

	exstart time.Time
	exend   time.Time

	series *data.TimeSeries
}

type PlotHandler struct {
	store  []*cache
	series *data.TimeSeries
}

func (p *PlotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.store == nil {
		p.store = make([]*cache, 0, 6)
	}

	switch r.Method {

	case http.MethodGet:
		const pattern = `/plot/practice/.+`
		regex := regexp.MustCompile(pattern)

		if !regex.MatchString(r.RequestURI) {
			http.NotFound(w, r)
			return
		}

		p.get(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (p *PlotHandler) get(w http.ResponseWriter, r *http.Request) {
	//var err error

	//const pattern = `/practice/([a-zA-Z0-9]+)/(h|d|w|m)/(time|frequency|forward|backward|info)*/*(\d{14})*/*(records)*/*(\d+)*`
	const pattern = `/practice/([a-zA-Z0-9]+)/(h|d|w|m)/(simple|refresh|forward|backward|info|inspect)*/*(\d{8}|\d{14})*/*(records)*/*(\d+)*`

	regex := regexp.MustCompile(pattern)
	match := regex.FindAllStringSubmatch(r.RequestURI, -1)
	if match == nil {
		http.Error(w, "unknown parameter", http.StatusNotFound)
		return
	}

	//symbol := match[0][1]
	//freq := data.ParseFrequency(match[0][2])
	function := match[0][3]
	//dtime := match[0][4]
	//showRecords := match[0][5] != ""

	//_ = match[0][6] // version

	//src := p.symbolSource(symbol)

	switch function {
	case "simple":
	case "refresh":
	case "forward":
		//p.quotes.Forward()

	case "backward":
		//p.quotes.Backward()
	//case "frequency":
	//case "frequency":
	//dt, err := time.Parse(timeFormatLong, dtime)
	//if err != nil {
	//http.Error(w, fmt.Sprintf("invalid time parameter: %s", dtime), http.StatusNotFound)
	//return
	//}

	//err = p.lookup(src, dt, symbol, freq)
	//if err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//return
	//}
	case "info":
		//if dtime == "" {
		//http.Error(w, "missing time parameter", http.StatusNotFound)
		//return
		//}

		//dt, err := time.Parse(timeFormatLong, dtime)
		//if err != nil {
		//http.Error(w, fmt.Sprintf("invalid time parameter: %s", dtime), http.StatusNotFound)
		//return
		//}

		//err = p.lookup(src, dt, symbol, freq)
		//if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
		//}

		//start, end := p.chartPeriod(dt, freq)

		//p.quotes.TimeSlice(start, end)

	case "inspect":
		//var err error
		//var tf string

		//if freq == data.Hourly {
		//tf = "20060102 15:04:05"
		//} else {
		//tf = "20060102"
		//}

		//msg := fmt.Sprintf(
		//"%s  %s  %s    O:%.2f  H:%.2f  L:%.2f  C:%.2f",
		//p.quotes.SliceEndTime().Format(tf),
		//strings.ToUpper(symbol),
		//strings.ToUpper(freq.FullString()),
		//p.quotes.SliceEndQuote().Open(),
		//p.quotes.SliceEndQuote().High(),
		//p.quotes.SliceEndQuote().Low(),
		//p.quotes.SliceEndQuote().Close(),
		//)

		//if !math.IsNaN(p.quotes.SliceEndQuote().Volume()) && p.quotes.SliceEndQuote().Volume() != 0 {
		//msg = fmt.Sprintf("%s  V:%.0f", msg, p.quotes.SliceEndQuote().Volume())
		//}

		//if !math.IsNaN(p.quotes.SliceEndQuote().OpenInterest()) && p.quotes.SliceEndQuote().OpenInterest() != 0 {
		//msg = fmt.Sprintf("%s  OI:%.0f", msg, p.quotes.SliceEndQuote().OpenInterest())
		//}

		//_, err = w.Write([]byte(msg))
		//if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
		//}

		//return

	default:
		http.Error(w, fmt.Sprintf("unknown function: %s", function), http.StatusNotFound)
		return
	}

	//if showRecords {
	//err = p.recordsLookup(symbol)
	//if err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//return
	//}
	//}

	//var buffer bytes.Buffer

	//err = p.plot(&buffer, p.quotes, freq, showRecords)
	//if err != nil {
	//http.Error(w, err.Error(), http.StatusInternalServerError)
	//return
	//}

	//w.Header().Add("Cache-Control", "no-cache")
	//w.Header().Add("Cache-Control", "no-store")

	//_, err = w.Write(buffer.Bytes())
	//if err != nil {
	//log.Println(err.Error())
	//return
	//}
}

func (p *PlotHandler) lookup(src data.DataSource, dt time.Time, symbol string, freq data.Frequency) error {

	start, end := p.chartPeriod(dt, freq)

	for _, q := range p.store {
		if q.symbol == symbol && q.frequency == freq {

			if (start.After(q.exstart) || start.Equal(q.exstart)) && (end.Before(q.exend) || end.Equal(q.exend)) {
				p.series = q.series
				return nil
			}

		}
	}

	exstart := start.Add(-500 * 24 * time.Hour)
	exend := end.Add(500 * 24 * time.Hour)

	ts, err := src.Read(exstart, exend, symbol, freq)
	if err != nil {
		return err
	}

	//ts.TimeSlice(start, end)

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

func (p *PlotHandler) symbolSource(symbol string) data.DataSource {
	const pattern = `^[a-zA-Z]+$`
	regex := regexp.MustCompile(pattern)

	if regex.MatchString(symbol) {
		return data.NewDataSource(data.Yahoo)
	} else {
		return data.NewDataSource(data.Barchart)
	}
}

func (p *PlotHandler) chartPeriod(end time.Time, freq data.Frequency) (time.Time, time.Time) {
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

//func (p *PlotHandler) recordsLookup(symbol string) error {

//switch symbol[:2] {
//case "es":
//symbol = "es"
//case "qr":
//symbol = "rty"
//case "zn":
//symbol = "ty"
//case "ge":
//symbol = "ed"
//default:
//return fmt.Errorf("unknown symbol: %s", symbol)
//}

//root := filepath.Join(os.Getenv("HOME"), "Documents/database/json/market_wizards")

//files, err := ioutil.ReadDir(root)
//if err != nil {
//return err
//}

//tcss := make([][]*model.FuturesTransaction, 0)
//rds := make([]*data.TradeRecord, 0)

//for _, file := range files {

//if m, _ := regexp.MatchString(`(?:paper|live)_trading_[a-zA-Z0-9]+\.json`, file.Name()); !m {
//continue
//}

//path := filepath.Join(root, file.Name())
//ctn, err := ioutil.ReadFile(path)
//if err != nil {
//return err
//}

//dbs := make([]map[string]string, 0)

//err = json.Unmarshal(ctn, &dbs)
//if err != nil {
//return err
//}

//if len(dbs) == 0 {
//return nil
//}

//tcs := make([]*model.FuturesTransaction, len(dbs))

//for i, db := range dbs {
//tc, err := model.NewFuturesTransactionFromEntity(db)
//if err != nil {
//return err
//}

//tcs[i] = tc
//}

//if symbol != tcs[0].Symbol() {
//continue
//}

//tcss = append(tcss, tcs)
//}

//sort.Slice(tcss, func(i, j int) bool {
//f := tcss[i][int((len(tcss[i])-1)/2.0)].Time()
//s := tcss[j][int((len(tcss[j])-1)/2.0)].Time()

//return f.Year() < s.Year()
//})

//for _, tcs := range tcss {
//for _, tc := range tcs {
//s := (tc.Time().After(p.quotes.SliceStartTime()) || tc.Time().Equal(p.quotes.SliceStartTime()))
//e := (tc.Time().Before(p.quotes.SliceEndTime()) || tc.Time().Equal(p.quotes.SliceEndTime()))

//if s && e {
//rds = append(rds, data.NewTradeRecord(tc.Time(), tc.Operation()))
//}
//}
//}

//p.records = rds

//return nil
//}

//func (p *PlotHandler) plot(out io.Writer, qs *data.Quotes, freq data.Frequency, showRecords bool) error {

//err := qs.NewPlot(config.ColorBackground)
//if err != nil {
//return err
//}

//qs.AddTickMarker(
//&plotter.TimeTicker{
//Quotes:    qs,
//Frequency: freq,
//},
//&plotter.PriceTicker{
//Quotes: qs,
//},
//config.ColorTick,
//config.FontSize,
//)

//qs.AddPlotter(
//utils.GridPlotter(
//utils.PixelsToPoints(config.GridLineWidth),
//config.ColorGrid,
//),
//utils.LinePlotter(
//indicator.SimpleMovingAverge(qs, 5),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorSMA5,
//),
//utils.LinePlotter(
//indicator.SimpleMovingAverge(qs, 20),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorSMA20,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, 1.5),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB15,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, 2),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB20,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, 2.5),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB25,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, 3),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB30,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, -1.5),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB15,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, -2),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB20,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, -2.5),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB25,
//),
//utils.LinePlotter(
//indicator.BollingerBand(qs, 20, -3),
//utils.PixelsToPoints(config.LineWidth),
//config.ColorBB30,
//),
//&plotter.CandleStick{
//Quotes:       qs,
//ColorUp:      config.ColorUp,
//ColorDown:    config.ColorDown,
//ColorNeutral: config.ColorNeutral,
//BodyWidth:    config.CandleBodyWidth,
//ShadowWidth:  config.CandleShadowWidth,
//},
//)

//if p.records != nil && showRecords {
//qs.AddPlotter(
//&plotter.TradesRecorder{
//Quotes:   qs,
//Records:  p.records,
//FontSize: config.FontSize,
//Color:    config.ColorRecord,
//},
//)
//}

//err = qs.Plot(
//out,
//utils.PixelsToPoints(config.ChartWidth),
//utils.PixelsToPoints(config.ChartHeight),
//)
//if err != nil {
//return err
//}

//return nil
//}
