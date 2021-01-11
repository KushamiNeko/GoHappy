package operation

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KushamiNeko/GoHappy/FilesBackup/shell"
)

const (
	root string = "/run/media/neko/HDD/TESTING_FIELDS/file_sync/TESTING"
)

var (
	src string = filepath.Join(root, "src")
	dst string = filepath.Join(root, "dst")
	tar string = filepath.Join(dst, "src")
)

func TestOnly(t *testing.T) {
	outb, err := shell.Diff(src, tar)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	err = onlyHelper(outb)
	if err != nil {
		t.Error(err)
	}

	outb, err = shell.Diff(tar, src)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	err = onlyHelper(outb)
	if err != nil {
		t.Error(err)
	}
}

func onlyHelper(outb string) error {
	for _, f := range only(outb) {
		fd := f[0]
		fs := f[1]

		if strings.Contains(fs, "only_src") {
			if !strings.Contains(fd, src) {
				return fmt.Errorf("%s should contains %s", fd, src)
			}
		} else if strings.Contains(fs, "only_dst") {
			if !strings.Contains(fd, dst) {
				return fmt.Errorf("%s should contains %s", fd, dst)
			}
		} else {
			return fmt.Errorf("unknown file: %s/%s", fd, fs)
		}
	}

	return nil
}

func TestDifferFolder(t *testing.T) {
	outb, err := shell.Diff(src, tar)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	for _, f := range differ(outb) {
		f1 := f[0]
		f2 := f[1]

		if !strings.Contains(f1, "differ") {
			t.Errorf("%s should contains differ", f1)
		}

		if !strings.Contains(f2, "differ") {
			t.Errorf("%s should contains differ", f2)
		}

		if !strings.Contains(f1, src) && !strings.Contains(f1, dst) {
			t.Errorf("%s should contains %s or %s", f1, src, dst)
		}

		if !strings.Contains(f2, src) && !strings.Contains(f2, dst) {
			t.Errorf("%s should contains %s or %s", f2, src, dst)
		}
	}
}

func filesDifferHelper(outb string) error {
	for _, f := range differ(outb) {
		f1 := f[0]
		f2 := f[1]

		if strings.Contains(f1, "src") {
			if !strings.Contains(f2, "differ") {
				return fmt.Errorf("%s should contains differ", f2)
			}
		} else if strings.Contains(f2, "src") {
			if !strings.Contains(f1, "differ") {
				return fmt.Errorf("%s should contains differ", f1)
			}
		} else {
			return fmt.Errorf("unknown file: %s, %s", f1, f2)
		}

	}

	return nil
}

func TestDifferFile(t *testing.T) {
	srcf := filepath.Join(root, "src.txt")
	tarf := filepath.Join(root, "dst.txt")

	outb, err := shell.Diff(srcf, tarf)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	if outb != "" {
		t.Errorf("outb should be empty but get %s", outb)
	}

	srcf = filepath.Join(root, "src.txt")
	tarf = filepath.Join(root, "differ.txt")

	outb, err = shell.Diff(srcf, tarf)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	err = filesDifferHelper(outb)
	if err != nil {
		t.Error(err)
	}

	srcf = filepath.Join(root, "differ.txt")
	tarf = filepath.Join(root, "src.txt")

	outb, err = shell.Diff(srcf, tarf)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	err = filesDifferHelper(outb)
	if err != nil {
		t.Error(err)
	}

}

func TestTarPathSame(t *testing.T) {
	_, err := tarPath(src, src)
	if err == nil {
		t.Errorf("err should not be empty")
	}
}

func TestTarPathSrcNonexist(t *testing.T) {
	var err error

	srcf := filepath.Join(root, "nonexist")
	tarf := filepath.Join(root, "dst")

	_, err = tarPath(srcf, tarf)
	if err == nil {
		t.Errorf("err should not be empty")
	}

	srcf = filepath.Join(root, "nonexist", "src.txt")
	tarf = filepath.Join(root, "dst")

	_, err = tarPath(srcf, tarf)
	if err == nil {
		t.Errorf("err should not be empty")
	}

	srcf = filepath.Join(root, "nonexist.txt")
	tarf = filepath.Join(root, "dst")

	_, err = tarPath(srcf, tarf)
	if err == nil {
		t.Errorf("err should not be empty")
	}
}

func TestTarPathFiles(t *testing.T) {
	var err error

	srcf := filepath.Join(root, "src.txt")
	tarf := filepath.Join(root, "dst.txt")

	tar, err := tarPath(srcf, tarf)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	if tar != tarf {
		t.Errorf("tar should be %s but get %s", tarf, tar)
	}

	srcf = filepath.Join(root, "src.txt")
	tarf = filepath.Join(root, "nonexist.txt")

	tar, err = tarPath(srcf, tarf)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	if tar != tarf {
		t.Errorf("tar should be %s but get %s", tarf, tar)
	}

	srcf = filepath.Join(root, "src.txt")
	tarf = filepath.Join(root, "nonexist", "nonexist.txt")

	tar, err = tarPath(srcf, tarf)
	if err == nil {
		t.Errorf("err should not be empty")
	}
}

func TestTarPathOneFileOneDir(t *testing.T) {
	var err error

	srcf := filepath.Join(root, "src.txt")
	tarf := filepath.Join(root, "nonexist")

	tar, err = tarPath(srcf, tarf)
	if err == nil {
		t.Errorf("err should not be empty")
	}

	srcf = filepath.Join(root, "src")
	tarf = filepath.Join(root, "dst.txt")

	tar, err = tarPath(srcf, tarf)
	if err == nil {
		t.Errorf("err should not be empty")
	}
}

func TestTarPathDirs(t *testing.T) {
	var err error

	srcf := filepath.Join(root, "src")
	tarf := filepath.Join(root, "dst")

	tar, err = tarPath(srcf, tarf)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	if tar != filepath.Join(tarf, "src") {
		t.Errorf("tar should be %s but get %s", filepath.Join(tarf, "src"), tar)
	}

	srcf = filepath.Join(root, "src")
	tarf = filepath.Join(root, "dst", "src")

	tar, err = tarPath(srcf, tarf)
	if err != nil {
		t.Errorf("err should be empty but get %s", err.Error())
	}

	if tar != tarf {
		t.Errorf("tar should be %s but get %s", tarf, tar)
	}

	srcf = filepath.Join(root, "src")
	tarf = filepath.Join(root, "nonexist")

	tar, err = tarPath(srcf, tarf)
	if err == nil {
		t.Errorf("err should not be empty")
	}

}
