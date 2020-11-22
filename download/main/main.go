package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/download/operator"
)

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
	start := flag.String("start", string((time.Now().Year()-1)*100+int(time.Now().Month())), "yyyymm")
	end := flag.String("end", string(time.Now().Year()*100+int(time.Now().Month())), "yyyymm")

	flag.Parse()

	validInput(*start, *end)

	pretty.ColorPrintln(pretty.PaperYellow400, fmt.Sprintf("start: %s", *start))
	pretty.ColorPrintln(pretty.PaperYellow400, fmt.Sprintf("end: %s", *end))

	istart, err := strconv.ParseInt(*start, 10, 64)
	if err != nil {
		panic(err)
	}

	iend, err := strconv.ParseInt(*end, 10, 64)
	if err != nil {
		panic(err)
	}

	operators := []operator.Operator{
		operator.NewBarchartFuturesOperator(int(istart), int(iend)),
		operator.NewBarchartGeneralOperator(),
		operator.NewYahooOperator(),
		operator.NewInvestingOperator(),
	}

	for _, op := range operators {
		op.Download()
		op.Rename()
	}

}
