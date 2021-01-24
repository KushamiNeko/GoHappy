package operator

import (
	"fmt"
	"path/filepath"
)

type barchartFuturesHourly struct {
	*barchartFutures
}

func NewBarchartFuturesHourlyOperator(start int, end int) *barchartFuturesHourly {
	b := &barchartFuturesHourly{
		barchartFutures: NewBarchartFuturesOperator(start, end),
	}

	b.symbols = b.source()

	return b
}

func (b *barchartFuturesHourly) source() []string {
	return []string{
		"zn",
		"zf",
		"zt",
		"zb",
	}
}

func (b *barchartFuturesHourly) dstPath(code string) string {
	return filepath.Join(b.dstDir, "continuous", fmt.Sprintf("%s@h", code[:2]), fmt.Sprintf("%s.csv", code))
}
