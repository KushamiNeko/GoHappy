package process

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var root = filepath.Join(os.Getenv("HOME"), "programming_projects/tools/go/web_factory", "testing")

func TestCacheDiff(t *testing.T) {
	t.SkipNow()

	tf := filepath.Join(root, "caching", "test.txt")

	f, err := os.Create(tf)
	if err != nil {
		t.Errorf("err should be nil but get: %s", err)
	}
	_, err = f.Write([]byte("hello"))
	if err != nil {
		t.Errorf("err should be nil but get: %s", err)
	}
	f.Close()

	diff, err := cacheDiff(tf)
	if diff != true {
		t.Errorf("diff should be true but get: %v", diff)
	}

	diff, err = cacheDiff(tf)
	if diff != false {
		t.Errorf("diff should be false but get: %v", diff)
	}

	f, err = os.OpenFile(tf, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		t.Errorf("err should be nil but get: %s", err)
	}

	_, err = f.Write([]byte(" world"))
	if err != nil {
		t.Errorf("err should be nil but get: %s", err)
	}
	f.Close()

	diff, err = cacheDiff(tf)
	if diff != true {
		t.Errorf("diff should be true but get: %v", diff)
	}
}

func TestIsShadowFile(t *testing.T) {

	tbs := []map[string]string{
		map[string]string{
			"root":     "testing",
			"input":    "testing",
			"expected": "false",
		},
		map[string]string{
			"root":     "./testing",
			"input":    "testing",
			"expected": "false",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/test",
			"expected": "false",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/test.txt",
			"expected": "false",
		},
		map[string]string{
			"root":     "./testing",
			"input":    "testing/test.txt",
			"expected": "false",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/_test",
			"expected": "true",
		},
		map[string]string{
			"root":     "./testing",
			"input":    "testing/_test",
			"expected": "true",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/.test",
			"expected": "true",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/_test.txt",
			"expected": "true",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/.test.txt",
			"expected": "true",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/_test/test.txt",
			"expected": "true",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/.test/test.txt",
			"expected": "true",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/_test/_test.txt",
			"expected": "true",
		},
		map[string]string{
			"root":     "testing",
			"input":    "testing/.test/.test.txt",
			"expected": "true",
		},
	}

	for _, tb := range tbs {
		p := Processor{
			Root: tb["root"],
		}

		shadow := p.isShadowFile(tb["input"])
		if strconv.FormatBool(shadow) != tb["expected"] {
			t.Errorf("expect %s but get %v", tb["expected"], shadow)
		}
	}
}
