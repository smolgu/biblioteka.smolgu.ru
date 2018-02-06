package config

import (
	"fmt"
	"gopkg.in/gcfg.v1"
)

type Config struct {
	App struct {
		Secret  string
		LogFile string
	}
	Web struct {
		Host string
	}
	Database struct {
		Driver, Setting, KVPath string
	}
	Mail struct {
		Sender, Password string
		Server           string
		Port             int
	}
	Admin struct {
		Login    string
		Password string
	}
}

func NewFromFile(filepath string) (*Config, error) {
	if filepath == "" {
		return nil, fmt.Errorf("Empty config path")
	}
	cnf := new(Config)
	e := gcfg.ReadFileInto(cnf, filepath)
	return cnf, e
}
