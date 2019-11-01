package shell

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/KushamiNeko/files_backup/config"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

func Diff(src, tar string) (string, error) {

	if src == "" || tar == "" {
		return "", fmt.Errorf("empty src or tar")
	}

	var outb bytes.Buffer
	var errb bytes.Buffer

	cmd := exec.Command("diff", "-rq", src, tar)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		statusRe := regexp.MustCompile(`exit status (\d)`)
		statusM := statusRe.FindAllStringSubmatch(err.Error(), -1)
		if len(statusM) != 1 {
			return "", fmt.Errorf("program error: unknown diff exit status")
		}

		status := strings.TrimSpace(statusM[0][1])

		switch status {

		case "0":
			return strings.TrimSpace(outb.String()), nil

		case "1":
			return strings.TrimSpace(outb.String()), nil

		case "2":
			return "", fmt.Errorf("diff error: %s", errb.String())

		default:
			return "", fmt.Errorf("program error: unknown diff exit status")
		}

	} else {
		return "", nil
	}
}

func Cp(src, tar string) error {

	if src == "" || tar == "" {
		return fmt.Errorf("empty src or tar")
	}

	if !strings.Contains(tar, config.SafeGuard) {
		//panic(fmt.Sprintf("%s is not under safeguard", tar))
		return fmt.Errorf("%s is not under safeguard", tar)
	}

	pretty.ColorPrintln(config.ColorInfo, fmt.Sprintf("cp: %s -> %s", src, tar))

	//var outb bytes.Buffer
	var errb bytes.Buffer

	cmd := exec.Command("cp", "-rp", src, tar)

	//cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func Rm(tar string) error {

	if tar == "" {
		return fmt.Errorf("empty src or tar")
	}

	if !strings.Contains(tar, config.SafeGuard) {
		//panic(fmt.Sprintf("%s is not under safeguard", tar))
		return fmt.Errorf("%s is not under safeguard", tar)
	}

	pretty.ColorPrintln(config.ColorWarning, fmt.Sprintf("rm: %s", tar))

	//var outb bytes.Buffer
	var errb bytes.Buffer

	cmd := exec.Command("rm", "-rf", tar)

	//cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
