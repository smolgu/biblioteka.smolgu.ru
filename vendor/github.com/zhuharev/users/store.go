package users

import (
	"fmt"
	"github.com/zhuharev/users/config"
	"sync"
)

type Store interface {
	Name() string
	Connect(string) error

	Save(*User) error
	Get(int64) (*User, error)
	Read(int, ...int) ([]*User, error)
	Delete(int64) error

	GetByUserName(string) (*User, error)
	GetByEmail(string) (*User, error)

	Count() (int64, error)

	SetLogFile(string) error
}

var (
	stores = map[string]Store{}
	smu    sync.Mutex
)

var (
	ErrNotFound              = fmt.Errorf("not found")
	ErrUsernameAlreadyExists = fmt.Errorf("user already exists")
)

func Register(name string, driver Store) {
	smu.Lock()
	defer smu.Unlock()
	stores[name] = driver
}

func NewStore(cnf *config.Config) (Store, error) {
	var (
		name    = cnf.Database.Driver
		setting = cnf.Database.Setting
	)
	if driver, ok := stores[name]; !ok {
		panic(fmt.Sprintf("driver %s not found, %s", name, stores))
	} else {
		e := driver.Connect(setting)
		if e != nil {
			return nil, e
		}
		e = driver.SetLogFile(cnf.App.LogFile)
		if e != nil {
			return nil, e
		}
		return driver, nil
	}
	panic("unrechable")
}
