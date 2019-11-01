package model

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/KushamiNeko/futures_wizards/config"
	"github.com/KushamiNeko/futures_wizards/utils"
)

type FuturesTransaction struct {
	index     string
	timeStamp string

	date      string
	symbol    string
	operation string
	quantity  string
	price     string
	note      string
}

func NewFuturesTransaction(
	date string,
	symbol string,
	operation string,
	quantity string,
	price string,
	note string,
) (*FuturesTransaction, error) {

	f := new(FuturesTransaction)
	f.date = date
	f.symbol = symbol
	f.operation = operation
	f.quantity = quantity
	f.price = price
	f.note = note

	f.index = utils.RandString(config.IdLen)
	f.timeStamp = strconv.FormatInt(time.Now().UnixNano(), 10)

	err := f.validateInput()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func NewFuturesTransactionFromInputs(entity map[string]string) (*FuturesTransaction, error) {

	date, ok := entity["date"]
	if !ok {
		return nil, fmt.Errorf("missing date")
	}

	symbol, ok := entity["symbol"]
	if !ok {
		return nil, fmt.Errorf("missing symbol")
	}

	operation, ok := entity["operation"]
	if !ok {
		return nil, fmt.Errorf("missing operation")
	}

	quantity, ok := entity["quantity"]
	if !ok {
		return nil, fmt.Errorf("missing quantity")
	}

	price, ok := entity["price"]
	if !ok {
		return nil, fmt.Errorf("missing price")
	}

	note, _ := entity["note"]
	//if !ok {
	//return nil, fmt.Errorf("missing note")
	//}

	f, err := NewFuturesTransaction(date, symbol, operation, quantity, price, note)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func NewFuturesTransactionFromEntity(entity map[string]string) (*FuturesTransaction, error) {

	date, ok := entity["date"]
	if !ok {
		return nil, fmt.Errorf("missing date")
	}

	symbol, ok := entity["symbol"]
	if !ok {
		return nil, fmt.Errorf("missing symbol")
	}

	operation, ok := entity["operation"]
	if !ok {
		return nil, fmt.Errorf("missing operation")
	}

	quantity, ok := entity["quantity"]
	if !ok {
		return nil, fmt.Errorf("missing quantity")
	}

	price, ok := entity["price"]
	if !ok {
		return nil, fmt.Errorf("missing price")
	}

	note, _ := entity["note"]
	//if !ok {
	//return nil, fmt.Errorf("missing note")
	//}

	index, ok := entity["index"]
	if !ok {
		return nil, fmt.Errorf("missing index")
	}

	timeStamp, ok := entity["time_stamp"]
	if !ok {
		return nil, fmt.Errorf("missing timeStamp")
	}

	f := new(FuturesTransaction)

	f.date = date
	f.symbol = symbol
	f.operation = operation
	f.quantity = quantity
	f.price = price
	f.note = note

	f.index = index
	f.timeStamp = timeStamp

	err := f.validateInput()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *FuturesTransaction) validateInput() error {

	const (
		reDate      = `^\d{8}$`
		reSymbol    = `^[a-z]+$`
		reOperation = `^[+-]$`
		reQuantity  = `^[0-9]+$`
		rePrice     = `^[0-9.]+$`
		reNote      = `^[^;]*$`
		reIndex     = `^[a-zA-Z0-9]+$`
		reTimeStamp = `^[0-9.]+$`
	)

	var re *regexp.Regexp

	re = regexp.MustCompile(reDate)
	if !re.MatchString(f.date) {
		return fmt.Errorf("invalid date: %s", f.date)
	}

	re = regexp.MustCompile(reSymbol)
	if !re.MatchString(f.symbol) {
		return fmt.Errorf("invalid symbol: %s", f.symbol)
	}

	c := NewContractSpecs()
	if !c.ValidateSymbol(f.symbol) {
		return fmt.Errorf("invalid symbol")
	}

	re = regexp.MustCompile(reOperation)
	if !re.MatchString(f.operation) {
		return fmt.Errorf("invalid operation: %s", f.operation)
	}

	re = regexp.MustCompile(reQuantity)
	if !re.MatchString(f.quantity) {
		return fmt.Errorf("invalid quantity: %s", f.quantity)
	}

	re = regexp.MustCompile(rePrice)
	if !re.MatchString(f.price) {
		return fmt.Errorf("invalid price: %s", f.price)
	}

	re = regexp.MustCompile(reNote)
	if !re.MatchString(f.note) {
		return fmt.Errorf("invalid note: %s", f.note)
	}

	re = regexp.MustCompile(reIndex)
	if !re.MatchString(f.index) {
		return fmt.Errorf("invalid index: %s", f.index)
	}

	re = regexp.MustCompile(reTimeStamp)
	if !re.MatchString(f.timeStamp) {
		return fmt.Errorf("invalid timeStamp: %s", f.timeStamp)
	}

	return nil
}

func (f *FuturesTransaction) Index() string {
	return f.index
}

func (f *FuturesTransaction) TimeStamp() int64 {
	t, _ := strconv.ParseInt(f.timeStamp, 10, 64)
	return t
}

func (f *FuturesTransaction) Date() time.Time {
	d, _ := time.Parse(config.TimeFormat, f.date)
	return d
}

func (f *FuturesTransaction) Symbol() string {
	return f.symbol
}

func (f *FuturesTransaction) Operation() string {
	return f.operation
}

func (f *FuturesTransaction) Quantity() int {
	q, _ := strconv.ParseInt(f.quantity, 10, 64)
	return int(q)
}

func (f *FuturesTransaction) Price() float64 {
	p, _ := strconv.ParseFloat(f.price, 64)
	return p
}

func (f *FuturesTransaction) Note() string {
	return f.note
}

func (f *FuturesTransaction) TotalPrice() float64 {
	return f.Price() * float64(f.Quantity())
}

func (f *FuturesTransaction) Action() int {
	a, _ := strconv.ParseInt(fmt.Sprintf("%s%s", f.operation, f.quantity), 10, 64)
	return int(a)
}

func (f *FuturesTransaction) Entity() map[string]string {
	return map[string]string{
		"index":      f.index,
		"time_stamp": f.timeStamp,
		"date":       f.date,
		"symbol":     f.symbol,
		"operation":  f.operation,
		"quantity":   f.quantity,
		"price":      f.price,
		"note":       f.note,
	}
}

const (
	futuresTransactionFmtString = "%-[2]*[4]s%-[1]*[5]s%-[2]*[6]s%-[2]*[7]s%-[2]*[8]s%[9]s"
)

func (f *FuturesTransaction) Fmt() string {
	return fmt.Sprintf(
		futuresTransactionFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		f.date,
		f.symbol,
		f.operation,
		f.quantity,
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.Price()),
		f.note,
	)
}

func FuturesTransactionFmtLabels() string {
	return fmt.Sprintf(
		futuresTransactionFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		"Date",
		"Symbol",
		"Operation",
		"Quantity",
		"Price",
		"Note",
	)
}
