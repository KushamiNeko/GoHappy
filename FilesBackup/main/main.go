package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/KushamiNeko/GoFun/utils/pretty"
	"github.com/KushamiNeko/GoHappy/FilesBackup/config"
	"github.com/KushamiNeko/GoHappy/FilesBackup/operation"
)

func main() {

	var from, to, syncFile, safeguard string
	var force, ensure bool

	flag.StringVar(&from, "from", "", "the source directory or file")
	flag.StringVar(&to, "to", "", "destination directory or file")

	flag.StringVar(&syncFile, "syncfile", "", "the file containing sync setups")

	flag.StringVar(&safeguard, "safeguard", "", "the dir path that must exist in dst.")

	flag.BoolVar(&force, "force", false, "force refresh all directories and files")
	flag.BoolVar(&ensure, "ensure", false, "ensure the results")

	flag.Parse()

	pretty.ColorPrintln(config.ColorTitle, "flags:")

	if syncFile != "" {
		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("syncfile: %s", syncFile))
	} else {
		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("from: %s", from))
		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("to: %s", to))
	}

	if safeguard != "" {
		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("safeguard: %s", safeguard))
	} else {
		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("safeguard: %s", config.SafeGuard))
	}

	pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("ensure: %v", ensure))
	pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("force: %v", force))

	pretty.ColorPrintln(config.ColorTitle, "start files backup...")
	start := time.Now()

	var err error

	if safeguard != "" {
		config.SafeGuard = safeguard
	}

	if syncFile == "" {
		if from == "" {
			colorExit(fmt.Errorf("please specify FROM"))
		}

		if to == "" {
			colorExit(fmt.Errorf("please specify TO"))
		}

		src := strings.TrimSpace(from)
		dst := strings.TrimSpace(to)

		err = operation.Sync(src, dst, force, ensure)
		if err != nil {
			colorExit(err)
		}

	} else {

		buffer, err := ioutil.ReadFile(syncFile)
		if err != nil {
			colorExit(err)
		}

		if bytes.Contains(buffer, []byte("\ufeff")) {
			buffer = bytes.Replace(buffer, []byte("\ufeff"), []byte(""), -1)
		}

		content := string(buffer)

		re := regexp.MustCompile(`(.+) -> (.+)`)

		matches := re.FindAllStringSubmatch(content, -1)

		for _, match := range matches {
			src := strings.TrimSpace(match[1])
			dst := strings.TrimSpace(match[2])

			err = operation.Sync(src, dst, force, ensure)
			if err != nil {
				colorExit(err)
			}

		}
	}

	end := time.Now()
	processing := end.Sub(start)

	pretty.ColorPrintln(config.ColorTitle, "files backup completed!!")
	pretty.ColorPrintln(config.ColorTitle, fmt.Sprintf("processing time: %f seconds", processing.Seconds()))
}

func colorExit(err error) {
	pretty.ColorPrintln(config.ColorWarning, err.Error())
	syscall.Exit(1)
}
