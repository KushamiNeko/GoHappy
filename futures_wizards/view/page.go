package view

import (
	"fmt"
	"os"
	"strings"

	"github.com/KushamiNeko/go_fun/utils/foreign"
	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/futures_wizards/config"
)

type Pages interface {
	Main()
}

type Page struct {
	actions  []string
	handlers map[string]func() error
	next     bool
}

func (p *Page) init(top bool) {
	if !top {
		p.handlers["back"] = p.cmdBack
		p.actions = append(p.actions, "back")
	}

	p.handlers["exit"] = p.cmdExit
	p.actions = append(p.actions, "exit")

	p.next = false
}

func (p *Page) Main() {
	for {
		cmd, err := p.readCommand()
		if err != nil {
			pretty.ColorPrintln(config.ColorWarnings, err.Error())
			continue
		}

		hdl, ok := p.handlers[cmd]
		if !ok {
			pretty.ColorPrintln(
				config.ColorWarnings,
				fmt.Sprintf("unknown command: %s", cmd),
			)
			continue
		}

		err = hdl()
		if err != nil {
			pretty.ColorPrintln(config.ColorWarnings, err.Error())
			continue
		}

		if p.next {
			break
		}
	}
}

func (p *Page) readCommand() (string, error) {
	cmd := foreign.ColorInput(
		config.ColorCommand,
		fmt.Sprintf("command: '%s'", strings.Join(p.actions, "' '")),
	)

	cmd = strings.TrimSpace(cmd)

	match := make([]string, 0)

	for _, a := range p.actions {
		if strings.HasPrefix(a, cmd) {
			match = append(match, a)
		}
	}

	if len(match) == 0 {
		return "", fmt.Errorf("unknown command: %s", cmd)
	}

	if len(match) > 1 {
		return "", fmt.Errorf("match multiple actions: '%s'", strings.Join(match, "' '"))
	}

	return match[0], nil
}

func (p *Page) cmdBack() error {
	p.next = true
	pretty.ColorPrintln(
		config.ColorInfo,
		"going back to the previous page...",
	)
	return nil
}

func (p *Page) cmdExit() error {
	pretty.ColorPrintln(
		config.ColorInfo,
		"thank you for using Futures Wizards!!!",
	)
	os.Exit(0)
	return nil
}
