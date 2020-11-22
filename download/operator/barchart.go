package operator

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KushamiNeko/go_fun/chart/futures"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

const (
	historicalPage = `https://www.barchart.com/futures/quotes/%s/historical-download`

	symbolPatternGeneral = `([\$\^][\w]+)`
	symbolPatternFutures = `([\w\d]{5})`

	//historicalPattern = `^([\w\d]{5})_([^_-]+)(?:-[^_-]+)*_[^_-]+-[^_-]+-\d{2}-\d{2}-\d{4}.csv$`
	//historicalPattern = `^([\$\^][\w]+)_([^_-]+)(?:-[^_-]+)*_[^_-]+-[^_-]+-\d{2}-\d{2}-\d{4}.csv$`
	historicalPattern = `^%s_([^_-]+)(?:-[^_-]+)*_[^_-]+-[^_-]+-\d{2}-\d{2}-\d{4}.csv$`

	interactivePage = `https://www.barchart.com/futures/quotes/%s/interactive-chart`
	//interactivePattern = `^([\w\d]{5})_[^_]+_[^_]+_[^_]+_([^_]+)(?:_[^_]+)*_\d{2}_\d{2}_\d{4}.csv$`
	interactivePattern = `^%s_[^_]+_[^_]+_[^_]+_([^_]+)(?:_[^_]+)*_\d{2}_\d{2}_\d{4}.csv$`

	futuresPageLimit = 10
)

type barchartFutures struct {
	operator

	page    string
	pattern string

	start int
	end   int

	symbols []string
	months  string
}

func NewBarchartFuturesOperator(start int, end int) *barchartFutures {
	b := &barchartFutures{start: start, end: end}
	b.FromHistoricalPage()

	b.initDir()

	b.symbols = b.source()

	return b
}

func NewBarchartFuturesOperatorCustom(start int, end int, symbols, months string) *barchartFutures {
	b := NewBarchartFuturesOperator(start, end)

	b.symbols = strings.Split(symbols, ",")
	b.months = months

	return b
}

func (b *barchartFutures) FromHistoricalPage() {
	b.page = historicalPage
	b.pattern = fmt.Sprintf(historicalPattern, symbolPatternFutures)
}

func (b *barchartFutures) FromInteractivePage() {
	b.page = interactivePage
	b.pattern = fmt.Sprintf(interactivePattern, symbolPatternFutures)
}

func (b *barchartFutures) source() []string {
	return []string{
		"es",
		"nq",
		"qr",
		"ym",
		"np",
		"fx",
		"zn",
		"ge",
		"tj",
		"gg",
		"dx",
		"j6",
		"e6",
		"b6",
		"a6",
		"n6",
		"d6",
		"s6",
		"gc",
		"cl",
	}
}

func (b *barchartFutures) Download() {
	count := 0

	//for _, symbol := range b.source() {
	for _, symbol := range b.symbols {
		symbol = strings.TrimSpace(symbol)

		startYear := int(b.start / 100)
		endYear := int(b.end / 100)

		for y := startYear; y <= endYear; y++ {
			var months futures.ContractMonths
			//switch symbol {
			//case "cl":
			//months = futures.AllContractMonths
			//case "gc":
			//months = futures.EvenContractMonths
			//default:
			//months = futures.FinancialContractMonths
			//}

			if b.months != "" {
				months = futures.ContractMonths(b.months)
			} else {
				switch symbol {
				case "cl":
					months = futures.AllContractMonths
				case "gc":
					months = futures.EvenContractMonths
				default:
					months = futures.FinancialContractMonths
				}
			}

			for _, m := range months {
				t := y*100 + int(futures.MonthCode(m).Month())
				if t >= b.start && t <= b.end {
					code := fmt.Sprintf("%s%s%02d", symbol, string(m), y%100)
					b.download(fmt.Sprintf(b.page, code), code)
				}

			}

		}

		count += 1
		if count >= futuresPageLimit {
			b.completed()
			count = 0
		}
	}

	b.downloadCompleted()
}

func (b *barchartFutures) Rename() {
	regex := regexp.MustCompile(b.pattern)

	fs, err := ioutil.ReadDir(b.srcDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		match := regex.FindAllStringSubmatch(f.Name(), -1)
		if len(match) != 0 {
			code := strings.ToLower(match[0][1])

			srcPath := filepath.Join(b.srcDir, f.Name())
			dstPath := filepath.Join(b.dstDir, "continuous", code[:2], fmt.Sprintf("%s.csv", code))

			b.rename(srcPath, dstPath)
		}

	}

	b.renameCompleted()
}

func (b *barchartFutures) Check() {}

type barchartGeneral struct {
	operator

	page    string
	pattern string
}

func NewBarchartGeneralOperator() *barchartGeneral {
	b := &barchartGeneral{}
	b.FromHistoricalPage()

	b.initDir()

	return b
}

func (b *barchartGeneral) FromHistoricalPage() {
	b.page = historicalPage
	b.pattern = fmt.Sprintf(historicalPattern, symbolPatternGeneral)
}

func (b *barchartGeneral) FromInteractivePage() {
	b.page = interactivePage
	b.pattern = fmt.Sprintf(interactivePattern, symbolPatternGeneral)
}

func (b *barchartGeneral) source() map[string]string {
	return map[string]string{
		// "$iqx": "spxew",
		// "$slew": "smlew",
		// "$sdew": "midew",
		// "$topx": "topix",
		// "$addn": "addn",
		// "$addq": "addq",
		// "$avdn": "avdn",
		// "$avdq": "avdq",
		// "$addt": "addt",
		// "$avdt": "avdt",
		"^btcusd": "btcusd",
		"^ethusd": "ethusd",
		"^ltcusd": "ltcusd",
		"^xrpusd": "xrpusd",
		"$dxy":    "dxy",
		"^eurusd": "eurusd",
		"^usdjpy": "usdjpy",
		"^gbpusd": "gbpusd",
		"^audusd": "audusd",
		"^usdcad": "usdcad",
		"^usdchf": "usdchf",
		"^nzdusd": "nzdusd",
		"^eurjpy": "eurjpy",
		"^eurgbp": "eurgbp",
		"^euraud": "euraud",
		"^eurcad": "eurcad",
		"^eurchf": "eurchf",
		"^gbpjpy": "gbpjpy",
		"^audjpy": "audjpy",
		"^cadjpy": "cadjpy",
	}
}

func (b *barchartGeneral) Download() {
	for symbol := range b.source() {
		b.download(fmt.Sprintf(b.page, symbol), symbol)
	}

	b.downloadCompleted()
}

func (b *barchartGeneral) Rename() {
	regex := regexp.MustCompile(b.pattern)

	fs, err := ioutil.ReadDir(b.srcDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		match := regex.FindAllStringSubmatch(f.Name(), -1)
		if len(match) != 0 {

			symbol := strings.ToLower(match[0][1])
			symbol, ok := b.source()[symbol]
			if !ok {
				pretty.ColorPrintln(
					pretty.PaperPink300,
					fmt.Sprintf("barchart general operator skips renaming symbol %s", symbol))
				continue
			}

			srcPath := filepath.Join(b.srcDir, f.Name())
			dstPath := filepath.Join(b.dstDir, "barchart", fmt.Sprintf("%s.csv", symbol))

			b.rename(srcPath, dstPath)
		}

	}

	b.renameCompleted()
}

func (b *barchartGeneral) Check() {}
