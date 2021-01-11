package operator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KushamiNeko/GoFun/utils/pretty"
	"github.com/KushamiNeko/GoHappy/Download/command"
)

const processLimit = 10

type Operator interface {
	Download()
	Rename()
	Check()
}

type operator struct {
	srcDir string
	dstDir string

	downloadCount int
	renameCount   int

	missingCount int

	processCount int
}

func (o *operator) initDir() {
	o.srcDir = filepath.Join(
		os.Getenv("HOME"),
		"Downloads",
	)

	if _, err := os.Stat(o.srcDir); os.IsNotExist(err) {
		panic("src dir does not exist")
	}

	o.dstDir = filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source",
	)

	if _, err := os.Stat(o.dstDir); os.IsNotExist(err) {
		panic("dst dir does not exist")
	}
}

func (o *operator) setDir(src, dst string) {
	o.srcDir = src
	o.dstDir = dst
}

func (o *operator) download(page, message string) {
	o.downloadMessage(message)

	command.Download(page)

	o.downloadCount++

	o.checkProcessLimit()
}

func (o *operator) rename(src, dst string) {
	if _, err := os.Stat(o.srcDir); os.IsNotExist(err) {
		panic(err)
	}

	if _, err := os.Stat(o.dstDir); os.IsNotExist(err) {
		panic(err)
	}

	if _, err := os.Stat(src); os.IsNotExist(err) {
		panic(err)
	}

	if _, err := os.Stat(filepath.Dir(dst)); os.IsNotExist(err) {
		panic(err)
	}

	o.renameMessage(src, dst)

	err := os.Rename(
		src,
		dst,
	)
	if err != nil {
		panic(err)
	}

	o.renameCount++
}

func (o *operator) check(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		o.checkMessage(strings.ReplaceAll(path, fmt.Sprintf(".%s", filepath.Ext(path)), ""))
		o.missingCount++
	}
}

func (o *operator) downloadMessage(symbol string) {
	pretty.ColorPrintln(
		pretty.PaperCyan300,
		fmt.Sprintf("downloading: %s", symbol),
	)
}

func (o *operator) renameMessage(src, dst string) {
	pretty.ColorPrintln(
		pretty.PaperOrange300,
		fmt.Sprintf(
			"%s -> %s",
			src,
			dst,
		),
	)
}

func (o *operator) checkMessage(symbol string) {
	pretty.ColorPrintln(
		pretty.PaperIndigo300,
		fmt.Sprintf("missing: %s", symbol),
	)
}

func (o *operator) downloadCompleted() {
	pretty.ColorPrintln(
		pretty.PaperLightGreenA200,
		fmt.Sprintf("%d files downloaded", o.downloadCount),
	)

	pretty.ColorPrintln(
		pretty.PaperGreen400,
		"download completed",
	)

	o.completed()
}

func (o *operator) renameCompleted() {
	pretty.ColorPrintln(
		pretty.PaperLightGreenA200,
		fmt.Sprintf("rename %d files", o.renameCount),
	)

	if o.downloadCount != o.renameCount {
		pretty.ColorPrintln(
			pretty.PaperRed600,
			"rename operation miss some downloaded files",
		)
	}

	pretty.ColorPrintln(
		pretty.PaperGreen400,
		"rename completed",
	)

	o.completed()
}

func (o *operator) checkCompleted() {
	pretty.ColorPrintln(
		pretty.PaperLightGreenA200,
		fmt.Sprintf("check missing %d files", o.missingCount),
	)

	pretty.ColorPrintln(
		pretty.PaperGreen400,
		"check completed",
	)

	o.completed()
}

func (o *operator) completed() {
	pretty.ColorPrintln(
		pretty.PaperBrown300,
		"press any key to continue",
	)

	fmt.Scanln()
	o.processCount = 0
}

func (o *operator) checkProcessLimit() {
	o.processCount++
	if o.processCount >= processLimit {
		o.completed()
		o.processCount = 0
	}
}
