package xorm

import (
	"github.com/zhuharev/users"
)

func (s *Store) GetByUserName(username string) (*users.User, error) {
	u := new(users.User)
	u.Name = username
	has, e := s.eng.Get(u)
	if e != nil {
		return nil, e
	}
	if !has {
		return nil, users.ErrNotFound
	}
	return u, nil
}

func (s *Store) GetByEmail(email string) (*users.User, error) {
	u := new(users.User)
	u.Email = email
	has, e := s.eng.Get(u)
	if e != nil {
		return nil, e
	}
	if !has {
		return nil, users.ErrNotFound
	}
	return u, nil
}
