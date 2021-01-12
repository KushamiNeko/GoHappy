package main

import (
	"flag"
	"fmt"
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

func main() {
	var futuresStart, futuresEnd string
	var cryptoStart, cryptoEnd string

	var futures, barchart, yahoo, investing, coinapi bool
	var download, rename, check bool

	var futuresSymbols string
	var cryptoSymbols string

	flag.StringVar(&futuresStart, "futures-start", fmt.Sprintf("%d", (time.Now().Year()-1)*100+int(time.Now().Month())), "yyyymm")
	flag.StringVar(&futuresEnd, "futures-end", fmt.Sprintf("%d", (time.Now().Year())*100+int(time.Now().Month())), "yyyymm")

	flag.StringVar(&futuresSymbols, "futures-symbols", "", "custom futures symbol list")

	flag.StringVar(&cryptoStart, "crypto-start", fmt.Sprintf("%d", (time.Now().Year()-1)*100+int(time.Now().Month())), "yyyymm")
	flag.StringVar(&cryptoEnd, "crypto-end", fmt.Sprintf("%d", (time.Now().Year())*100+int(time.Now().Month())), "yyyymm")

	flag.StringVar(&cryptoSymbols, "crypto-symbols", "", "custom crypto symbol list")

	flag.BoolVar(&futures, "futures", false, "download futures data from Barchart")
	flag.BoolVar(&barchart, "barchart", false, "download data from Barchart")
	flag.BoolVar(&yahoo, "yahoo", false, "download from Yahoo")
	flag.BoolVar(&investing, "investing", false, "download from Investing.com")
	flag.BoolVar(&coinapi, "coinapi", false, "download from CoinAPI")

	flag.BoolVar(&download, "download", false, "download operation")
	flag.BoolVar(&rename, "rename", false, "rename operation")
	flag.BoolVar(&check, "check", false, "check operation")

	flag.Parse()

	const pattern = `^\d{6}$`

	if futures && !input.ValidateWithRegex(futuresStart, pattern) {
		pretty.ColorPrintln(errorColor, "invalid futures start")
		return
	}

	if futures && !input.ValidateWithRegex(futuresEnd, pattern) {
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

	if !futures && !barchart && !yahoo && !investing && !coinapi {
		futures = true
		barchart = true
		yahoo = true
		investing = true
		coinapi = true
	}

	if !download && !rename && !check {
		download = true
		rename = true
		check = false
	}

	if futures {
		if futuresStart != "" && futuresEnd != "" {
			pretty.ColorPrintln(introColor, fmt.Sprintf("futures start: %s", futuresStart))
			pretty.ColorPrintln(introColor, fmt.Sprintf("futures end: %s", futuresEnd))
		}

		if futuresSymbols != "" {
			pretty.ColorPrintln(introColor, fmt.Sprintf("futures symbols: %s", strings.ReplaceAll(futuresSymbols, ",", ", ")))
		}
	}

	if coinapi && (cryptoStart != "" && cryptoEnd != "") {
		pretty.ColorPrintln(introColor, fmt.Sprintf("crypto start: %s", cryptoStart))
		pretty.ColorPrintln(introColor, fmt.Sprintf("crypto end: %s", cryptoEnd))
	}

	b := make([]string, 0, 4)

	if futures {
		b = append(b, "Barchart Futures")
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
		b = append(b, "download")
	}

	if rename {
		b = append(b, "rename")
	}

	if check {
		b = append(b, "check")
	}

	pretty.ColorPrintln(introColor, fmt.Sprintf("operation: %s", strings.Join(b, introSeparator)))

	var istart, iend int
	var err error

	operators := make([]operator.Operator, 0, 4)

	if futures {
		istart, err = strconv.Atoi(futuresStart)
		if err != nil {
			panic(err)
		}

		iend, err = strconv.Atoi(futuresEnd)
		if err != nil {
			panic(err)
		}

		o := operator.NewBarchartFuturesOperator(istart, iend)
		if futuresSymbols != "" {
			o.SetCustomSymbols(strings.Split(futuresSymbols, ","))
		}
		operators = append(operators, o)
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
			o.SetCustomSymbols(strings.Split(cryptoSymbols, ","))
		}
		operators = append(operators, o)
	}

	for _, op := range operators {
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
