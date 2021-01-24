package operator

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KushamiNeko/GoFun/Utils/pretty"
)

const (
	//symbolPatternGeneral = `([\$\^][\w]+)`
	symbolPatternGeneral = `([\$\^]*[a-zA-Z0-9.]+)`
	symbolPatternFutures = `([\w\d]{5})`

	historicalPage    = `https://www.barchart.com/futures/quotes/%s/historical-download`
	historicalPattern = `^%s_([^_-]+)(?:-[^_-]+)*_[^_-]+-[^_-]+-\d{2}-\d{2}-\d{4}.csv$`

	interactivePage    = `https://www.barchart.com/futures/quotes/%s/interactive-chart`
	interactivePattern = `^%s_[^_]+_[^_]+_[^_]+_([^_]+)(?:_[^_]+)*_\d{2}_\d{2}_\d{4}.csv$`
)

type barchartGeneral struct {
	*operator

	page    string
	pattern string
}

func NewBarchartGeneralOperator() *barchartGeneral {
	b := &barchartGeneral{operator: new(operator)}
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
		//
		"$dxy":    "dxy",
		"^eurusd": "eurusd",
		"^usdjpy": "usdjpy",
		"^gbpusd": "gbpusd",
		"^audusd": "audusd",
		"^usdcad": "usdcad",
		"^usdchf": "usdchf",
		"^nzdusd": "nzdusd",
		//
		//"^eurjpy": "eurjpy",
		//"^eurgbp": "eurgbp",
		//"^euraud": "euraud",
		//"^eurcad": "eurcad",
		//"^eurchf": "eurchf",
		//"^gbpjpy": "gbpjpy",
		//"^audjpy": "audjpy",
		//"^cadjpy": "cadjpy",
		//
		// "fedfunds.rt": "fedfunds",
		// "ustm1.rt": "ustm1",
		"ustm3.rt": "ustm3",
		// "ustm6.rt": "ustm6",
		"usty2.rt": "usty2",
		// "usty5.rt": "usty5",
		"usty10.rt": "usty10",
		// "usty30.rt": "usty30",
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

func (b *barchartGeneral) Check() {
	for _, symbol := range b.source() {
		path := filepath.Join(b.dstDir, "barchart", fmt.Sprintf("%s.csv", symbol))
		b.check(path)
	}

	b.checkCompleted()
}
