package model

import "fmt"

type ContractSpecs struct {
	ts   map[string]string
	ibkr map[string]string

	unit map[string]float64
}

var (
	specs *ContractSpecs
)

func init() {

	specs = new(ContractSpecs)

	specs.ts = map[string]string{
		"es":  "es",
		"ym":  "ym",
		"nq":  "nq",
		"rty": "rty",
		"nk":  "nyd",
		"us":  "zb",
		"ty":  "zn",
		"fv":  "zf",
		"tu":  "zt",
		"ed":  "ge",
		"ec":  "6e",
		"jy":  "6j",
		"bp":  "6b",
		"ad":  "6a",
		"cd":  "6c",
		"sf":  "6s",
		"mp1": "6m",
		"ne1": "6n",
		"gc":  "gc",
		"si":  "si",
		"hg":  "hg",
		"pl":  "pl",
		"cl":  "cl",
		"ng":  "ng",
		"rb":  "rb",
		"ho":  "ho",
		"s":   "zs",
		"w":   "zw",
		"c":   "zc",
		"bo":  "zl",
		"sm":  "zm",
		"lc":  "le",
		"lh":  "he",
		"fc":  "gf",
	}

	specs.ibkr = map[string]string{}

	specs.unit = map[string]float64{
		"es":  50.0,
		"ym":  5.0,
		"nq":  20.0,
		"rty": 50.0,
		"nyd": 5.0,
		"zb":  1000.0,
		"zn":  1000.0,
		"zf":  1000.0,
		"zt":  1000.0,
		"ge":  2500.0,
		"6e":  125000.0,
		"6j":  12500000.0,
		"6a":  100000.0,
		"6b":  62500.0,
		"6c":  100000.0,
		"6s":  125000.0,
		"6n":  100000.0,
		"6m":  500000.0,
		"gc":  100.0,
		"si":  5000.0,
		"hg":  25000.0,
		"pl":  50.0,
		"cl":  1000.0,
		"rb":  42000.0,
		"ng":  10000.0,
		"ho":  42000.0,
		"zs":  5000.0,
		"zc":  5000.0,
		"zw":  5000.0,
		"zl":  60000.0,
		"zm":  100.0,
		"le":  40000.0,
		"he":  40000.0,
		"gf":  50000.0,
	}

}

func NewContractSpecs() *ContractSpecs {
	return specs
}

func (c *ContractSpecs) lookupSymbol(symbol string) (string, error) {
	cme, ok := c.ts[symbol]
	if !ok {
		return "", fmt.Errorf("invalid TradeStation symbol: %s", symbol)
	}

	return cme, nil
}

func (c *ContractSpecs) LookupContractUnit(symbol string) (float64, error) {
	cme, ok := c.ts[symbol]
	if !ok {
		return 0, fmt.Errorf("invalid TradeStation symbol: %s", symbol)
	}

	unit, ok := c.unit[cme]
	if !ok {
		panic(fmt.Sprintf("%s should exist in unit but does not", cme))
	}

	return unit, nil
}

func (c *ContractSpecs) ValidateSymbol(symbol string) bool {
	_, err := c.lookupSymbol(symbol)
	if err == nil {
		return true
	} else {
		return false
	}

}
