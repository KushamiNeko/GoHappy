package operator

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type yahoo struct {
	operator
}

func NewYahooOperator() *yahoo {
	y := &yahoo{}
	y.initDir()

	return y
}

func (y *yahoo) source() map[string]string {
	return map[string]string{
		"btc-usd": "20141001",
		"eth-usd": "20150901",
		"ltc-usd": "20141001",
		"xrp-usd": "20141001",
		//
		"^vix": "19900101",
		"^vxn": "20000101",
		//"^sml":  "19890101",
		"^rut":  "19880101",
		"^dji":  "19860101",
		"^n225": "19650101",
		"^gspc": "19270101",
		"^ixic": "19710101",
		"^ndx":  "19850101",
		//"^nya":  "19650101",
		"^hsi": "19860101",
		//
		"ezu": "20000101",
		"eem": "20030101",
		"fxi": "20040101",
		// "^ovx": "20070101",
		// "^gvz": "20100101",
		// ##########
		"near": "20130101",
		"jpst": "20170701",
		"icsh": "20140101",
		"gsy":  "20080301",
		// "mbb": "20070401",
		// "flot": "20110701",
		"shv": "20070101",
		// "shyg": "20140101",
		"hyg":  "20070101",
		"jnk":  "20080101",
		"ushy": "20180101",
		// "faln": "20160701",
		"emb": "20070101",
		// "emhy": "20120501",
		// "slqd": "20140101",
		"lqd": "20020101",
		// "usig": "20070201",
		// "igsb": "20070201",
		// "igib": "20070201",
		// "iglb": "20100101",
		// "qlta": "20120301",
		// "lqdh": "20140701",
		"shy": "20020801",
		// "iei": "20070201",
		"ief": "20020801",
		// "tlh": "20070201",
		// "tlt": "20020801",
		// "govt": "20120301",
		// "igov": "20090201",
		// "stip": "20110101",
		// "tip": "20040101",
		// "sub": "20090101",
		// "mub": "20071001",
		"iyr":  "20000101",
		"rem":  "20070101",
		"reet": "20140801",
		// "icvt": "20150701",
		// "istb": "20130101",
		// "iusb": "20140701",
		// "agg": "20040101",
		// "byld": "20140501",

	}
}

func (y *yahoo) url(symbol, datetime string) string {

	dt, err := time.Parse("20060102", datetime)
	if err != nil {
		panic(err)
	}

	var b strings.Builder

	b.WriteString(
		fmt.Sprintf(
			"https://finance.yahoo.com/quote/%s/history?",
			url.PathEscape(symbol),
		),
	)

	b.WriteString(
		fmt.Sprintf(
			"period1=%d&",
			dt.Unix(),
		),
	)

	b.WriteString(
		fmt.Sprintf(
			"period2=%d&",
			time.Now().Add(5*24*time.Hour).Unix(),
		),
	)

	b.WriteString("interval=1d&filter=history&frequency=1d")

	return b.String()
}

func (y *yahoo) Download() {

	for symbol, datetime := range y.source() {
		y.download(y.url(symbol, datetime), symbol)
	}

	y.downloadCompleted()
}

func (y *yahoo) renameSymbol(symbol string) string {
	switch symbol {
	case "n225":
		return "nikk"
	case "gspc":
		return "spx"
	case "ixic":
		return "compq"
	default:
		return symbol
	}
}

func (y *yahoo) Rename() {
	indexRegex := regexp.MustCompile(`(\^(\w+)).csv`)

	fs, err := ioutil.ReadDir(y.srcDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		symbol := strings.ToLower(strings.ReplaceAll(f.Name(), filepath.Ext(f.Name()), ""))

		if _, ok := y.source()[symbol]; ok {
			match := indexRegex.FindAllStringSubmatch(f.Name(), -1)
			if len(match) != 0 {
				symbol = strings.ToLower(match[0][2])
			}

			symbol = strings.ReplaceAll(symbol, "-", "")

			symbol = y.renameSymbol(symbol)

			srcPath := filepath.Join(y.srcDir, f.Name())
			dstPath := filepath.Join(y.dstDir, "yahoo", fmt.Sprintf("%s.csv", symbol))

			y.rename(srcPath, dstPath)

		}
	}

	y.renameCompleted()

}

func (y *yahoo) Check() {}
