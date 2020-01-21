package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KushamiNeko/go_fun/chart/futures"
	"github.com/KushamiNeko/go_fun/utils/input"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

type Page int

const (
	interactive = `https://www.barchart.com/futures/quotes/%s/interactive-chart`
	historical  = `https://www.barchart.com/futures/quotes/%s/historical-download`

	historicalPage Page = iota
	interactivePage

	historicalPattern  = `([\w\d]{5})_(\w+)_\w+-\w+-\d{2}-\d{2}-\d{4}.csv`
	interactivePattern = `([\w\d]{5})_\w+_\w+_\w+_(\w+)_\d{2}_\d{2}_\d{4}.csv`
)

func barchartPage(symbols []string, ys, ye int, months futures.ContractMonths, page Page) {
	var root string
	switch page {
	case historicalPage:
		root = historical
	case interactivePage:
		root = interactive
	default:
		panic("unknown page")
	}

	for _, symbol := range symbols {
		for y := ys; y < ye; y++ {
			for _, month := range months {

				code := fmt.Sprintf(
					"%s%s%02d",
					symbol,
					string(month),
					y%100,
				)

				pretty.ColorPrintln(pretty.PaperCyan300, code)

				var outb bytes.Buffer
				var errb bytes.Buffer

				cmd := exec.Command(
					"google-chrome",
					fmt.Sprintf(
						root,
						code,
					),
				)

				cmd.Stdout = &outb
				cmd.Stderr = &errb

				err := cmd.Start()
				if err != nil {
					panic(err)
				}

			}

			reader := bufio.NewReader(os.Stdin)
			pretty.ColorPrintln(pretty.PaperAmber300, "press enter when you are ready to proceed")
			reader.ReadString('\n')
		}
	}
}

func validInput(symbols, years, months, page string) error {
	var regex *regexp.Regexp

	regex = regexp.MustCompile(`^\w+(?:,\w+)*$`)
	if !regex.MatchString(symbols) {
		return fmt.Errorf("invalid symbols: %s", symbols)
	}

	regex = regexp.MustCompile(`^(\d{4})(?:(?:\-|\~)(\d{4}))*$`)
	if !regex.MatchString(years) {
		return fmt.Errorf("invalid years: %s", years)
	}

	regex = regexp.MustCompile(`^(?:[fghjkmnquvxz]+|all|even|financial)$`)
	if !regex.MatchString(months) {
		return fmt.Errorf("invalid months: %s", months)
	}

	regex = regexp.MustCompile(`^(?:historical|interactive)$`)
	if !regex.MatchString(page) {
		return fmt.Errorf("invalid page: %s", page)
	}

	return nil
}

func main() {
	symbols := flag.String("symbols", "", "symbols")
	years := flag.String("years", "", "years")
	months := flag.String("months", "", "months")
	page := flag.String("page", "historical", "barchart page")

	rename := flag.Bool("rename", false, "rename downloaded files in Download folder")
	check := flag.Bool("check", false, "check if downloaded files are complete in data source folder")

	flag.Parse()

	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("symbols: %s", *symbols))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("years: %s", *years))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("months: %s", *months))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("page: %s", *page))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("rename: %v", *rename))

	if *symbols != "" && *years != "" && *months != "" {
		err := validInput(*symbols, *years, *months, *page)
		if err != nil {
			pretty.ColorPrintln(pretty.PaperRed400, err.Error())
			return
		}

		ss := strings.Split(*symbols, ",")
		ys, ye := input.YearsInput(*years)

		var p Page
		switch *page {
		case "historical":
			p = historicalPage
		case "interactive":
			p = interactivePage
		default:
			panic("unknown page")
		}

		var m futures.ContractMonths
		switch *months {
		case "all":
			m = futures.AllContractMonths
		case "even":
			m = futures.EvenContractMonths
		case "financial":
			m = futures.FinancialContractMonths
		default:
		}

		if *check {
			checkDownload(ss, ys, ye, m)
		} else {
			barchartPage(
				ss,
				ys,
				ye,
				m,
				p,
			)
		}
	}

	if *rename {
		renameDownload()
	}

}

func checkDownload(symbols []string, ys, ye int, months futures.ContractMonths) {
	root := filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source/continuous",
	)

	for _, symbol := range symbols {
		for y := ys; y < ye; y++ {
			for _, month := range months {

				file := fmt.Sprintf("%s%s%02d.csv", symbol, string(month), y%100)

				if _, err := os.Stat(filepath.Join(root, symbol, file)); os.IsNotExist(err) {
					pretty.ColorPrintln(pretty.PaperRed400, fmt.Sprintf("file does not exist: %s", file))
				}
			}
		}
	}
}

func renameDownload() {
	root := filepath.Join(
		os.Getenv("HOME"),
		"Downloads",
	)

	fs, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		if !rename(root, f.Name(), historicalPattern) {
			rename(root, f.Name(), interactivePattern)
		}
	}

}

func rename(root, file, pattern string) bool {
	regex := regexp.MustCompile(pattern)
	match := regex.FindAllStringSubmatch(file, -1)
	if len(match) != 0 {

		newName := fmt.Sprintf(
			"%s.csv",
			strings.ToLower(match[0][1]),
		)

		pretty.ColorPrintln(pretty.PaperDeepOrange400, fmt.Sprintf("%s -> %s", file, newName))

		err := os.Rename(
			filepath.Join(
				root,
				file,
			),
			filepath.Join(
				root,
				newName,
			),
		)

		if err != nil {
			panic(err)
		}

		return true
	} else {
		return false
	}
}
