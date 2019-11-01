package config

import (
	"github.com/KushamiNeko/go_fun/utils/pretty"
)

const (
	ColorInfo     = pretty.PaperAmber300
	ColorCommand  = pretty.PaperTeal300
	ColorPages    = pretty.PaperPurple300
	ColorWarnings = pretty.PaperRed400
	ColorWhite    = pretty.PaperGrey300

	FloatDecimals  = 4
	DollarDecimals = 4

	PerContractCommissionFee = 1.5

	IdLen = 16

	FmtWidth    = 10
	FmtWidthM   = 6
	FmtWidthL   = 14
	FmtWidthXL  = 18
	FmtWidthXXL = 20

	TimeFormat = "20060102"

	DbAdmin        = "admin"
	DbTradingBooks = "trading_books"
	DbLiveTrading  = "live_trading"
	DbPaperTrading = "paper_trading"
	DbWatchList    = "watch_list"

	ColUser = "user"

	TestDb  = "test"
	TestCol = "test"
)
