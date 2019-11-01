package view

import (
	"github.com/KushamiNeko/futures_wizards/context"
)

type Entry struct {
	Page

	ctx *context.Context
}

func NewEntry(ctx *context.Context) *Entry {
	e := new(Entry)
	e.ctx = ctx

	e.Page.actions = []string{
		"calculator",
		"paper trading",
		"live trading",
	}

	e.Page.handlers = map[string]func() error{
		"calculator":    e.cmdCalculator,
		"paper trading": e.cmdPaperTrading,
		"live trading":  e.cmdLiveTrading,
	}

	e.Page.init(true)

	return e
}

func (e *Entry) cmdCalculator() error {
	p := NewCalculator(e.ctx)

	p.Main()

	return nil
}

func (e *Entry) cmdPaperTrading() error {

	p, err := NewPaperTrading(e.ctx)
	if err != nil {
		return err
	}

	p.Main()

	return nil
}

func (e *Entry) cmdLiveTrading() error {

	p, err := NewLiveTrading(e.ctx)
	if err != nil {
		return err
	}

	p.Main()

	return nil
}
