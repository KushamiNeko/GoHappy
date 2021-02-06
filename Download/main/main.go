package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/KushamiNeko/GoFun/Utils/input"
	"github.com/KushamiNeko/GoFun/Utils/pretty"
	"github.com/KushamiNeko/GoHappy/Download/operator"
)

const introColor = pretty.PaperYellow300
const errorColor = pretty.PaperRed500
const introSeparator = ", "

const symbolSeparator = `[,&\s]`

var symbolsRegex = regexp.MustCompile(symbolSeparator)

func main() {

	var futuresStart, futuresEnd string
	var cryptoStart, cryptoEnd string

	var futuresDaily, futuresIntraday60, futuresIntraday30, barchart, yahoo, investing, coinapi bool
	var download, rename, check bool

	var futuresSymbols string
	var cryptoSymbols string

	flag.StringVar(&futuresStart, "futures-start", fmt.Sprintf("%d", (time.Now().Year()-1)*100+int(time.Now().Month())), "yyyymm")
	flag.StringVar(&futuresEnd, "futures-end", fmt.Sprintf("%d", (time.Now().Year())*100+int(time.Now().Month())), "yyyymm")

	flag.StringVar(&futuresSymbols, "futures-symbols", "", "custom futures symbol list")

	flag.StringVar(&cryptoStart, "crypto-start", fmt.Sprintf("%d", (time.Now().Year()-1)*100+int(time.Now().Month())), "yyyymm")
	flag.StringVar(&cryptoEnd, "crypto-end", fmt.Sprintf("%d", (time.Now().Year())*100+int(time.Now().Month())), "yyyymm")

	flag.StringVar(&cryptoSymbols, "crypto-symbols", "", "custom crypto symbol list")

	flag.BoolVar(&futuresDaily, "futures", false, "download futures daily data from Barchart")
	flag.BoolVar(&futuresIntraday60, "futures-intraday-60", false, "download futures intraday data from Barchart")
	flag.BoolVar(&futuresIntraday30, "futures-intraday-30", false, "download futures intraday data from Barchart")
	flag.BoolVar(&barchart, "barchart", false, "download data from Barchart")
	flag.BoolVar(&yahoo, "yahoo", false, "download from Yahoo")
	flag.BoolVar(&investing, "investing", false, "download from Investing.com")
	flag.BoolVar(&coinapi, "coinapi", false, "download from CoinAPI")

	flag.BoolVar(&download, "download", false, "download operation")
	flag.BoolVar(&rename, "rename", false, "rename operation")
	flag.BoolVar(&check, "check", false, "check operation")

	flag.Parse()

	const pattern = `^\d{6}$`

	if futuresDaily && !input.ValidateWithRegex(futuresStart, pattern) {
		pretty.ColorPrintln(errorColor, "invalid futures start")
		return
	}

	if futuresDaily && !input.ValidateWithRegex(futuresEnd, pattern) {
		pretty.ColorPrintln(errorColor, "invalid futures end")
		return
	}

	if coinapi && !input.ValidateWithRegex(cryptoStart, pattern) {
		pretty.ColorPrintln(errorColor, "invalid crypto start")
		return
	}

	if coinapi && !input.ValidateWithRegex(cryptoEnd, pattern) {
		pretty.ColorPrintln(errorColor, "invalid crypto end")
		return
	}

	if !futuresDaily && !futuresIntraday60 && !futuresIntraday30 && !barchart && !yahoo && !investing && !coinapi {
		futuresDaily = true
		futuresIntraday60 = true
		futuresIntraday30 = true
		barchart = true
		yahoo = true
		investing = true
	}

	if !download && !rename && !check {
		download = true
		rename = true
		check = false
	}

	if futuresDaily || futuresIntraday60 || futuresIntraday30 {
		if futuresStart != "" && futuresEnd != "" {
			pretty.ColorPrintln(introColor, fmt.Sprintf("futures start: %s", futuresStart))
			pretty.ColorPrintln(introColor, fmt.Sprintf("futures end: %s", futuresEnd))
		}

		if futuresSymbols != "" {
			pretty.ColorPrintln(introColor, fmt.Sprintf("futures symbols: %s", strings.Join(symbolsRegex.Split(futuresSymbols, -1), ", ")))
		}
	}

	if coinapi && (cryptoStart != "" && cryptoEnd != "") {
		pretty.ColorPrintln(introColor, fmt.Sprintf("crypto start: %s", cryptoStart))
		pretty.ColorPrintln(introColor, fmt.Sprintf("crypto end: %s", cryptoEnd))
	}

	b := make([]string, 0, 4)

	if futuresDaily {
		b = append(b, "Barchart Futures")
	}

	if futuresIntraday60 {
		b = append(b, "Barchart Futures Intraday 60 Minutes")
	}

	if futuresIntraday30 {
		b = append(b, "Barchart Futures Intraday 30 Minutes")
	}

	if barchart {
		b = append(b, "Barchart")
	}

	if yahoo {
		b = append(b, "Yahoo")
	}

	if investing {
		b = append(b, "Investing.com")
	}

	if coinapi {
		b = append(b, "CoinAPI")
	}

	pretty.ColorPrintln(introColor, fmt.Sprintf("source: %s", strings.Join(b, introSeparator)))

	b = make([]string, 0, 3)

	if download {
		b = append(b, "Download")
	}

	if rename {
		b = append(b, "Rename")
	}

	if check {
		b = append(b, "Check")
	}

	pretty.ColorPrintln(introColor, fmt.Sprintf("operation: %s", strings.Join(b, introSeparator)))

	var istart, iend int
	var err error

	operators := make([]operator.Operator, 0, 4)

	if futuresDaily || futuresIntraday60 || futuresIntraday30 {
		istart, err = strconv.Atoi(futuresStart)
		if err != nil {
			panic(err)
		}

		iend, err = strconv.Atoi(futuresEnd)
		if err != nil {
			panic(err)
		}

		if istart >= iend {
			pretty.ColorPrintln(errorColor, "range start should be smaller than range end")
			return
		}

		symbols := symbolsRegex.Split(futuresSymbols, -1)

		if futuresDaily {
			o := operator.NewBarchartFuturesOperator(istart, iend)
			if futuresSymbols != "" {
				o.SetCustomSymbols(symbols)
			}

			operators = append(operators, o)
		}

		if futuresIntraday60 {
			o := operator.NewBarchartFuturesOperator(istart, iend).IntradaySixtyMinutes()
			if futuresSymbols != "" {
				o.SetCustomSymbols(symbols)
			}

			operators = append(operators, o)
		}

		if futuresIntraday30 {
			o := operator.NewBarchartFuturesOperator(istart, iend).IntradayThirtyMinutes()
			if futuresSymbols != "" {
				o.SetCustomSymbols(symbols)
			}

			operators = append(operators, o)
		}

	}

	if barchart {
		operators = append(
			operators,
			operator.NewBarchartGeneralOperator(),
		)
	}

	if yahoo {
		operators = append(
			operators,
			operator.NewYahooOperator(),
		)
	}

	if investing {
		operators = append(
			operators,
			operator.NewInvestingOperator(),
		)
	}

	if coinapi {
		istart, err = strconv.Atoi(cryptoStart)
		if err != nil {
			panic(err)
		}

		iend, err = strconv.Atoi(cryptoEnd)
		if err != nil {
			panic(err)
		}

		o := operator.NewCoinAPI(istart, iend)
		if cryptoSymbols != "" {
			o.SetCustomSymbols(symbolsRegex.Split(cryptoSymbols, -1))
		}
		operators = append(operators, o)
	}

	for _, op := range operators {
		op.Greeting()

		if download {
			op.Download()
		}

		if rename {
			op.Rename()
		}

		if check {
			op.Check()
		}
	}

}
