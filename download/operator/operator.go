package operator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KushamiNeko/go_fun/utils/pretty"
)

type Operator interface {
	Download()
	Rename()
	Check()
}

type operator struct {
	downloadCount int
	renameCount   int
}

func (o *operator) rename(src, dst string) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		panic(err)
	}

	if _, err := os.Stat(filepath.Dir(dst)); os.IsNotExist(err) {
		panic(err)
	}

	o.renameMessage(src, dst)

	//err := os.Rename(
	//srcPath,
	//dstPath,
	//)
	//if err != nil {
	//panic(err)
	//}
}

func (o *operator) downloadCountIncrement() {
	o.downloadCount += 1
}

func (o *operator) renameCountIncrement() {
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

	o.completed()
}

func (o *operator) completed() {
	pretty.ColorPrintln(
		pretty.PaperLime300,
		"press any key to continue",
	)

	fmt.Scanln()
}
