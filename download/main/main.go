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
	"time"

	"github.com/KushamiNeko/go_fun/chart/futures"
	"github.com/KushamiNeko/go_fun/utils/input"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

const (
	interactive = `https://www.barchart.com/futures/quotes/%s/interactive-chart`
	historical  = `https://www.barchart.com/futures/quotes/%s/historical-download`

	historicalPattern  = `^([\w\d]{5})_([^_-]+)(?:-[^_-]+)*_[^_-]+-[^_-]+-\d{2}-\d{2}-\d{4}.csv$`
	interactivePattern = `^([\w\d]{5})_[^_]+_[^_]+_[^_]+_([^_]+)(?:_[^_]+)*_\d{2}_\d{2}_\d{4}.csv$`
)

func command(url string) {
	var outb bytes.Buffer
	var errb bytes.Buffer

	cmd := exec.Command(
		// "google-chrome",
		"firefox",
		url,
	)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Start()
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed400, outb.String())
		pretty.ColorPrintln(pretty.PaperRed400, errb.String())
		panic(err)
	}
}

func barchartPage(root string, symbols []string, ys, ye int, months futures.ContractMonths) {
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

				command(
					fmt.Sprintf(
						root,
						code,
					),
				)
			}

			reader := bufio.NewReader(os.Stdin)
			pretty.ColorPrintln(pretty.PaperAmber300, "press enter when you are ready to proceed")
			reader.ReadString('\n')
		}
	}
}

func validInput(symbols, years, months, page string) error {
	var regex *regexp.Regexp

	if symbols != "" {
		regex = regexp.MustCompile(`^\w+(?:,\w+)*$`)
		if !regex.MatchString(symbols) {
			return fmt.Errorf("invalid symbols: %s", symbols)
		}
	}

	if years != "" {
		regex = regexp.MustCompile(`^(\d{4})(?:(?:\-|\~)(\d{4}))*$`)
		if !regex.MatchString(years) {
			return fmt.Errorf("invalid years: %s", years)
		}
	}

	if months != "" {
		regex = regexp.MustCompile(`^(?:[fghjkmnquvxz]+|all|even|financial)$`)
		if !regex.MatchString(months) {
			return fmt.Errorf("invalid months: %s", months)
		}
	}

	if page != "" {
		regex = regexp.MustCompile(`^(?:historical|interactive)$`)
		if !regex.MatchString(page) {
			return fmt.Errorf("invalid page: %s", page)
		}
	}

	return nil
}

func main() {
	symbols := flag.String("symbols", "", "symbols")
	years := flag.String("years", "", "years")
	months := flag.String("months", "", "months")
	page := flag.String("page", "historical", "barchart page")

	download := flag.Bool("download", false, "download files from barchart")
	front := flag.Bool("front", false, "download front contract")
	rename := flag.Bool("rename", false, "rename downloaded files in Download folder")
	check := flag.Bool("check", false, "check if downloaded files are complete in data source folder")

	flag.Parse()

	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("symbols: %s", *symbols))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("years: %s", *years))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("months: %s", *months))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("page: %s", *page))

	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("download: %v", *download))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("front: %v", *front))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("check: %v", *check))
	pretty.ColorPrintln(pretty.PaperLime400, fmt.Sprintf("rename: %v", *rename))

	err := validInput(*symbols, *years, *months, *page)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed400, err.Error())
		return
	}

	var root string
	switch *page {
	case "historical":
		root = historical
	case "interactive":
		root = interactive
	default:
		panic("unknown page")
	}

	switch {
	case *download:
		if *symbols == "" || *months == "" {
			pretty.ColorPrintln(pretty.PaperRed400, "empty symbols or months")
			return
		}

		ss := strings.Split(*symbols, ",")

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

		if *front {
			for _, symbol := range ss {
				contract := futures.FrontContract(
					time.Now(),
					symbol,
					m,
					futures.BarchartSymbolFormat,
				)

				command(fmt.Sprintf(root, contract))
			}
		} else {
			if *years == "" {
				pretty.ColorPrintln(pretty.PaperRed400, "empty years")
				return
			}

			ys, ye := input.YearsInput(*years)

			barchartPage(root, ss, ys, ye, m)
		}

		if *rename {
			renameDownload()
		}

	case *check:
		if *symbols == "" || *months == "" {
			pretty.ColorPrintln(pretty.PaperRed400, "empty symbols or months")
			return
		}

		if *years == "" {
			pretty.ColorPrintln(pretty.PaperRed400, "empty years")
			return
		}

		ss := strings.Split(*symbols, ",")
		ys, ye := input.YearsInput(*years)

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

		checkDownload(ss, ys, ye, m)
	case *rename:
		renameDownload()
	default:
		panic("unknown case")
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
		//if !rename(root, f.Name(), historicalPattern) {
		//rename(root, f.Name(), interactivePattern)
		//}

		switch {
		case renameBarchart(root, f.Name(), historicalPattern):
		case renameBarchart(root, f.Name(), interactivePattern):
		case renameYahoo(root, f.Name()):
		case renameInvesting(root, f.Name()):
		default:
			continue
		}
	}

}

func renameBarchart(root, file, pattern string) bool {
	dst := filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source/continuous",
	)

	regex := regexp.MustCompile(pattern)
	match := regex.FindAllStringSubmatch(file, -1)
	if len(match) != 0 {

		symbol := strings.ToLower(match[0][1])

		newName := fmt.Sprintf(
			"%s.csv",
			//strings.ToLower(match[0][1]),
			symbol,
		)

		oldPath := filepath.Join(root, file)
		newPath := filepath.Join(dst, symbol[:2], newName)

		//pretty.ColorPrintln(pretty.PaperDeepOrange400, fmt.Sprintf("%s -> %s", file, newName))
		pretty.ColorPrintln(
			pretty.PaperDeepOrange400,
			fmt.Sprintf(
				"%s -> %s",
				//file,
				//newName,
				oldPath,
				newPath,
			),
		)

		err := os.Rename(
			oldPath,
			newPath,
			//filepath.Join(
			//src,
			//file,
			//),
			//filepath.Join(
			//src,
			//newName,
			//),
		)

		if err != nil {
			panic(err)
		}

		return true
	} else {
		return false
	}
}

func renameYahoo(root, file string) bool {
	dst := filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source/yahoo",
	)

	regex := regexp.MustCompile(`\^(\w+)`)
	match := regex.FindAllStringSubmatch(file, -1)
	if len(match) != 0 {

		symbol := strings.ToLower(match[0][1])

		newName := fmt.Sprintf(
			"%s.csv",
			//strings.ToLower(match[0][1]),
			symbol,
		)

		oldPath := filepath.Join(root, file)
		newPath := filepath.Join(dst, newName)

		pretty.ColorPrintln(
			pretty.PaperDeepOrange400,
			fmt.Sprintf(
				"%s -> %s",
				oldPath,
				newPath,
			),
		)

		err := os.Rename(
			oldPath,
			newPath,
		)

		//pretty.ColorPrintln(pretty.PaperDeepOrange400, fmt.Sprintf("%s -> %s", file, newName))

		//err := os.Rename(
		//filepath.Join(
		//src,
		//file,
		//),
		//filepath.Join(
		//src,
		//newName,
		//),
		//)

		if err != nil {
			panic(err)
		}

		return true
	} else {
		return false
	}
}

func renameInvesting(root, file string) bool {
	dst := filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source/investing.com",
	)

	var symbol string

	switch {
	case strings.Contains(file, "STOXX 50 Volatility"):
		symbol = "vstx"
	case strings.Contains(file, "Nikkei Volatility"):
		symbol = "jniv"
	default:
		return false
	}

	newName := fmt.Sprintf(
		"%s.csv",
		symbol,
	)

	oldPath := filepath.Join(root, file)
	newPath := filepath.Join(dst, newName)

	pretty.ColorPrintln(
		pretty.PaperDeepOrange400,
		fmt.Sprintf(
			"%s -> %s",
			oldPath,
			newPath,
		),
	)

	err := os.Rename(
		oldPath,
		newPath,
	)

	//pretty.ColorPrintln(pretty.PaperDeepOrange400, fmt.Sprintf("%s -> %s", file, newName))

	//err := os.Rename(
	//filepath.Join(
	//src,
	//file,
	//),
	//filepath.Join(
	//src,
	//newName,
	//),
	//)
	if err != nil {
		panic(err)
	}

	return true
}
