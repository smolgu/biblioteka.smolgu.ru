package users

import (
	"github.com/zhuharev/users/config"
)

type Service struct {
	Store  Store
	Config *config.Config
	KV     *KV
}

func New(filepath string) (*Service, error) {
	cnf, e := config.NewFromFile(filepath)
	if e != nil {
		return nil, e
	}
	return NewFromConfig(cnf)
}

func NewFromConfig(cnf *config.Config) (*Service, error) {
	s := new(Service)
	s.Config = cnf

	store, e := NewStore(cnf)
	if e != nil {
		return nil, e
	}
	s.Store = store

	kv, e := NewKV(cnf.Database.KVPath)
	if e != nil {
		return nil, e
	}
	s.KV = kv

	e = s.checkInstall()
	if e != nil {
		return nil, e
	}

	return s, nil
}

func (s *Service) IsExistUserName(username string) (bool, error) {
	u, e := s.Store.GetByUserName(username)
	if e == ErrNotFound || u == nil {
		return false, nil
	}
	if e != nil {
		return false, e
	}
	return true, nil
}

func (s *Service) CreateUser(username, password string) (*User, error) {
	if ok, e := s.IsExistUserName(username); ok {
		return nil, ErrUsernameAlreadyExists
	} else {
		if e != nil {
			return nil, e
		}
	}
	u := new(User)
	u.Name = username
	u.Email = username
	u.Data = map[string]interface{}{}
	e := u.SetPassword(password)
	if e != nil {
		return nil, e
	}
	e = s.Store.Save(u)
	return u, e
}

func (s *Service) MakeAdminPermission(uid int64) error {
	u, e := s.Store.Get(uid)
	if e != nil {
		return e
	}
	u.Status = u.Status.Add(Admin)
	e = s.Store.Save(u)
	if e != nil {
		return e
	}
	return nil
}

func (s *Service) DeleteAdminPermission(uid int64) error {
	u, e := s.Store.Get(uid)
	if e != nil {
		return e
	}
	u.Status = u.Status.Remove(Admin)
	e = s.Store.Save(u)
	if e != nil {
		return e
	}
	return nil
}

func (s *Service) checkInstall() error {
	count, e := s.Store.Count()
	if e != nil {
		return e
	}
	if count == 0 {
		u, e := s.CreateUser(s.Config.Admin.Login, s.Config.Admin.Password)
		if e != nil {
			return e
		}
		return s.MakeAdminPermission(u.Id)
	}
	return nil
}
