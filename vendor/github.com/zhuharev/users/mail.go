package users

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"time"
)

var (
	ErrInvalidEmail = fmt.Errorf("Invalid email")
)

type EmailConfirmation struct {
	Email string
	TS    time.Time
	Code  string

	srv *Service
}

func NewConfirmation(email string, srv *Service) *EmailConfirmation {
	cf := new(EmailConfirmation)
	cf.Email = email
	cf.TS = time.Now()
	cf.Code = RandString(4)
	cf.srv = srv
	return cf
}

func (c *EmailConfirmation) Hash() string {
	return SecretHash(c.Email+c.Code, c.srv.Config.App.Secret)
}

func (c *EmailConfirmation) Message() string {
	messagefmt := `
<a href="http://%s/confirm?email=%s&code=%s&token=%s">Подтвердить e-mail</a><br>
Ваш код подтверждения <b>%s</b>
	`
	message := fmt.Sprintf(messagefmt, c.srv.Config.Web.Host, c.Email, c.Code, c.Hash(), c.Code)
	return message
}

func (s *Service) NewConfirmation(email string) *EmailConfirmation {
	return NewConfirmation(email, s)
}

func (s *Service) SendConfirmEmail(u *User) (*EmailConfirmation, error) {
	if u == nil {
		return nil, fmt.Errorf("user is nil")
	}
	if u.Email == "" {
		return nil, ErrInvalidEmail
	}
	cf := s.NewConfirmation(u.Email)
	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.Config.Mail.Sender, "Library")
	m.SetAddressHeader("To", u.Email, "lib")
	m.SetHeader("Subject", "Email confirmation")
	m.SetBody("text/html", cf.Message())

	d := gomail.NewPlainDialer(s.Config.Mail.Server,
		s.Config.Mail.Port, s.Config.Mail.Sender, s.Config.Mail.Password)

	if e := d.DialAndSend(m); e != nil {
		return nil, e
	}

	return cf, nil
}

func (s *Service) ConfirmEmail(email, code, hash string) error {
	str := email + code
	token := SecretHash(str, s.Config.App.Secret)
	if token != hash {
		return fmt.Errorf("bad request")
	}
	u, e := s.Store.GetByEmail(email)
	if e != nil {
		return e
	}
	u.Status = u.Status.Add(EmailConfirmed)
	s.Store.Save(u)
	return nil
}
