package model

import (
	"fmt"
)

type User struct {
	name string
	uid  string
}

func NewUser(name, uid string) *User {

	u := new(User)
	u.name = name
	u.uid = uid

	return u
}

func NewUserFromEntity(e map[string]string) (*User, error) {
	name, ok := e["name"]
	if !ok {
		return nil, fmt.Errorf("no name in entity %v", e)
	}

	uid, ok := e["uid"]
	if !ok {
		return nil, fmt.Errorf("no uid in entity %v", e)
	}

	u := new(User)
	u.name = name
	u.uid = uid

	return u, nil
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Uid() string {
	return u.uid
}

func (u *User) Entity() map[string]string {
	return map[string]string{
		"name": u.name,
		"uid":  u.uid,
	}
}
