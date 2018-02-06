package xorm

import (
	"os"

	"github.com/go-xorm/xorm"
	"github.com/zhuharev/users"
)

type Store struct {
	Driver string
	eng    *xorm.Engine
}

func (s *Store) Name() string {
	return "xorm"
}

func (s *Store) Connect(setting string) error {
	if s.Driver == "" {
		panic("not driver setted")
	}
	eng, e := xorm.NewEngine(s.Driver, setting)
	if e != nil {
		return e
	}
	eng.SetLogger(nil)
	e = eng.Sync2(&users.User{})
	if e != nil {
		return e
	}

	s.eng = eng
	return nil
}

func (s *Store) SetLogFile(path string) error {
	if path == "" {
		return nil
	}
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	s.eng.SetLogger(xorm.NewSimpleLogger(f))
	s.eng.ShowSQL()
	return nil
}

func (s *Store) Save(u *users.User) error {
	if u.Id == 0 {
		_, e := s.eng.Insert(u)
		if e != nil {
			return e
		}
	} else {
		_, e := s.eng.Id(u.Id).Update(u)
		if e != nil {
			return e
		}
	}
	return nil
}

func (s *Store) Get(id int64) (*users.User, error) {
	u := new(users.User)
	has, e := s.eng.Id(id).Get(u)
	if e != nil {
		return nil, e
	}
	if !has {
		return nil, users.ErrNotFound
	}
	return u, nil
}

func (s *Store) Read(limit int, offset ...int) ([]*users.User, error) {
	var us []*users.User
	e := s.eng.Limit(limit, offset...).Find(&us)
	return us, e
}

func (s *Store) Delete(id int64) error {
	_, e := s.eng.Id(id).Delete(new(users.User))
	return e
}

func (s *Store) Count() (int64, error) {
	return s.eng.Count(new(users.User))
}

func init() {
	users.Register("xorm", &Store{})
	users.Register("sqlite3", &Store{Driver: "sqlite3"})
}
