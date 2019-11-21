package agent

import (
	"fmt"
	"sort"
	"strings"

	"github.com/KushamiNeko/go_fun/trading/model"
	"github.com/KushamiNeko/go_fun/utils/input"
	"github.com/KushamiNeko/go_happy/futures_wizards/config"
	"github.com/KushamiNeko/go_happy/futures_wizards/context"
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
	n := input.InputsAbbreviation(inputs, map[string]string{
		//"t": "time",
		//"n": "note",
		//"v": "version",
		"t": "title",
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
			//"index": t.reading.Index(),
			"record_index": t.reading.RecordIndex(),
		},
		t.reading.Entity(),
	)

	return err
}

func (t *TradingAgent) ChangeBook(inputs map[string]string) error {
	n := input.InputsAbbreviation(inputs, map[string]string{
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
		//if strings.HasPrefix(b.Index(), i) {
		if strings.HasPrefix(b.RecordIndex(), i) {
			t.reading = b
			return nil
		}
	}

	return fmt.Errorf("unknown book index: %s", i)
}

func (t *TradingAgent) NewTransaction(inputs map[string]string) error {
	n := input.InputsAbbreviation(inputs, map[string]string{
		"t": "time",
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
		//t.reading.Index(),
		t.reading.RecordIndex(),
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
		//t.reading.Index(),
		t.reading.RecordIndex(),
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

	tb := make(map[string]*model.FuturesTrade)

	for i := len(trades) - 1; i >= 0; i-- {
		t := trades[i]

		if v, ok := tb[t.Symbol()]; !ok {
			tb[t.Symbol()] = t
		} else {
			if t.CloseTime().Equal(v.CloseTime()) {
				if t.CloseTimeStamp() > v.CloseTimeStamp() {
					tb[t.Symbol()] = t
				}
			} else {
				if t.CloseTime().After(v.CloseTime()) {
					tb[t.Symbol()] = t
				}
			}
		}
	}

	positions := make([]*model.FuturesTransaction, 0)

	for _, t := range transactions {
		if v, ok := tb[t.Symbol()]; !ok {
			positions = append(positions, t)
		} else {
			if t.Time().Equal(v.CloseTime()) {
				if t.TimeStamp() > v.CloseTimeStamp() {
					positions = append(positions, t)
				}
			} else {
				if t.Time().After(v.CloseTime()) {
					positions = append(positions, t)
				}
			}
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
		//t.reading.Index(),
		t.reading.RecordIndex(),
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
		//t.reading.Index(),
		t.reading.RecordIndex(),
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
		if transactions[i].Time().Equal(transactions[j].Time()) {
			return transactions[i].TimeStamp() < transactions[j].TimeStamp()
		} else {
			return transactions[i].Time().Before(transactions[j].Time())
		}
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
		//if b.Time().Equal(book.Time()) && b.Note() == book.Note() {
		if b.Title() == book.Title() {
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
