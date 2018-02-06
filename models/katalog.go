package models

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/sisteamnik/ruslanparser"
	"time"
)

type Katalog struct {
	Id       int64
	SourceId string `xorm:"index unique"`

	Title  string
	Author string

	BookId int64

	Updated time.Time
	Created time.Time
	Deleted time.Time
}

func NewKatalogFromRuslanBook(b ruslanparser.Book) (*Katalog, error) {
	var (
		k = new(Katalog)
	)

	e := copier.Copy(k, &b)
	if e != nil {
		return k, e
	}

	// clean id after copier
	k.Id = 0

	return k, nil
}

func KatalogSave(k *Katalog) error {
	return Save(k)
}

func KatalogGet(id int64) (k *Katalog, e error) {
	k = new(Katalog)
	_, e = x.Id(id).Get(k)
	return k, e
}

func KatalogGetByOldId(id string) (k *Katalog, e error) {
	k = new(Katalog)
	has, e := x.Where("? = source_id ", id).Get(k)
	if !has {
		return nil, fmt.Errorf("not found")
	}
	return k, e
}
