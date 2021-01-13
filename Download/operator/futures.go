package operator

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KushamiNeko/GoFun/Chart/futures"
)

type barchartFutures struct {
	operator

	page    string
	pattern string

	start int
	end   int

	symbols []string
}

func NewBarchartFuturesOperator(start int, end int) *barchartFutures {
	b := &barchartFutures{start: start, end: end}
	b.FromHistoricalPage()

	b.initDir()

	b.symbols = b.source()

	return b
}

func (b *barchartFutures) SetCustomSymbols(symbols []string) *barchartFutures {
	if len(symbols) != 0 {
		b.symbols = symbols
	}

	return b
}

func (b *barchartFutures) FromHistoricalPage() *barchartFutures {
	b.page = historicalPage
	b.pattern = fmt.Sprintf(historicalPattern, symbolPatternFutures)
	return b
}

func (b *barchartFutures) FromInteractivePage() *barchartFutures {
	b.page = interactivePage
	b.pattern = fmt.Sprintf(interactivePattern, symbolPatternFutures)
	return b
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
		"zf",
		"zt",
		"zb",
		"ge",
		"tj",
		"gg",
		"dx",
		"e6",
		"j6",
		"b6",
		"a6",
		"d6",
		"s6",
		"n6",
		"gc",
		"si",
		"cl",
		"ng",
		"zs",
		"zc",
		"zw",
	}
}

func (b *barchartFutures) process(fun func(code string)) {
	for _, symbol := range b.symbols {
		symbol = strings.TrimSpace(symbol)
		symbol = strings.ToLower(symbol)

		startYear := int(b.start / 100)
		endYear := int(b.end / 100)

		for y := startYear; y <= endYear; y++ {
			months := futures.DefaultContractMonths(symbol)

			months.ForEach(func(m futures.MonthCode) {
				t := y*100 + int(futures.MonthCode(m).Month())
				if t >= b.start && t <= b.end {
					code := fmt.Sprintf("%s%s%02d", symbol, string(m), y%100)
					fun(code)
				}
			})
		}
	}
}

func (b *barchartFutures) Download() {
	b.process(func(code string) {
		b.download(fmt.Sprintf(b.page, code), code)
	})

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

func (b *barchartFutures) Check() {
	b.process(func(code string) {
		path := filepath.Join(b.dstDir, "continuous", code[:2], fmt.Sprintf("%s.csv", code))
		b.check(path)
	})

	b.checkCompleted()
}