package operator

import (
	"fmt"
	"path/filepath"
)

type futuresWorker interface {
	name() string
	source() []string
	dstPath(dstDir, code string) string
}

type dailyWorker struct {
}

func (b *dailyWorker) name() string {
	return "Futures Daily"
}

func (b *dailyWorker) source() []string {
	return []string{
		"es",
		"nq",
		"qr",
		"ym",
		"np",
		"fx",
		"zn",
		"zf",
		"zt",
		"zb",
		"ge",
		"tj",
		"gg",
		"dx",
		"e6",
		"j6",
		"b6",
		"a6",
		"d6",
		"s6",
		"n6",
		"gc",
		"si",
		"cl",
		"ng",
		"zs",
		"zc",
		"zw",
	}
}

func (b *dailyWorker) dstPath(dstDir, code string) string {
	return filepath.Join(dstDir, "continuous", code[:2], fmt.Sprintf("%s.csv", code))
}

type intraday60MinWorker struct {
}

func (b *intraday60MinWorker) name() string {
	return "Futures Intraday 60 Minutes"
}

func (b *intraday60MinWorker) source() []string {
	return []string{
		"zn",
		"zf",
		//"zt",
		"zb",
		"gg",
		"e6",
		"j6",
		//"b6",
		//"a6",
	}
}

func (b *intraday60MinWorker) dstPath(dstDir, code string) string {
	return filepath.Join(dstDir, "continuous", fmt.Sprintf("%s@h", code[:2]), fmt.Sprintf("%s.csv", code))
}

type intraday30MinWorker struct {
}

func (b *intraday30MinWorker) name() string {
	return "Futures Intraday 30 Minutes"
}

func (b *intraday30MinWorker) source() []string {
	return []string{
		"zn",
		"zf",
		"zb",
		"gg",
	}
}

func (b *intraday30MinWorker) dstPath(dstDir, code string) string {
	return filepath.Join(dstDir, "continuous", fmt.Sprintf("%s@30m", code[:2]), fmt.Sprintf("%s.csv", code))
}
