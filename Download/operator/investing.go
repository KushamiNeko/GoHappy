package operator

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type investing struct {
	operator
}

func NewInvestingOperator() *investing {
	i := &investing{}
	i.initDir()

	return i
}

func (i *investing) source() map[string]string {
	return map[string]string{
		"https://www.investing.com/indices/stoxx-50-volatility-vstoxx-eur-historical-data": "vstx",
		"https://www.investing.com/indices/jpx-nikkei-400-historical-data":                 "nk400",
		"https://www.investing.com/indices/nikkei-volatility-historical-data":              "jniv",
		"https://www.investing.com/indices/hsi-volatility-historical-data":                 "vhsi",
		"https://www.investing.com/indices/cboe-china-etf-volatility-historical-data":      "vxfxi",
	}
}

func (i *investing) Download() {
	for page, symbol := range i.source() {
		i.download(page, symbol)
	}

	i.downloadCompleted()
}

func (i *investing) Rename() {
	fs, err := ioutil.ReadDir(i.srcDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		for page, symbol := range i.source() {
			_, n := filepath.Split(page)
			n = strings.ToLower(strings.ReplaceAll(n, "-", " "))

			if strings.ToLower(strings.ReplaceAll(f.Name(), "-", " ")) == fmt.Sprintf("%s.csv", n) {
				srcPath := filepath.Join(i.srcDir, f.Name())
				dstPath := filepath.Join(i.dstDir, "investing.com", fmt.Sprintf("%s.csv", symbol))

				i.rename(srcPath, dstPath)
				break
			}
		}
	}

	i.renameCompleted()
}

func (i *investing) Check() {
	for _, symbol := range i.source() {
		path := filepath.Join(i.dstDir, "investing.com", fmt.Sprintf("%s.csv", symbol))
		i.check(path)
	}

	i.checkCompleted()
}
