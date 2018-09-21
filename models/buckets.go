package models

import (
	"time"
)

type Bucket struct {
	ID    int64
	Title string

	Files []*File `xorm:"-"`

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func BucketCreate(title string) (*Bucket, error) {
	buck := new(Bucket)
	buck.Title = title

	e := Save(buck)

	return buck, e
}

func BucketGet(id int64) (*Bucket, error) {
	var (
		e     error
		buck  = new(Bucket)
		files []*File
	)
	_, e = x.Id(id).Get(buck)
	if e != nil {
		return nil, e
	}
	e = x.Where("bucket_id = ?", id).Find(&files)
	if e != nil {
		return nil, e
	}
	buck.Files = files
	return buck, e
}

func BucketList(limit int, offsets ...int) (res []*Bucket, err error) {
	err = x.Limit(limit, offsets...).OrderBy("id desc").Find(&res)
	return res, err
}
