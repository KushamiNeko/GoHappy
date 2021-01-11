package shell

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KushamiNeko/GoHappy/FilesBackup/config"
)

func TestCpEmpty(t *testing.T) {

	var err error

	err = Cp("", "/home/neko/programming_tools/python_ml")
	if err == nil {
		t.Errorf("err should not be empty due to empty src")
	}

	err = Cp("/home/neko/programming_tools/python_ml", "")
	if err == nil {
		t.Errorf("err should not be empty due to empty tar")
	}

}

func TestRmEmpty(t *testing.T) {

	var err error

	err = Rm("")
	if err == nil {
		t.Errorf("err should not be empty due to empty tar")
	}
}

func TestDiffEmpty(t *testing.T) {

	var err error

	_, err = Diff("", "/home/neko/programming_tools/python_ml")
	if err == nil {
		t.Errorf("err should not be empty due to empty src")
	}

	_, err = Diff("/home/neko/programming_tools/python_ml", "")
	if err == nil {
		t.Errorf("err should not be empty due to empty tar")
	}
}

func TestCpSafeGuard(t *testing.T) {

	var err error

	err = Cp("/home/neko/programming_tools/python_ml", "/home/neko/programming_tools/python_ml")
	if err == nil {
		t.Errorf("err should not be empty due to the safeguard")
	}

}

func TestRmSafeGuard(t *testing.T) {

	var err error

	err = Rm("/home/neko/programming_tools/python_ml")
	if err == nil {
		t.Errorf("err should not be empty due to the safeguard")
	}
}

func TestCpRm(t *testing.T) {

	config.SafeGuard = "/run/media/neko/HDD/TESTING_FIELDS/file_sync"

	var err error

	root := "/run/media/neko/HDD/TESTING_FIELDS/file_sync/TESTING"
	shell := filepath.Join(root, "shell")

	src := filepath.Join(shell, "src")
	dst := filepath.Join(shell, "dst")

	srcf := filepath.Join(shell, "src.txt")
	dstf := filepath.Join(shell, "dst.txt")

	if _, err := os.Stat(src); os.IsNotExist(err) {
		t.Errorf("%s does not exist", src)
	}

	if _, err := os.Stat(srcf); os.IsNotExist(err) {
		t.Errorf("%s does not exist", srcf)
	}

	if _, err := os.Stat(dst); err == nil {
		t.Errorf("%s already exist", dst)
	}

	if _, err := os.Stat(dstf); err == nil {
		t.Errorf("%s already exist", dstf)
	}

	err = Cp(src, dst)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	err = Cp(srcf, dstf)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Errorf("%s should exist", dst)
	}

	if _, err := os.Stat(dstf); os.IsNotExist(err) {
		t.Errorf("%s should exist", dstf)
	}

	diff, err := Diff(src, dst)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
	if diff != "" {
		t.Errorf("%s and %s should be the same", src, dst)
	}

	diff, err = Diff(srcf, dstf)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
	if diff != "" {
		t.Errorf("%s and %s should be the same", srcf, dstf)
	}

	err = Rm(dst)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	err = Rm(dstf)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	if _, err := os.Stat(dst); err == nil {
		t.Errorf("%s should not exist", dst)
	}

	if _, err := os.Stat(dstf); err == nil {
		t.Errorf("%s should not exist", dstf)
	}
}
