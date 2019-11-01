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

type DartOperator struct {
	Optimized bool
}

func (d *DartOperator) Operate(src, dst string) (string, error) {
	ext := filepath.Ext(src)

	if ext != ".dart" {
		return "", fmt.Errorf("unknown file extension: %s", ext)
	}

	var err error

	stat, err := os.Stat(dst)
	if err != nil {
		return "", err
	}

	if stat.Mode().IsDir() {
		dst = filepath.Join(dst, strings.ReplaceAll(filepath.Base(src), filepath.Ext(src), ".js"))
	} else {
		return "", fmt.Errorf("dst should be a directory")
	}

	pretty.ColorPrintln(pretty.PaperBlue300, fmt.Sprintf("processing file: %s -> %s", src, dst))

	var cmd *exec.Cmd

	if d.Optimized {
		cmd = exec.Command("dart2js", "-O2", src, "-o", dst)
	} else {
		cmd = exec.Command("dart2js", src, "-o", dst, "--enable-asserts")
	}

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
