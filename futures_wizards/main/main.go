package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/KushamiNeko/go_fun/utils/database"
	"github.com/KushamiNeko/go_fun/utils/foreign"
	"github.com/KushamiNeko/go_fun/utils/pretty"
	"github.com/KushamiNeko/go_happy/futures_wizards/config"
	"github.com/KushamiNeko/go_happy/futures_wizards/context"
	"github.com/KushamiNeko/go_happy/futures_wizards/view"
)

func main() {
	db := database.NewFileDB(
		filepath.Join(
			os.Getenv("HOME"),
			"Documents/database/yaml/futures_wizards"),
		database.YamlEngine,
	)

	ctx := context.NewContext(db)

	var err error

	var entry *view.Entry

	pretty.ColorPrintln(config.ColorInfo, "Welcome to Futures Wizards Wealthy Interface")

	user := foreign.ColorInput(config.ColorWhite, "user name:")
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
