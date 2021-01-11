package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/web_factory/process"
)

func main() {

	input := flag.String("input", "", "the input file or directory")
	output := flag.String("output", "", "the output directory")

	operations := flag.String("operations", "scss,dart", "scss,dart,ts")

	interval := flag.Int("interval", 5, "sleep interval between each loop")

	templated := flag.Bool("templated", false, "generate go templates")
	optimized := flag.Bool("optimized", false, "optimize the output when available")

	flag.Parse()

	err := validateInputs(*input, *output, *operations)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed400, err.Error())
		return
	}

	src := filepath.Clean(*input)
	dst := filepath.Clean(getDst(*input, *output))

	pretty.ColorPrintln(pretty.PaperGreen300, fmt.Sprintf("src: %s", src))
	pretty.ColorPrintln(pretty.PaperGreen300, fmt.Sprintf("dst: %s", dst))
	pretty.ColorPrintln(pretty.PaperGreen300, fmt.Sprintf("operations: %s", *operations))
	pretty.ColorPrintln(pretty.PaperGreen300, fmt.Sprintf("templated: %v", *templated))
	pretty.ColorPrintln(pretty.PaperGreen300, fmt.Sprintf("optimized: %v", *optimized))
	pretty.ColorPrintln(pretty.PaperGreen300, fmt.Sprintf("interval: %d", *interval))

	ops := strings.Split(*operations, ",")

	p := &process.Processor{
		Root:       src,
		Dst:        dst,
		Operations: ops,
		Templated:  *templated,
		Optimized:  *optimized,
	}

	err = p.Caching(src)
	if err != nil {
		pretty.ColorPrintln(pretty.PaperRed400, err.Error())
		return
	}

	for {
		err = p.Process()
		if err != nil {
			pretty.ColorPrintln(pretty.PaperRed400, err.Error())
			return
		}

		time.Sleep(time.Duration(*interval))
	}
}

func validateInputs(input, output, operations string) error {
	if input == "" {
		return fmt.Errorf("please specify the input")
	} else {
		stat, err := os.Stat(input)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("input does not exist: %s", input)
			} else {
				return err
			}
		}

		if !stat.Mode().IsDir() {
			return fmt.Errorf("input should be a path of directory")
		}
	}

	if filepath.Ext(output) != "" {
		return fmt.Errorf("output should be a path of directory")
	}

	if operations == "" {
		return fmt.Errorf("please specify the operations")
	} else {
		regex := regexp.MustCompile(`^(?:scss|dart|ts)(?:,(?:scss|dart|ts))*$`)
		if !regex.MatchString(operations) {
			return fmt.Errorf("invalid operation")
		}
	}

	return nil
}

func getDst(input, output string) string {
	var stat os.FileInfo
	var err error

	if output == "" {
		stat, err = os.Stat(input)
		if err != nil {
			panic(err)
		}

		if stat.Mode().IsDir() {
			return input
		}

		if stat.Mode().IsRegular() {
			return filepath.Dir(input)
		}
	}

	return output
}
