package context

import (
	"fmt"

	"github.com/KushamiNeko/futures_wizards/config"
	"github.com/KushamiNeko/futures_wizards/model"
	"github.com/KushamiNeko/futures_wizards/utils"
)

type database interface {
	Find(db, col string, q map[string]string) ([]map[string]string, error)

	Insert(db, col string, es ...map[string]string) error

	Replace(db, col string,
		q map[string]string,
		e map[string]string) error

	Delete(db, col string,
		q map[string]string) error

	DropCol(db, col string) error
}

type Context struct {
	db   database
	user *model.User
}

func NewContext(database database) *Context {
	c := new(Context)
	c.db = database

	return c

}

func (c *Context) Db() database {
	return c.db
}

func (c *Context) User() *model.User {
	return c.user
}

func (c *Context) Login(name string) error {

	es, err := c.db.Find(
		config.DbAdmin,
		config.ColUser,
		map[string]string{
			"name": name,
		},
	)
	if err != nil {
		return err
	}

	if len(es) == 0 {
		return fmt.Errorf("invalid user")
	} else {
		if len(es) != 1 {
			panic(fmt.Sprintf("user should be unique: %v", es))
		}
	}

	c.user, err = model.NewUserFromEntity(es[0])

	return nil
}

func (c *Context) NewUser(name string) error {

	es, err := c.db.Find(
		config.DbAdmin,
		config.ColUser,
		map[string]string{
			"name": name,
		},
	)
	if err != nil {
		return err
	}

	if len(es) == 0 {
		user := model.NewUser(name, utils.RandString(config.IdLen))

		c.db.Insert(
			config.DbAdmin,
			config.ColUser,
			user.Entity(),
		)

		c.user = user
	} else {
		return fmt.Errorf("user has already exist: %s", name)
	}

	return nil
}
