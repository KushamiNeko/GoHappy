package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/download/operator"
)

const introColor = pretty.PaperYellow300
const introSeparator = ", "

func validInput(start, end string) {

	regex := regexp.MustCompile(`^\d{6}$`)

	if !regex.MatchString(start) {
		panic("invalid start")
	}

	if !regex.MatchString(end) {
		panic("invalid end")
	}

}

func main() {
	var start, end string
	var futures, barchart, yahoo, investing bool
	var download, rename, check bool

	var symbols, months string

	flag.StringVar(&start, "start", fmt.Sprintf("%d", (time.Now().Year()-1)*100+int(time.Now().Month())), "yyyymm")
	flag.StringVar(&end, "end", fmt.Sprintf("%d", (time.Now().Year())*100+int(time.Now().Month())), "yyyymm")

	flag.StringVar(&symbols, "symbols", "", "custom futures symbol list")
	flag.StringVar(&months, "months", "", "custom futures contract months")

	flag.BoolVar(&futures, "futures", false, "download from barchart futures")
	flag.BoolVar(&barchart, "barchart", false, "download from barchart")
	flag.BoolVar(&yahoo, "yahoo", false, "download from yahoo")
	flag.BoolVar(&investing, "investing", false, "download from investing.com")

	flag.BoolVar(&download, "download", false, "download operation")
	flag.BoolVar(&rename, "rename", false, "rename operation")
	flag.BoolVar(&check, "check", false, "check operation")

	flag.Parse()

	validInput(start, end)

	if !futures && !barchart && !yahoo && !investing {
		futures = true
		barchart = true
		yahoo = true
		investing = true
	}

	if !download && !rename && !check {
		download = true
		rename = true
		check = false
	}

	pretty.ColorPrintln(introColor, fmt.Sprintf("start: %s", start))
	pretty.ColorPrintln(introColor, fmt.Sprintf("end: %s", end))

	b := make([]string, 0, 4)

	if futures {
		b = append(b, "futures")
	}

	if barchart {
		b = append(b, "forex")
	}

	if yahoo {
		b = append(b, "yahoo")
	}

	if investing {
		b = append(b, "investing.com")
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

	switch {
	case symbols == "" && months == "":
	case symbols == "" && months != "":
		pretty.ColorPrintln(pretty.PaperRed400, "symbols and months should both be empty or specified")
		return
	case symbols != "" && months == "":
		pretty.ColorPrintln(pretty.PaperRed400, "symbols and months should both be empty or specified")
		return
	case symbols != "" && months != "":
		if !futures {
			pretty.ColorPrintln(pretty.PaperRed400, "symbols and months are used to download futures")
			return
		}

		pretty.ColorPrintln(introColor, fmt.Sprintf("symbols: %s", strings.ReplaceAll(symbols, ",", ", ")))
		pretty.ColorPrintln(introColor, fmt.Sprintf("months: %s", months))
	}

	istart, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		panic(err)
	}

	iend, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		panic(err)
	}

	operators := make([]operator.Operator, 0, 4)

	if futures {
		if symbols != "" && months != "" {
			operators = append(
				operators,
				operator.NewBarchartFuturesOperatorCustom(int(istart), int(iend), symbols, months),
			)
		} else {
			operators = append(
				operators,
				operator.NewBarchartFuturesOperator(int(istart), int(iend)),
			)
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
