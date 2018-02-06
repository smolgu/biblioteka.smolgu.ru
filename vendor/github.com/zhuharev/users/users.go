package users

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id   int64  `json:"id"`
	Name string `xorm:"unique index"`

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`

	Phone string
	Email string `xorm:"unique index" json:"email"`

	Status Status `json:"status"`

	HashedPassword []byte `json:"-"`

	Data UserData `json:"data"`

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted time.Time `xorm:"deleted" json:"-"`
}

type UserData map[string]interface{}

func (u *User) ValidatePassword(password string) bool {
	e := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	if e != nil {
		return false
	}
	return true
}

func (u *User) SetPassword(password string) error {
	hashedPassword, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return e
	}
	u.HashedPassword = hashedPassword
	return nil
}
