package models

import (
	"time"
)

type Banner struct {
	Id      int64
	ImgPath string
	Url     string
	Title   string

	Color string

	Enabled bool

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted time.Time `xorm:"deleted"`
}

func GetEnabledBanners() ([]*Banner, error) {
	var (
		res []*Banner
	)
	e := x.Where("enabled = ?", true).Find(&res)
	return res, e
}

func BannerGet(id int64) (*Banner, error) {
	var (
		res = new(Banner)
	)
	_, e := x.Id(id).Get(res)
	return res, e
}
