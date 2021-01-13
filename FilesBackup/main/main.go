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

	"github.com/KushamiNeko/GoFun/Utils/pretty"
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

	if safeguard != "" {
		config.SafeGuard = safeguard
	}

	pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("safeguard: %s", config.SafeGuard))

	pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("ensure: %v", ensure))
	pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("force: %v", force))

	start := time.Now()

	var err error

	if syncFile == "" {

		if from == "" {
			colorExit(fmt.Errorf("please specify FROM"))
		}

		if to == "" {
			colorExit(fmt.Errorf("please specify TO"))
		}

		src := strings.TrimSpace(from)
		dst := strings.TrimSpace(to)

		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("from: %s", src))
		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("to: %s", dst))
		pretty.ColorPrintln(config.ColorTitle, "start files backup...")

		err = operation.Sync(src, dst, force, ensure)
		if err != nil {
			colorExit(err)
		}

	} else {

		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("syncfile: %s", syncFile))
		pretty.ColorPrintln(config.ColorTitle, "start files backup...")

		buffer, err := ioutil.ReadFile(syncFile)
		if err != nil {
			colorExit(err)
		}

		if bytes.Contains(buffer, []byte("\ufeff")) {
			buffer = bytes.ReplaceAll(buffer, []byte("\ufeff"), []byte(""))
		}

		content := string(buffer)

		re := regexp.MustCompile(`(.+)\s*->\s*(.+)`)

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
