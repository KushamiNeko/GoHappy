package process

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/KushamiNeko/go_fun/utils/pretty"
)

type CssOperator struct{}

//func (c *CssOperator) DstPath(src, dst string) string {
//return filepath.Join(dst, strings.ReplaceAll(filepath.Base(src), filepath.Ext(src), ".css"))
//}

func (c *CssOperator) Operate(src, dst string) (string, error) {
	ext := filepath.Ext(src)

	if ext != ".scss" && ext != "sass" {
		return "", fmt.Errorf("unknown file extension: %s", ext)
	}

	var err error

	stat, err := os.Stat(dst)
	if err != nil {
		return "", err
	}

	if stat.Mode().IsDir() {
		dst = filepath.Join(dst, strings.ReplaceAll(filepath.Base(src), filepath.Ext(src), ".css"))
	} else {
		return "", fmt.Errorf("dst should be a directory")
	}

	pretty.ColorPrintln(pretty.PaperBlue300, fmt.Sprintf("processing file: %s -> %s", src, dst))

	cmd := exec.Command("sass", "--no-source-map", src, dst)

	var cmdErr bytes.Buffer
	var cmdOut bytes.Buffer

	cmd.Stderr = &cmdErr
	cmd.Stdout = &cmdOut

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s:\n%s\n%s\n", err.Error(), cmdErr.String(), cmdOut.String())
	}

	return dst, nil
}
