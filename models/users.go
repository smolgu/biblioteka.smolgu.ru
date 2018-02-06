package models

import (
	"encoding/gob"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/smolgu/lib/modules/setting"
	"github.com/zhuharev/users"
	"github.com/zhuharev/users/config"
	_ "github.com/zhuharev/users/store/xorm"
)

var (
	UserService *users.Service
)

func init() {
	gob.Register(new(User))
}

func NewUserService() error {
	cnf := new(config.Config)
	cnf.Admin.Login = setting.Users.AdminLogin
	cnf.Admin.Password = setting.Users.AdminPassword
	cnf.App.Secret = setting.Users.Secret
	cnf.Database.Driver = setting.Users.DbDriver
	cnf.Database.KVPath = setting.Users.KvPath
	cnf.Database.Setting = setting.Users.DbSetting
	cnf.Web.Host = setting.HostName

	log.Println(cnf.Database)
	s, err := users.NewFromConfig(cnf)
	if err != nil {
		return errors.Wrap(err, "from config")
	}
	err = s.Store.SetLogFile(filepath.Join(setting.LogDir, "users.log"))
	if err != nil {
		return errors.Wrap(err, "set log file")
	}
	UserService = s
	return nil
}

type User struct {
	*users.User
	Status Status
}

func (u User) IsAdmin() bool {
	return u.Status&Status(users.Admin) != 0 || u.Status&ElectronicResources != 0 || u.Status&Director != 0
}

func assignUser(usersUser *users.User) *User {
	u := new(User)
	u.User = usersUser

	u.Status = Status(usersUser.Status)

	return u
}

func CreateUser(username, password string) (*User, error) {
	us, e := UserService.CreateUser(username, password)
	if e != nil {
		return nil, e
	}
	return assignUser(us), nil
}

func GetByUserName(username string) (*User, error) {
	us, e := UserService.Store.GetByUserName(username)
	if e != nil {
		return nil, e
	}
	return assignUser(us), nil
}

func GetUser(id int64) (*User, error) {
	us, e := UserService.Store.Get(id)
	if e != nil {
		return nil, e
	}
	return assignUser(us), nil
}

func SaveUser(u *User) error {
	u.User.Status = users.Status(u.Status)
	return UserService.Store.Save(u.User)
}

func GetUserList(offset int, limit int) (u []*User, e error) {
	var us []*users.User
	us, e = UserService.Store.Read(offset, limit)
	if e != nil {
		return
	}
	for _, v := range us {
		u = append(u, assignUser(v))
	}
	return
}
