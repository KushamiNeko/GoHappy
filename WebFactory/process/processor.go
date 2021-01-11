package process

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/KushamiNeko/go_fun/utils/pretty"
)

type Operator interface {
	Operate(src, dst string) (string, error)
}

var (
	caches map[string][]byte
)

func init() {
	caches = make(map[string][]byte)
}

func hashFile(src string) ([]byte, error) {
	h := sha512.New()

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	h.Write(content)
	hash := h.Sum(nil)
	return hash, nil
}

func cacheFile(file string) error {
	hash, err := hashFile(file)
	if err != nil {
		return err
	}

	pretty.ColorPrintln(pretty.PaperAmber400, fmt.Sprintf("caching file: %s", file))
	caches[file] = hash
	return nil
}

func cacheDiff(file string) (bool, error) {
	hash, err := hashFile(file)
	if err != nil {
		return true, err
	}

	v, ok := caches[file]
	if !ok || !bytes.Equal(hash, v) {
		return true, nil
	}

	return false, nil
}

type Processor struct {
	Root       string
	Dst        string
	Operations []string
	Templated  bool
	Optimized  bool

	Initialized bool
}

func (p *Processor) Caching(src string) error {
	var stat os.FileInfo
	var err error

	stat, err = os.Stat(src)
	if err != nil {
		return err
	}

	switch {
	case stat.Mode().IsRegular():
		tar := false
		for _, op := range p.Operations {
			if strings.HasSuffix(filepath.Ext(src), strings.TrimSpace(op)) {
				tar = true
			}
		}

		if tar == false {
			return nil
		}

		caches[src] = []byte{}

	case stat.Mode().IsDir():
		dir, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}

		for _, f := range dir {
			srcr := filepath.Join(src, f.Name())
			err = p.Caching(srcr)
			if err != nil {
				return err
			}
		}

	default:
		panic(fmt.Sprintf("unknown file mode: %s", stat.Mode().String()))
	}

	return nil
}
func (p *Processor) isShadowFile(path string) bool {
	for {
		if path == p.Root || path == "" {
			break
		}
		base := filepath.Base(path)
		if strings.HasPrefix(base, "_") || strings.HasPrefix(base, ".") {
			return true
		} else {
			path = filepath.Dir(path)
		}
	}

	return false
}

func (p *Processor) Process() error {
	shadowChanged := false

	nfs := make([]string, 0, len(caches))
	cfs := make([]string, 0, len(caches))

	sc := make(map[string]bool)

	for k, _ := range caches {
		stat, err := os.Stat(k)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				return err
			}
		}

		if !stat.Mode().IsRegular() {
			panic(fmt.Sprintf("cache should only apply to files: %s", k))
		}

		diff, err := cacheDiff(k)
		if err != nil {
			return err
		}

		if p.isShadowFile(k) {
			if diff {
				shadowChanged = true

				ext := filepath.Ext(k)
				if _, ok := sc[ext]; !ok {
					sc[ext] = true
				}

				err = cacheFile(k)
				if err != nil {
					return err
				}
			}
		} else {
			nfs = append(nfs, k)

			if diff {
				cfs = append(cfs, k)
			}
		}
	}

	if !p.Initialized {
		for _, f := range nfs {
			p.operate(f)
		}
	} else if shadowChanged {
		for ext, _ := range sc {
			for _, f := range nfs {
				if filepath.Ext(f) == ext {
					p.operate(f)
				}
			}
		}
	} else {
		for _, c := range cfs {
			p.operate(c)
		}
	}

	if p.Initialized == false {
		p.Initialized = true
	}

	return nil
}

func (p *Processor) operator(ext string) Operator {
	switch ext {
	case ".scss":
		return &CSSOperator{}
	case ".dart":
		return &DartOperator{p.Optimized}
	case ".ts":
		panic("not implement")
	default:
		panic(fmt.Sprintf("unknown extension: %s", ext))
	}
}

func (p *Processor) operate(src string) error {

	var err error
	operator := p.operator(filepath.Ext(src))

	dst := filepath.Join(
		p.Dst,
		strings.ReplaceAll(
			filepath.Dir(src),
			p.Root,
			"",
		),
	)

	out, err := operator.Operate(src, dst)
	if err != nil {
		return err
	}

	if p.Templated {
		to := &TemplateOperator{
			CleanSrc: true,
		}
		err = to.Operate(out)
		if err != nil {
			return err
		}
	}

	err = cacheFile(src)
	if err != nil {
		return err
	}

	return nil
}
