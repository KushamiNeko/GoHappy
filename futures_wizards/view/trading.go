package view

import (
	"fmt"

	"github.com/KushamiNeko/go_fun/trading/model"
	"github.com/KushamiNeko/go_fun/utils/input"
	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/futures_wizards/agent"
	"github.com/KushamiNeko/go_happy/futures_wizards/config"
	"github.com/KushamiNeko/go_happy/futures_wizards/context"
)

type Trading struct {
	Page

	ctx *context.Context

	agent *agent.TradingAgent
}

func newTrading(ctx *context.Context) (*Trading, error) {
	t := new(Trading)
	t.ctx = ctx

	t.Page.actions = []string{
		"calculator",
		"books",
		"reading",
		"new book",
		"change book",
		"new transaction",
		"position",
		"transactions",
		"trades",
		"statistic",
		"plot",
		//"plot all",
	}

	t.Page.handlers = map[string]func() error{
		"calculator":      t.cmdCalculator,
		"books":           t.cmdBooks,
		"reading":         t.cmdReading,
		"new book":        t.cmdNewBook,
		"change book":     t.cmdChangeBook,
		"new transaction": t.cmdNewTransaction,
		"position":        t.cmdPosition,
		"transactions":    t.cmdTransactions,
		"trades":          t.cmdTrades,
		"statistic":       t.cmdStatistic,
		"plot":            t.cmdPlot,
		//"plot all":        t.cmdPlotAll,
	}

	t.Page.init(false)

	return t, nil
}

func NewPaperTrading(ctx *context.Context) (*Trading, error) {
	var err error

	t, err := newTrading(ctx)
	if err != nil {
		return nil, err
	}

	t.agent, err = agent.NewTradingAgent(ctx, false)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewLiveTrading(ctx *context.Context) (*Trading, error) {
	var err error

	t, err := newTrading(ctx)
	if err != nil {
		return nil, err
	}

	t.agent, err = agent.NewTradingAgent(ctx, true)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Trading) Main() {

	err := t.cmdReading()
	if err != nil {
		pretty.ColorPrintln(
			config.ColorWarnings,
			err.Error(),
		)
	}

	t.Page.Main()
}

func (t *Trading) cmdCalculator() error {
	p := NewCalculator(t.ctx)

	p.Main()

	return nil
}

func (t *Trading) cmdBooks() error {
	books, err := t.agent.Books()
	if err != nil {
		return err
	}

	pretty.ColorPrintln(config.ColorInfo, model.TradingBookFmtLabels())

	for _, book := range books {
		pretty.ColorPrintln(config.ColorWhite, book.Fmt())
	}

	return nil
}

func (t *Trading) cmdReading() error {
	book, err := t.agent.Reading()
	if err != nil {
		return err
	}

	pretty.ColorPrintln(config.ColorInfo, model.TradingBookFmtLabels())
	pretty.ColorPrintln(config.ColorWhite, book.Fmt())

	return nil
}

func (t *Trading) cmdNewBook() error {
	inputs, err := input.KeyValueInput(
		config.ColorInfo,
		"create a new trading book: (time, note, version)",
	)
	if err != nil {
		return err
	}

	err = t.agent.NewBook(inputs)
	if err != nil {
		return err
	}

	t.cmdReading()

	return nil
}

func (t *Trading) cmdChangeBook() error {
	inputs, err := input.KeyValueInput(
		config.ColorInfo,
		"change to another book: (index)",
	)
	if err != nil {
		return err
	}

	err = t.agent.ChangeBook(inputs)
	if err != nil {
		return err
	}

	t.cmdReading()

	return nil
}

func (t *Trading) cmdNewTransaction() error {
	inputs, err := input.KeyValueInput(
		config.ColorInfo,
		"new transaction: (time, symbol, operation, quantity, price, note)",
	)
	if err != nil {
		return err
	}

	err = t.agent.NewTransaction(inputs)
	if err != nil {
		return err
	}

	err = t.cmdPosition()
	if err != nil {
		return err
	}

	return nil
}

func (t *Trading) cmdPosition() error {
	positions, err := t.agent.Positions()
	if err != nil {
		return err
	}

	if len(positions) == 0 {
		trades, err := t.agent.Trades()
		if err != nil {
			return err
		}

		last := trades[len(trades)-1]

		pretty.ColorPrintln(
			config.ColorInfo,
			fmt.Sprintf(
				"last trade GL($, %%): $%.[1]*[2]f, %.[1]*[3]f%%",
				config.FloatDecimals,
				last.GL(),
				last.GLP(),
			),
		)
		return fmt.Errorf("empty position")
	}

	pretty.ColorPrintln(config.ColorInfo, model.FuturesTransactionFmtLabels())

	for _, p := range positions {
		pretty.ColorPrintln(config.ColorWhite, p.Fmt())
	}

	return nil
}

func (t *Trading) cmdTransactions() error {
	transactions, err := t.agent.Transactions()
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return fmt.Errorf("empty transaction")
	}

	pretty.ColorPrintln(config.ColorInfo, model.FuturesTransactionFmtLabels())

	for _, v := range transactions {
		pretty.ColorPrintln(config.ColorWhite, v.Fmt())
	}

	return nil
}

func (t *Trading) cmdTrades() error {
	trades, err := t.agent.Trades()
	if err != nil {
		return err
	}

	if len(trades) == 0 {
		return fmt.Errorf("empty trade")
	}

	pretty.ColorPrintln(config.ColorInfo, model.FuturesTradeFmtLabels())

	for _, v := range trades {
		pretty.ColorPrintln(config.ColorWhite, v.Fmt())
	}

	return nil
}

func (t *Trading) cmdStatistic() error {
	statistic, err := t.agent.Statistic()
	if err != nil {
		return err
	}

	pretty.ColorPrintln(config.ColorInfo, model.StatisticFmtLabels())
	pretty.ColorPrintln(config.ColorWhite, statistic.Fmt())

	pretty.ColorPrintln(config.ColorInfo, model.StatisticFmtLabelsL())
	pretty.ColorPrintln(config.ColorWhite, statistic.FmtL())

	pretty.ColorPrintln(config.ColorInfo, model.StatisticFmtLabelsS())
	pretty.ColorPrintln(config.ColorWhite, statistic.FmtS())

	return nil
}

func (t *Trading) cmdPlot() error {
	//inputs, err := input.KeyValueInput(
	//config.ColorInfo,
	//"plot candlestick: (period, file)",
	//)
	//if err != nil {
	//return err
	//}

	//return t.agent.Plot(inputs)
	return nil
}

func (t *Trading) cmdPlotAll() error {
	//inputs, err := input.KeyValueInput(
	//config.ColorInfo,
	//"plot all books: (folder)",
	//)
	//if err != nil {
	//return err
	//}

	//return t.agent.PlotAllBooks(inputs)
	return nil
}
