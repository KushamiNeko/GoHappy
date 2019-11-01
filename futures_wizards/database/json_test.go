package database

import (
	"testing"

	"github.com/KushamiNeko/futures_wizards/config"
)

func TestInsert(t *testing.T) {
	db := NewJsonDB(true)
	e := map[string]string{
		"a": "b",
	}
	err := db.Insert(config.TestDb, config.TestCol, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	db.DropCol(config.TestDb, config.TestCol)
}

func TestReplace(t *testing.T) {
	db := NewJsonDB(true)
	e := map[string]string{
		"a": "b",
		"b": "b",
	}
	err := db.Insert(config.TestDb, config.TestCol, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	q := map[string]string{
		"a": "b",
	}

	e = map[string]string{
		"a": "c",
	}

	err = db.Replace(config.TestDb, config.TestCol, q, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	n, err := db.Find(config.TestDb, config.TestCol, e)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if len(n) != 1 {
		t.Errorf("len should be 1 but get %d", len(n))
	}

	if n[0]["a"] != "c" {
		t.Errorf("new value should be c, but get %s", n[0]["a"])
	}

	if n[0]["b"] != "b" {
		t.Errorf("new value should be b, but get %s", n[0]["b"])
	}

	db.DropCol(config.TestDb, config.TestCol)
}

func TestDelete(t *testing.T) {
	db := NewJsonDB(true)
	es := []map[string]string{
		map[string]string{
			"a": "b",
			"b": "b",
		},
		map[string]string{
			"c": "d",
			"b": "b",
		},
	}
	err := db.Insert(config.TestDb, config.TestCol, es...)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	q := map[string]string{
		"b": "b",
	}

	rs, err := db.Find(config.TestDb, config.TestCol, q)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if len(rs) != 2 {
		t.Errorf("len should be 2 but get %d", len(rs))
	}

	dq := map[string]string{
		"a": "b",
	}

	err = db.Delete(config.TestDb, config.TestCol, dq)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	rs, err = db.Find(config.TestDb, config.TestCol, q)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if len(rs) != 1 {
		t.Errorf("len should be 1 but get %d", len(rs))
	}

	db.DropCol(config.TestDb, config.TestCol)
}
