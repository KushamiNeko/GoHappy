package command

import (
	"bytes"
	"os/exec"

	"github.com/KushamiNeko/go_fun/utils/pretty"
)

const (
	Google  = "google-chrome"
	Firefox = "firefox"
)

func Download(url string) {
	var outb bytes.Buffer
	var errb bytes.Buffer

	cmd := exec.Command(
		// "google-chrome",
		//"firefox",
		Firefox,
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
