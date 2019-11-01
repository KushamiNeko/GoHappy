package utils

import (
	"testing"
)

func TestKeyValuePair(t *testing.T) {
	inputs := "d=20200101; s=rty"
	pair, err := keyValuePair(inputs)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	d, ok := pair["d"]
	if !ok {
		t.Errorf("d should contain date")
	}

	if d != "20200101" {
		t.Errorf("d should be 20200101 but get %s", d)
	}

	s, ok := pair["s"]
	if !ok {
		t.Errorf("s should contain symbol")
	}

	if s != "rty" {
		t.Errorf("s should be rty but get %s", s)
	}

	inputs = " d=20200101;s=rty; "
	pair, err = keyValuePair(inputs)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	d, ok = pair["d"]
	if !ok {
		t.Errorf("d should contain date")
	}

	if d != "20200101" {
		t.Errorf("d should be 20200101 but get %s", d)
	}

	s, ok = pair["s"]
	if !ok {
		t.Errorf("s should contain symbol")
	}

	if s != "rty" {
		t.Errorf("s should be rty but get %s", s)
	}
}

func TestInputAbbreviation(t *testing.T) {
	a := map[string]string{
		"s": "symbol",
		"p": "price",
	}

	i := map[string]string{
		"s": "rty",
		"p": "100.25",
		"d": "20190521",
	}

	v := InputsAbbreviation(i, a)
	s, ok := v["symbol"]
	if !ok {
		t.Errorf("v should contain symbol")
	}

	if s != "rty" {
		t.Errorf("s should be rty")
	}

	p, ok := v["price"]
	if !ok {
		t.Errorf("v should contain price")
	}

	if p != "100.25" {
		t.Errorf("p should be 100.25")
	}

	d, ok := v["d"]
	if !ok {
		t.Errorf("v should contain date")
	}

	if d != "20190521" {
		t.Errorf("d should be 20190521")
	}
}
