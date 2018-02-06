package models

import (
	"github.com/zhuharev/users"
)

type Status users.Status

const (
	Student Status = 2 << (iota + 10)
	Guru
	Stuff
	Librarian
	ElectronicResources
	ChiefElectronicResources
	ViceDirector
	Director
)

var Statuses = map[Status]string{
	Student:                  "Студент",
	Guru:                     "Преподаватель",
	Stuff:                    "Сотрудник",
	Librarian:                "Библиотекарь",
	ElectronicResources:      "Сотрудник ОЭР",
	ChiefElectronicResources: "Начальник ОЭР",
	ViceDirector:             "Заместитель директора",
	Director:                 "Директор",
}

func (st Status) String() (s string) {
	switch users.Status(st) {
	case users.Admin:
		return "Администратор"
	}
	if r, ok := Statuses[st]; ok {
		return r
	} else {
		return "Гость"
	}
}

func (s Status) Add(toAdd Status) Status {
	return s | toAdd
}

func (s Status) Remove(toRemove Status) Status {
	if s&toRemove == 0 {
		return s
	}
	return s ^ toRemove
}
