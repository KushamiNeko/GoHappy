package model

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/KushamiNeko/futures_wizards/config"
	"github.com/KushamiNeko/futures_wizards/utils"
)

type TradingBook struct {
	index        string
	lastModified string
	date         string
	note         string
	bookType     string
}

func NewTradingBook(date, note, bookType string) (*TradingBook, error) {
	t := new(TradingBook)

	t.date = date
	t.note = note
	t.bookType = bookType

	t.index = utils.RandString(config.IdLen)
	t.lastModified = strconv.FormatInt(time.Now().UnixNano(), 10)

	err := t.validateInput()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewTradingBookFromInputs(entity map[string]string) (*TradingBook, error) {
	date, ok := entity["date"]
	if !ok {
		return nil, fmt.Errorf("missing date")
	}

	note, ok := entity["note"]
	if !ok {
		return nil, fmt.Errorf("missing note")
	}

	bookType, ok := entity["book_type"]
	if !ok {
		return nil, fmt.Errorf("missing bookType")
	}

	t, err := NewTradingBook(date, note, bookType)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewTradingBookFromEntity(entity map[string]string) (*TradingBook, error) {
	date, ok := entity["date"]
	if !ok {
		return nil, fmt.Errorf("missing date")
	}

	note, ok := entity["note"]
	if !ok {
		return nil, fmt.Errorf("missing note")
	}

	bookType, ok := entity["book_type"]
	if !ok {
		return nil, fmt.Errorf("missing bookType")
	}

	index, ok := entity["index"]
	if !ok {
		return nil, fmt.Errorf("missing index")
	}

	lastModified, ok := entity["last_modified"]
	if !ok {
		return nil, fmt.Errorf("missing last_modified")
	}

	t := new(TradingBook)

	t.date = date
	t.note = note
	t.bookType = bookType

	t.index = index
	t.lastModified = lastModified

	err := t.validateInput()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TradingBook) validateInput() error {

	const (
		reIndex        = `^[0-9a-zA-Z]+$`
		reLastModified = `^[0-9.]+$`
		reDate         = `^\d{8}$`
		reNote         = `^[^;]+$`
		reBookType     = `(?:paper|live|PAPER|LIVE)$`
	)

	var re *regexp.Regexp

	re = regexp.MustCompile(reDate)
	if !re.MatchString(t.date) {
		return fmt.Errorf("invalid date: %s", t.date)
	}

	re = regexp.MustCompile(reNote)
	if !re.MatchString(t.note) {
		return fmt.Errorf("invalid note: %s", t.note)
	}

	re = regexp.MustCompile(reBookType)
	if !re.MatchString(t.bookType) {
		return fmt.Errorf("invalid bookType: %s", t.bookType)
	}

	re = regexp.MustCompile(reIndex)
	if !re.MatchString(t.index) {
		return fmt.Errorf("invalid index: %s", t.index)
	}

	re = regexp.MustCompile(reLastModified)
	if !re.MatchString(t.lastModified) {
		return fmt.Errorf("invalid lastModified: %s", t.lastModified)
	}

	return nil
}

func (t *TradingBook) Index() string {
	return t.index
}

func (t *TradingBook) Date() time.Time {
	d, _ := time.Parse(config.TimeFormat, t.date)
	return d
}

func (t *TradingBook) BookType() string {
	return t.bookType
}

func (t *TradingBook) LastModified() int64 {
	d, _ := strconv.ParseInt(t.lastModified, 10, 64)
	return d
}

func (t *TradingBook) Note() string {
	return t.note
}

func (t *TradingBook) Modified() {
	t.lastModified = strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (t *TradingBook) Entity() map[string]string {
	return map[string]string{
		"index":         t.index,
		"last_modified": t.lastModified,
		"date":          t.date,
		"note":          t.note,
		"book_type":     t.bookType,
	}
}

const (
	tradingBookFmtString = "%-[3]*[4]s%-[1]*[5]s%-[1]*[6]s%[7]s"
)

func (t *TradingBook) Fmt() string {
	return fmt.Sprintf(
		tradingBookFmtString,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		t.index,
		t.bookType,
		t.date,
		t.note,
	)
}

func TradingBookFmtLabels() string {
	return fmt.Sprintf(
		tradingBookFmtString,
		config.FmtWidthL,
		config.FmtWidthXL,
		config.FmtWidthXXL,
		"Index",
		"Book Type",
		"Date",
		"Note",
	)
}
