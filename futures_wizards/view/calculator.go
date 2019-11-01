package view

import (
	"fmt"

	"github.com/KushamiNeko/futures_wizards/agent"
	"github.com/KushamiNeko/futures_wizards/config"
	"github.com/KushamiNeko/futures_wizards/context"
	"github.com/KushamiNeko/futures_wizards/utils"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

type Calculator struct {
	Page

	ctx *context.Context

	agent *agent.CalculatingAgent
}

func NewCalculator(ctx *context.Context) *Calculator {
	c := new(Calculator)
	c.ctx = ctx

	c.Page.actions = []string{
		"profit",
		"stop",
		"depth",
	}

	c.Page.handlers = map[string]func() error{
		"profit":          c.cmdProfit,
		"stop":            c.cmdStop,
		"depth":           c.cmdDepth,
		"pivot":           c.cmdPivot,
		"fib retracemane": c.cmdFib,
		"keltner channel": c.cmdKeltner,
	}

	c.Page.init(false)

	c.agent = agent.NewCalculatingAgent()

	return c
}

func (c *Calculator) cmdProfit() error {
	inputs, err := utils.KeyValueInput(
		config.ColorInfo,
		"calculating profit target: (price, %, operation)",
	)
	if err != nil {
		return err
	}

	tar, err := c.agent.Profit(inputs)
	if err != nil {
		return err
	}

	pretty.ColorPrintln(config.ColorWhite, fmt.Sprintf("%.[1]*f", config.FloatDecimals, tar))

	return nil
}

func (c *Calculator) cmdStop() error {
	inputs, err := utils.KeyValueInput(
		config.ColorInfo,
		"calculating stop target: (price, %, operation)",
	)
	if err != nil {
		return err
	}

	tar, err := c.agent.Stop(inputs)
	if err != nil {
		return err
	}

	pretty.ColorPrintln(config.ColorWhite, fmt.Sprintf("%.[1]*f", config.FloatDecimals, tar))

	return nil
}

func (c *Calculator) cmdDepth() error {
	inputs, err := utils.KeyValueInput(
		config.ColorInfo,
		"calculating depth: (start, end)",
	)
	if err != nil {
		return err
	}

	tar, err := c.agent.Depth(inputs)
	if err != nil {
		return err
	}

	pretty.ColorPrintln(config.ColorWhite, fmt.Sprintf("%.[1]*f%%", config.FloatDecimals, tar))

	return nil
}

func (c *Calculator) cmdPivot() error {

	return nil
}

func (c *Calculator) cmdFib() error {

	return nil
}

func (c *Calculator) cmdKeltner() error {

	return nil
}
