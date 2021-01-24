package operator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/KushamiNeko/GoFun/Chart/futures"
)

type coinAPI struct {
	*operator

	start int
	end   int

	symbols []string
}

func NewCoinAPI(start int, end int) *coinAPI {
	c := &coinAPI{
		operator: new(operator),
		start:    start,
		end:      end,
	}

	c.initDir()

	c.symbols = c.source()

	return c
}

func (c *coinAPI) SetCustomSymbols(symbols []string) *coinAPI {
	if len(symbols) != 0 {
		c.symbols = symbols
	}

	return c
}

func (b *coinAPI) source() []string {
	return []string{
		"BTCUSD",
		"ETHUSD",
		"LTCUSD",
		"XRPUSD",
		"BCHUSD",
		"USDCUSD",
		"DOTUSD",
		"LINKUSD",
		"EOSUSD",
		"ADAUSD",
		"XLMUSD",
		"TRXUSD",
		"UNIUSD",
		// "BNBUSD",
		// "DOGEUSD",

	}
}

func (c *coinAPI) nextContract(ym string) string {
	const pattern = `^(\d{4})(\d{2})$`
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(ym) {
		panic(fmt.Sprintf("invalid year month format: %s", ym))
	}

	match := regex.FindAllStringSubmatch(ym, -1)

	ys := match[0][1]
	ms := match[0][2]

	year, err := strconv.Atoi(ys)
	if err != nil {
		panic(err)
	}

	month, err := strconv.Atoi(ms)
	if err != nil {
		panic(err)
	}

	idt := (year * 100) + month

	for i := year; i <= year+1; i++ {
		for _, m := range futures.CryptoContractMonths {
			ic := (i * 100) + futures.MonthCode(m).MonthValue()
			if ic > idt {
				return fmt.Sprintf("%d", ic)
			}
		}
	}

	panic("logic error in next contract function")
}

func (c *coinAPI) url(symbol, startYM string) string {
	const baseUrl = `https://rest.coinapi.io/v1/ohlcv/%s/%s/history?period_id=%s`

	endYM := c.nextContract(startYM)

	const (
		startTime = "00:00:00"
		endTime   = "00:00:00"
		period    = "1DAY"
	)

	start := fmt.Sprintf("%s-%s-01", startYM[:4], startYM[4:6])
	end := fmt.Sprintf("%s-%s-01", endYM[:4], endYM[4:6])

	const pattern = `^(\w+?)(\w{3})$`
	regex := regexp.MustCompile(pattern)

	symbol = strings.ToUpper(symbol)

	if !regex.MatchString(symbol) {
		panic(fmt.Sprintf("invalid symbol: %s", symbol))
	}

	match := regex.FindAllStringSubmatch(symbol, -1)
	crypto := match[0][1]
	base := match[0][2]

	var b strings.Builder
	b.WriteString(
		fmt.Sprintf(
			baseUrl,
			strings.ToUpper(crypto),
			strings.ToUpper(base),
			strings.ToUpper(period),
		),
	)

	b.WriteString(fmt.Sprintf("&time_start=%sT%s", start, startTime))
	b.WriteString(fmt.Sprintf("&time_end=%sT%s", end, endTime))
	//b.WriteString("&limit=100000")

	return b.String()
}

func (c *coinAPI) process(fun func(symbol, startYM string)) {
	for _, symbol := range c.symbols {
		symbol = strings.TrimSpace(symbol)
		symbol = strings.ToUpper(symbol)

		startYear := int(c.start / 100)
		endYear := int(c.end / 100)

		for y := startYear; y <= endYear; y++ {
			months := futures.CryptoContractMonths

			months.ForEach(func(m futures.MonthCode) {
				t := y*100 + int(futures.MonthCode(m).Month())
				if t >= c.start && t <= c.end {
					fun(symbol, fmt.Sprintf("%d%02d", y, m.MonthValue()))
				}
			})
		}
	}
}

func (c *coinAPI) Download() {
	apiKey := os.Getenv("COIN_API")
	if apiKey == "" {
		panic("invalid coin api key")
	}

	c.process(func(symbol, startYM string) {
		url := c.url(symbol, startYM)

		symbol = strings.ToLower(symbol)

		path := filepath.Join(c.dstDir, "coinapi", symbol)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic(fmt.Sprintf("path does not exist: %s", path))
		}

		file := fmt.Sprintf("%s@%s.json", symbol, startYM)
		path = filepath.Join(path, file)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("X-CoinAPI-Key", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = f.Write(body)
		if err != nil {
			panic(err)
		}

		var b strings.Builder
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("src: %s", url))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("dst: %s", path))

		c.downloadMessage(b.String())
		c.downloadCount++

		c.checkProcessLimit()
	})

	c.downloadCompleted()
}

func (c *coinAPI) Rename() {
	c.renameCount = c.downloadCount
	c.renameCompleted()
}

func (c *coinAPI) Check() {
	c.process(func(symbol, startYM string) {
		symbol = strings.ToLower(symbol)

		path := filepath.Join(c.dstDir, "coinapi", symbol)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic(fmt.Sprintf("path does not exist: %s", path))
		}

		file := fmt.Sprintf("%s@%s.json", symbol, startYM)
		path = filepath.Join(path, file)

		c.check(path)
	})

	c.checkCompleted()
}
