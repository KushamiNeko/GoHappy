package operator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/download/command"
)

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
}

func (o *operator) initDir() {
	o.srcDir = filepath.Join(
		os.Getenv("HOME"),
		"Downloads",
	)

	o.dstDir = filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source",
	)
}

func (o *operator) setDir(src, dst string) {
	o.srcDir = src
	o.dstDir = dst
}

func (o *operator) download(page, message string) {
	o.downloadMessage(message)

	command.Download(page)

	o.downloadCount += 1
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

	o.renameCount += 1
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

	o.completed()
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

func (o *operator) completed() {
	pretty.ColorPrintln(
		pretty.PaperBrown300,
		"press any key to continue",
	)

	fmt.Scanln()
}
