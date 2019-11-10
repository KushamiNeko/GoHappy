package agent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/KushamiNeko/go_fun/utils/input"
)

type CalculatingAgent struct{}

func NewCalculatingAgent() *CalculatingAgent {
	c := new(CalculatingAgent)

	return c
}

func (c *CalculatingAgent) validateInputs(inputs map[string]string, keys []string) error {
	const (
		reFloat     = `^[0-9.]+$`
		reOperation = `^[+-]$`
	)

	book := map[string]string{
		"price":     reFloat,
		"percent":   reFloat,
		"operation": reOperation,
		"start":     reFloat,
		"end":       reFloat,
	}

	var re *regexp.Regexp

	for _, k := range keys {
		v, ok := inputs[k]
		if !ok {
			return fmt.Errorf("missing key: %s", k)
		}

		r, ok := book[k]
		if !ok {
			panic(fmt.Sprintf("program error: missing regex for key: %s", k))
		}

		re = regexp.MustCompile(r)
		if !re.MatchString(v) {
			return fmt.Errorf("invalid input: %s=%s", k, v)
		}
	}

	return nil
}

func (c *CalculatingAgent) Profit(inputs map[string]string) (float64, error) {
	n := input.InputsAbbreviation(
		inputs,
		map[string]string{
			"p": "price",
			"%": "percent",
			"o": "operation",
		},
	)

	err := c.validateInputs(n, []string{"price", "percent", "operation"})
	if err != nil {
		return 0, err
	}

	price, _ := strconv.ParseFloat(n["price"], 64)

	percent, _ := strconv.ParseFloat(n["percent"], 64)

	op := n["operation"]

	var tar float64

	if op == "+" {
		tar = price * (1.0 + (percent / 100.0))
	} else {
		tar = price * (1.0 - (percent / 100.0))
	}

	return tar, nil
}

func (c *CalculatingAgent) Stop(inputs map[string]string) (float64, error) {
	n := input.InputsAbbreviation(
		inputs,
		map[string]string{
			"p": "price",
			"%": "percent",
			"o": "operation",
		},
	)

	err := c.validateInputs(n, []string{"price", "percent", "operation"})
	if err != nil {
		return 0, err
	}

	price, _ := strconv.ParseFloat(n["price"], 64)

	percent, _ := strconv.ParseFloat(n["percent"], 64)

	op := n["operation"]

	var tar float64

	if op == "+" {
		tar = price * (1.0 - (percent / 100.0))
	} else {
		tar = price * (1.0 + (percent / 100.0))
	}

	return tar, nil
}

func (c *CalculatingAgent) Depth(inputs map[string]string) (float64, error) {
	n := input.InputsAbbreviation(
		inputs,
		map[string]string{
			"s": "start",
			"e": "end",
		},
	)

	err := c.validateInputs(n, []string{"start", "end"})
	if err != nil {
		return 0, err
	}

	start, _ := strconv.ParseFloat(n["start"], 64)

	end, _ := strconv.ParseFloat(n["end"], 64)

	tar := ((end - start) / start) * 100.0

	return tar, nil
}

func (c *CalculatingAgent) Pivot() error {

	return nil
}

func (c *CalculatingAgent) Fib() error {

	return nil
}

func (c *CalculatingAgent) Keltner() error {

	return nil
}
