package agent

import (
	"fmt"
	"sort"
	"strings"

	"github.com/KushamiNeko/futures_wizards/config"
	"github.com/KushamiNeko/futures_wizards/context"
	"github.com/KushamiNeko/futures_wizards/model"
	"github.com/KushamiNeko/futures_wizards/utils"
)

type TradingAgent struct {
	ctx *context.Context

	books   []*model.TradingBook
	reading *model.TradingBook

	bookType  string
	tradingDB string
}

func NewTradingAgent(ctx *context.Context, live bool) (*TradingAgent, error) {
	t := new(TradingAgent)
	t.ctx = ctx

	if live {
		t.bookType = "live"
		t.tradingDB = config.DbLiveTrading
	} else {
		t.bookType = "paper"
		t.tradingDB = config.DbPaperTrading
	}

	books, err := t.ctx.Db().Find(
		config.DbTradingBooks,
		t.ctx.User().Uid(),
		map[string]string{
			"book_type": t.bookType,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(books) > 0 {

		t.books = make([]*model.TradingBook, len(books))

		for i, b := range books {
			t.books[i], err = model.NewTradingBookFromEntity(b)
			if err != nil {
				return nil, err
			}
		}

		sort.Slice(t.books, func(i, j int) bool {
			return t.books[i].LastModified() > t.books[j].LastModified()
		})

		t.reading = t.books[0]
	}

	return t, nil
}

func (t *TradingAgent) NewBook(inputs map[string]string) error {
	n := utils.InputsAbbreviation(inputs, map[string]string{
		"d": "date",
		"n": "note",
	})

	n["book_type"] = t.bookType

	book, err := model.NewTradingBookFromInputs(n)
	if err != nil {
		return err
	}

	err = t.ctx.Db().Insert(config.DbTradingBooks, t.ctx.User().Uid(), book.Entity())
	if err != nil {
		return err
	}

	t.books = append(t.books, book)
	sort.Slice(t.books, func(i, j int) bool {
		return t.books[i].LastModified() > t.books[j].LastModified()
	})

	t.reading = book

	return nil
}

func (t *TradingAgent) UpdateBook() error {

	t.reading.Modified()

	err := t.ctx.Db().Replace(
		config.DbTradingBooks,
		t.ctx.User().Uid(),
		map[string]string{
			"index": t.reading.Index(),
		},
		t.reading.Entity(),
	)

	return err
}

func (t *TradingAgent) ChangeBook(inputs map[string]string) error {
	n := utils.InputsAbbreviation(inputs, map[string]string{
		"i": "index",
	})

	i, ok := n["index"]
	if !ok {
		return fmt.Errorf("missing index")
	}

	if i == "" {
		return fmt.Errorf("invalid book index")
	}

	for _, b := range t.books {
		if strings.HasPrefix(b.Index(), i) {
			t.reading = b
			return nil
		}
	}

	return fmt.Errorf("unknown book index: %s", i)
}

func (t *TradingAgent) NewTransaction(inputs map[string]string) error {
	n := utils.InputsAbbreviation(inputs, map[string]string{
		"d": "date",
		"s": "symbol",
		"o": "operation",
		"q": "quantity",
		"p": "price",
		"n": "note",
	})

	transaction, err := model.NewFuturesTransactionFromInputs(n)
	if err != nil {
		return err
	}

	err = t.ctx.Db().Insert(
		t.tradingDB,
		t.reading.Index(),
		transaction.Entity(),
	)
	if err != nil {
		return err
	}

	err = t.UpdateBook()
	if err != nil {
		return err
	}

	return nil
}

func (t *TradingAgent) Positions() ([]*model.FuturesTransaction, error) {
	results, err := t.ctx.Db().Find(
		t.tradingDB,
		t.reading.Index(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.FuturesTransaction, len(results))

	for i, r := range results {
		t, err := model.NewFuturesTransactionFromEntity(r)
		if err != nil {
			return nil, err
		}

		transactions[i] = t
	}

	trades, err := t.processTrades(transactions)
	if err != nil {
		return nil, err
	}

	if len(trades) == 0 {
		return transactions, nil
	}

	positions := make([]*model.FuturesTransaction, 0)

	for _, t := range transactions {
		if t.TimeStamp() > trades[len(trades)-1].CloseTimeStamp() {
			positions = append(positions, t)
		}
	}

	return positions, nil
}

func (t *TradingAgent) Books() ([]*model.TradingBook, error) {
	if t.books == nil || len(t.books) == 0 {
		return nil, fmt.Errorf("empty books")
	} else {
		return t.books, nil
	}
}

func (t *TradingAgent) Transactions() ([]*model.FuturesTransaction, error) {
	results, err := t.ctx.Db().Find(
		t.tradingDB,
		t.reading.Index(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.FuturesTransaction, len(results))

	for i, r := range results {
		t, err := model.NewFuturesTransactionFromEntity(r)
		if err != nil {
			return nil, err
		}

		transactions[i] = t
	}

	return transactions, nil
}

func (t *TradingAgent) Trades() ([]*model.FuturesTrade, error) {
	results, err := t.ctx.Db().Find(
		t.tradingDB,
		t.reading.Index(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.FuturesTransaction, len(results))

	for i, r := range results {
		t, err := model.NewFuturesTransactionFromEntity(r)
		if err != nil {
			return nil, err
		}

		transactions[i] = t
	}

	trades, err := t.processTrades(transactions)
	if err != nil {
		return nil, err
	}

	return trades, nil
}

func (t *TradingAgent) Statistic() (*model.Statistic, error) {
	trades, err := t.Trades()
	if err != nil {
		return nil, err
	}

	return model.NewStatistic(trades)
}

func (t *TradingAgent) processTrades(transactions []*model.FuturesTransaction) ([]*model.FuturesTrade, error) {

	type orders struct {
		t []*model.FuturesTransaction
		p int
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].TimeStamp() < transactions[j].TimeStamp()
	})

	trades := make([]*model.FuturesTrade, 0)
	book := make(map[string]*orders)

	for _, t := range transactions {
		o, ok := book[t.Symbol()]
		if !ok {
			book[t.Symbol()] = &orders{
				t: []*model.FuturesTransaction{t},
				p: t.Action(),
			}
		} else {
			o.p += t.Action()
			o.t = append(book[t.Symbol()].t, t)
		}

		o = book[t.Symbol()]

		if o.p == 0 {
			trade, err := model.NewFuturesTrade(o.t)
			if err != nil {
				return nil, err
			}

			trades = append(trades, trade)
			delete(book, t.Symbol())
		}
	}

	return trades, nil
}

func (t *TradingAgent) Reading() (*model.TradingBook, error) {
	if t.reading != nil {
		return t.reading, nil
	} else {
		return nil, fmt.Errorf("empty books")
	}
}

func (t *TradingAgent) SetReading(book *model.TradingBook) error {
	found := false
	for _, b := range t.books {
		if b.Date().Equal(book.Date()) && b.Note() == book.Note() {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid book")
	}

	t.reading = book

	return nil
}

//func (t *TradingAgent) Plot(inputs map[string]string) error {
//n := utils.InputsAbbreviation(inputs, map[string]string{
//"p": "period",
//"f": "file",
//})

//p := n["period"]

//file := n["file"]
//if _, err := os.Stat(filepath.Dir(file)); os.IsNotExist(err) {
//return fmt.Errorf("parent folder does not exist: %s", file)
//}

//ext := filepath.Ext(file)

//if ext == "" {
//file = fmt.Sprintf("%s.png", file)
//} else if ext != "png" {
//file = strings.Replace(file, ext, ".png", -1)
//}

//tcs, _ := t.Transactions()
//if len(tcs) == 0 {
//return fmt.Errorf("empty book: %s", t.reading)
//}

//sort.Slice(tcs, func(i, j int) bool {
//return tcs[i].TimeStamp() < tcs[j].TimeStamp()
//})

//tds, _ := t.Trades()

//trs := make([]*data.TradeRecord, 0, len(tcs))

//cur := 0
//lastOp := ""
//for _, tc := range tcs {
//date := tc.Date()

//var op string

//if cur >= len(tds) {
//if lastOp == "" {
//switch tc.Operation() {
//case "+":
//op = "long"
//case "-":
//op = "short"
//default:
//panic(fmt.Sprintf("unknown operation: %s", tc.Operation()))
//}

//trs = append(trs, data.NewTradeRecord(date, op))
//lastOp = tc.Operation()
//} else {
//switch {
//case tc.Operation() == lastOp:
//op = "increase"
//case tc.Operation() != lastOp:
//op = "decrease"
//default:
//panic(fmt.Sprintf("unknown operation: %s", tc.Operation()))
//}

//trs = append(trs, data.NewTradeRecord(date, op))
//}

//continue
//}

//s := tds[cur].OpenDate()
//e := tds[cur].CloseDate()

//switch {
//case tc.Date().Equal(s):
//switch tc.Operation() {
//case "+":
//op = "long"
//case "-":
//op = "short"
//default:
//panic(fmt.Sprintf("unknown operation: %s", tc.Operation()))
//}

//trs = append(trs, data.NewTradeRecord(date, op))

//case tc.Date().Equal(e):
//op = "close"
//trs = append(trs, data.NewTradeRecord(date, op))
//cur++

//case tc.Date().After(s) && tc.Date().Before(e):
//switch tds[cur].Operation() {
//case "+":
//switch tc.Operation() {
//case "+":
//op = "increase"
//case "-":
//op = "decrease"
//default:
//panic(fmt.Sprintf("unknown operation: %s", tc.Operation()))
//}
//trs = append(trs, data.NewTradeRecord(date, op))
//case "-":
//switch tc.Operation() {
//case "+":
//op = "decrease"
//case "-":
//op = "increase"
//default:
//panic(fmt.Sprintf("unknown operation: %s", tc.Operation()))
//}
//trs = append(trs, data.NewTradeRecord(date, op))
//}
//}
//}

//var s, e time.Time

//d := t.reading.Date()

//rd := time.Date(d.Year()+1, d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
//td := tcs[len(tcs)-1].Date()

//if rd.After(td) {
//e = rd
//} else {
//e = td
//}

////if p == chart.Weekly {
////for j := 0; j < 31; j++ {
////e = e.Add(24 * time.Hour)
////if e.Weekday() != time.Friday {
////break
////}
////}
////}

//switch p {
//case data.DailyChart:
//rd = t.reading.Date()
//td = tcs[0].Date()

//if rd.Before(td) {
//s = rd
//} else {
//s = td
//}

//for j := 0; j > -31; j-- {
//s = s.Add(-24 * time.Hour)
//if s.Weekday() != time.Saturday && s.Weekday() != time.Sunday {
//break
//}
//}
//case data.WeeklyChart:
//s = time.Date(d.Year()-3, d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
//for j := 0; j > -31; j-- {
//s = s.Add(-24 * time.Hour)
//if s.Weekday() == time.Monday {
//break
//}
//}
//case data.MonthlyChart:
//s = time.Date(d.Year()-18, d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())

//default:
//return fmt.Errorf("unknown period: %s", p)
//}

//var sym string

//switch tds[0].Symbol() {
//case "rty":
//sym = "rut"
//case "es":
//sym = "spx"
//default:
//return fmt.Errorf("unknown symbol: %s", tds[0].Symbol())
//}

//c, err := painter.NewChart(sym, p, s, e, trs)
//if err != nil {
//return err
//}

//f, err := os.Create(file)
//if err != nil {
//return err
//}
//defer f.Close()

//err = c.Plot(f)
//if err != nil {
//return err
//}

//return nil
//}

//func (t *TradingAgent) PlotAllBooks(inputs map[string]string) error {
//n := utils.InputsAbbreviation(inputs, map[string]string{
//"f": "folder",
//})

//folder := n["folder"]

//if _, err := os.Stat(folder); os.IsNotExist(err) {
//return fmt.Errorf("folder does not exist: %s", folder)
//}

//reading := t.reading

//books, err := t.Books()
//if err != nil {
//return err
//}

//for _, book := range books {
//t.reading = book

//tcs, _ := t.Transactions()
//if len(tcs) == 0 {
//continue
//}

//for _, p := range []string{"d", "w"} {
//notes := strings.Split(book.Note(), " ")
//var name string

//if len(notes) == 2 {
//name = fmt.Sprintf("%d_%s_%s.png", book.Date().Year(), strings.Join(notes, "_"), p)
//} else if len(notes) == 3 {
//name = fmt.Sprintf("%d_%s_%s_%s.png", book.Date().Year(), strings.Join(notes[:2], "_"), p, notes[2])
//} else {
//return fmt.Errorf("program error: unknown book note: %s", book.Note())
//}
//file := filepath.Join(folder, name)

//in := map[string]string{
//"p": p,
//"f": file,
//}

//err = t.Plot(in)
//if err != nil {
//return err
//}

//}
//}

//t.reading = reading

//return nil
//}
