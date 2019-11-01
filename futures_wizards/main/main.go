package main

import (
	"os"
	"strings"

	"github.com/KushamiNeko/futures_wizards/config"
	"github.com/KushamiNeko/futures_wizards/context"
	"github.com/KushamiNeko/futures_wizards/database"
	"github.com/KushamiNeko/futures_wizards/view"
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

func main() {
	db := database.NewJsonDB(false)
	ctx := context.NewContext(db)

	var err error

	var entry *view.Entry

	pretty.ColorPrintln(config.ColorInfo, "Welcome to Futures Wizards Wealthy Interface")

	user := pretty.ColorInput(config.ColorWhite, "user name:")
	user = strings.TrimSpace(user)

	err = ctx.Login(user)
	if err != nil {
		goto exit
	}

	entry = view.NewEntry(ctx)
	entry.Main()

	return

exit:
	pretty.ColorPrintln(config.ColorWarnings, err.Error())
	os.Exit(1)

}
