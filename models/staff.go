package models

import (
	"time"
)

type Staff struct {
	Id     int64
	UserId int64

	Order int

	JobFrom     string
	Education   string
	Position    string
	Departament int64

	Deleted time.Time
}

type Departament struct {
	Id       int64
	ParentId int64

	Title string

	Phone    string
	Location string
}
