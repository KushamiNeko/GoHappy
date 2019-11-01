package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type JsonDB struct {
	testing bool
	root    string
}

func NewJsonDB(testing bool) *JsonDB {

	const (
		dbRoot     string = "/home/neko/Documents/database/json/futures_wizards"
		testDbRoot string = "/home/neko/Documents/database/json/testing"
	)

	j := new(JsonDB)
	j.testing = testing

	if testing {
		j.root = testDbRoot
	} else {
		j.root = dbRoot
	}

	return j
}

func (j *JsonDB) dbPath(db, col string) string {
	fn := fmt.Sprintf("%s_%s.json", db, col)
	p := filepath.Join(j.root, fn)

	return p
}

func (j *JsonDB) read(db, col string) ([]map[string]string, error) {
	r := make([]map[string]string, 0)

	p := j.dbPath(db, col)
	content, err := ioutil.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		} else {
			return nil, err
		}
	}

	if len(content) == 0 {
		return r, nil
	}

	err = json.Unmarshal(content, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (j *JsonDB) write(db, col string, entities []map[string]string) error {
	p := j.dbPath(db, col)
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(entities, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (j *JsonDB) Insert(db, col string, entities ...map[string]string) error {

	r, err := j.read(db, col)
	if err != nil {
		return err
	}

	r = append(r, entities...)

	err = j.write(db, col, r)
	if err != nil {
		return err
	}

	return nil
}

func (j *JsonDB) Replace(
	db, col string,
	query map[string]string,
	entity map[string]string) error {

	if query == nil || len(query) == 0 {
		return fmt.Errorf("invalid query")
	}

	if len(entity) == 0 {
		return fmt.Errorf("invalid entity")
	}

	r, err := j.read(db, col)
	if err != nil {
		return err
	}

	updated := false

	for i, e := range r {
		found := true
		for k, v := range query {
			if val, ok := e[k]; !(ok && val == v) {
				found = false
			}
		}

		if found {
			for k, v := range entity {
				r[i][k] = v
			}

			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("no entity match query: %v", query)
	}

	err = j.write(db, col, r)
	if err != nil {
		return err
	}

	return nil
}

func (j *JsonDB) Find(
	db, col string,
	query map[string]string) ([]map[string]string, error) {

	r, err := j.read(db, col)
	if err != nil {
		return nil, err
	}

	if query == nil || len(query) == 0 {
		return r, nil
	}

	n := make([]map[string]string, 0)

	for i, e := range r {
		found := true
		for k, v := range query {
			if val, ok := e[k]; !(ok && val == v) {
				found = false
			}
		}

		if found {
			n = append(n, r[i])
		}
	}

	return n, nil
}

func (j *JsonDB) Delete(
	db, col string,
	query map[string]string) error {

	if query == nil || len(query) == 0 {
		return fmt.Errorf("invalid query or entity")
	}

	r, err := j.read(db, col)
	if err != nil {
		return err
	}

	n := make([]map[string]string, 0, len(r)-1)

	for _, e := range r {
		found := true
		for k, v := range query {
			if val, ok := e[k]; !(ok && val == v) {
				found = false
			}
		}

		if found {
			continue
		} else {
			n = append(n, e)
		}
	}

	err = j.write(db, col, n)
	if err != nil {
		return err
	}

	return nil
}

func (j *JsonDB) DropCol(db, col string) error {
	p := j.dbPath(db, col)

	err := os.Remove(p)
	if err != nil {
		return err
	}

	return nil
}
