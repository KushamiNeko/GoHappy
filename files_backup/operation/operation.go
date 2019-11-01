package operation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/files_backup/config"
	"github.com/KushamiNeko/go_happy/files_backup/shell"
)

func Sync(src, dst string, force, ensure bool) error {
	var err error

	tar, err := tarPath(src, dst)
	if err != nil {
		return err
	}

	if force {
		err = refresh(src, tar, ensure)
		if err != nil {
			return err
		}
	} else {
		err = sync(src, tar, ensure)
		if err != nil {
			return err
		}
	}

	if ensure {
		outb, err := shell.Diff(src, tar)
		if err != nil {
			return err
		}

		pretty.ColorPrintln(config.ColorSpecial, fmt.Sprintf("ensure: %s -> %s", src, tar))

		if outb != "" {
			return fmt.Errorf("failed to ensure %s -> %s", src, tar)
		}
	}

	return nil
}

func sync(src, tar string, ensure bool) error {

	if _, err := os.Stat(tar); os.IsNotExist(err) {
		err = shell.Cp(src, filepath.Dir(tar))
		if err != nil {
			return err
		}
	} else {
		outb, err := shell.Diff(src, tar)
		if err != nil {
			return err
		}

		switch {
		case outb == "":
			return nil

		default:
			for _, f := range differ(outb) {
				f1 := f[0]
				f2 := f[1]

				var s, d string
				switch {
				case strings.Contains(f1, src) && strings.Contains(f2, tar):
					s = f1
					d = f2
				case strings.Contains(f1, tar) && strings.Contains(f2, src):
					s = f2
					d = f1
				default:
					return fmt.Errorf("Unknow file path:\n%s\n%s", f1, f2)
				}

				err := shell.Rm(d)
				if err != nil {
					return err
				}

				err = shell.Cp(s, d)
				if err != nil {
					return err
				}
			}

			for _, f := range only(outb) {
				folder := f[0]
				file := f[1]

				path := filepath.Join(folder, file)

				switch {
				case strings.Contains(path, src):
					s := path
					d := strings.Replace(path, src, tar, -1)
					err = shell.Cp(s, d)
					if err != nil {
						return err
					}

				case strings.Contains(path, tar):
					err := shell.Rm(path)
					if err != nil {
						return err
					}

				default:
					return fmt.Errorf("Unknow file path:\n%s", path)
				}

			}

			return nil

		}
	}

	return nil
}

func refresh(src, tar string, ensure bool) error {
	var err error

	if _, err := os.Stat(tar); err == nil {
		err = shell.Rm(tar)
		if err != nil {
			return err
		}
	}

	err = shell.Cp(src, tar)
	if err != nil {
		return err
	}

	return nil
}

func tarPath(src, dst string) (string, error) {

	if src == dst {
		return "", fmt.Errorf("src path and the dst path are the same")
	}

	srcIn, err := os.Stat(src)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("src path does not exist: %s", src)
	}

	dstIn, err := os.Stat(dst)
	if os.IsNotExist(err) {
		if filepath.Ext(dst) == "" {
			return "", fmt.Errorf("dst path does not exist: %s", src)
		} else {
			if _, err := os.Stat(filepath.Dir(dst)); os.IsNotExist(err) {
				return "", fmt.Errorf("dst path does not exist: %s", src)
			}
		}
	}

	var tar string

	if srcIn.Mode().IsRegular() {
		tar = filepath.Join(dst, filepath.Base(src))
	} else if srcIn.Mode().IsDir() {
		if dstIn.Mode().IsRegular() {
			return "", fmt.Errorf("dst path should not be a file while src path is a dir: %s", dst)
		}

		if filepath.Base(src) == filepath.Base(dst) {
			tar = dst
		} else {
			tar = filepath.Join(dst, filepath.Base(src))
		}
	}

	return tar, nil
}

func differ(outb string) [][]string {

	re := regexp.MustCompile(`Files (.+) and (.+) differ`)
	m := re.FindAllStringSubmatch(outb, -1)

	pair := make([][]string, len(m))

	for i, match := range m {
		f := make([]string, 2)
		f[0] = strings.TrimSpace(match[1])
		f[1] = strings.TrimSpace(match[2])
		pair[i] = f
	}

	return pair
}

func only(outb string) [][]string {

	re := regexp.MustCompile(`Only in (.+): (.+)`)
	m := re.FindAllStringSubmatch(outb, -1)

	pair := make([][]string, len(m))

	for i, match := range m {
		f := make([]string, 2)
		f[0] = strings.TrimSpace(match[1])
		f[1] = strings.TrimSpace(match[2])
		pair[i] = f
	}

	return pair
}
