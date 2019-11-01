package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/KushamiNeko/futures_wizards/config"
)

type FuturesTrade struct {
	transactions []*FuturesTransaction
	size         int

	o []*FuturesTransaction
	c []*FuturesTransaction
}

func NewFuturesTrade(transactions []*FuturesTransaction) (*FuturesTrade, error) {
	if len(transactions) == 0 {
		return nil, fmt.Errorf("invalid transactions")
	}

	f := new(FuturesTrade)

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].TimeStamp() < transactions[j].TimeStamp()
	})

	f.transactions = transactions
	err := f.processing()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *FuturesTrade) processing() error {
	o := make([]*FuturesTransaction, 0)
	c := make([]*FuturesTransaction, 0)

	oq := 0
	cq := 0

	for _, t := range f.transactions {
		if t.Symbol() != f.Symbol() {
			return fmt.Errorf("inconsistence symbols: %s, %s", t.Symbol(), f.Symbol())
		}

		op := t.Operation()

		if op == f.Operation() {
			oq = oq + t.Quantity()
			o = append(o, t)
		} else {
			cq = cq + t.Quantity()
			c = append(c, t)
		}
	}

	if oq != cq {
		return fmt.Errorf("inconsistence quantity: %d, %d", oq, cq)
	}

	f.o = o
	f.c = c

	f.size = oq

	return nil
}

func (f *FuturesTrade) averagePrice(transactions []*FuturesTransaction) float64 {
	q := 0
	var tp float64 = 0.0

	for _, t := range transactions {
		q += t.Quantity()
		tp += t.Price() * float64(t.Quantity())
	}

	return tp / float64(q)
}

func (f *FuturesTrade) Operation() string {
	return f.transactions[0].Operation()
}

func (f *FuturesTrade) Symbol() string {
	return f.transactions[0].Symbol()
}

func (f *FuturesTrade) OpenDate() time.Time {
	return f.transactions[0].Date()
}

func (f *FuturesTrade) CloseDate() time.Time {
	return f.transactions[len(f.transactions)-1].Date()
}

func (f *FuturesTrade) OpenTimeStamp() int64 {
	return f.transactions[0].TimeStamp()
}

func (f *FuturesTrade) CloseTimeStamp() int64 {
	return f.transactions[len(f.transactions)-1].TimeStamp()
}

func (f *FuturesTrade) Size() int {
	return f.size
}

func (f *FuturesTrade) CommissionFees() float64 {
	return config.PerContractCommissionFee * float64(f.Size()) * 2
}

func (f *FuturesTrade) AvgOpenPrice() float64 {
	c := NewContractSpecs()
	unit, _ := c.LookupContractUnit(f.Symbol())
	return f.averagePrice(f.o) * float64(f.Size()) * unit
}

func (f *FuturesTrade) AvgClosePrice() float64 {
	c := NewContractSpecs()
	unit, _ := c.LookupContractUnit(f.Symbol())
	return f.averagePrice(f.c) * float64(f.Size()) * unit
}

func (f *FuturesTrade) GL() float64 {

	var o float64
	var c float64

	if f.Operation() == "+" {
		o = -1 * f.AvgOpenPrice()
		c = f.AvgClosePrice()
	} else {
		o = f.AvgOpenPrice()
		c = -1 * f.AvgClosePrice()
	}

	return o + c - f.CommissionFees()
}

func (f *FuturesTrade) GLP() float64 {
	glp := (f.GL() / f.AvgOpenPrice()) * 100.0
	return glp
}

func (f *FuturesTrade) Entity() map[string]string {
	return map[string]string{
		"operation":       f.Operation(),
		"symbol":          f.Symbol(),
		"open_date":       fmt.Sprintf("%s", f.transactions[0].date),
		"close_date":      fmt.Sprintf("%s", f.transactions[len(f.transactions)-1].date),
		"average_open":    fmt.Sprintf("%.2f", f.AvgOpenPrice()),
		"average_close":   fmt.Sprintf("%.2f", f.AvgClosePrice()),
		"gl":              fmt.Sprintf("%.2f", f.GL()),
		"glp":             fmt.Sprintf("%.2f", f.GLP()),
		"commission_fees": fmt.Sprintf("%.2f", f.CommissionFees()),
		"size":            fmt.Sprintf("%d", f.Size()),
	}
}

const (
	futuresTradeFmtString = "%-[1]*[4]s%-[2]*[5]s%-[1]*[6]s%-[2]*[7]s%-[2]*[8]s%-[3]*[9]s%-[3]*[10]s%-[3]*[11]s%-[3]*[12]s%-[3]*[13]s"
)

func (f *FuturesTrade) Fmt() string {
	return fmt.Sprintf(
		futuresTradeFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		f.Symbol(),
		f.Operation(),
		fmt.Sprintf("%d", f.Size()),
		fmt.Sprintf("%s", f.transactions[0].date),
		fmt.Sprintf("%s", f.transactions[len(f.transactions)-1].date),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.AvgOpenPrice()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.AvgClosePrice()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.CommissionFees()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.GL()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.GLP()),
	)
}

func FuturesTradeFmtLabels() string {
	return fmt.Sprintf(
		futuresTradeFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		"Symbol",
		"Operation",
		"Size",
		"Open Date",
		"Close Date",
		"Avg Open Price",
		"Avg Close Price",
		"Commission Fees",
		"GL($)",
		"GL(%)",
	)
}
