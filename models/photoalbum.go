package models

import (
	"time"
)

type Album struct {
	Id    int64
	Title string
	Text  string

	Cat int64

	Photos []*Photo `xorm:"-"`

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func AlbumCreate(title string, text string, cat int64) (*Album, error) {
	album := new(Album)
	album.Title = title
	album.Text = text
	album.Cat = cat

	e := Save(album)

	return album, e
}

func AlbumGet(id int64) (*Album, error) {
	var (
		e      error
		album  = new(Album)
		photos []*Photo
	)
	_, e = x.Id(id).Get(album)
	if e != nil {
		return nil, e
	}
	e = x.Where("album_id = ?", id).Find(&photos)
	if e != nil {
		return nil, e
	}
	album.Photos = photos
	return album, e
}

func AlbumList(cat int, limit int, offsets ...int) (res []*Album, err error) {
	if cat == 0 {
		err = x.Limit(limit, offsets...).OrderBy("id desc").Find(&res)
	} else {
		err = x.Where("cat = ?", cat).Limit(limit, offsets...).OrderBy("id desc").Find(&res)
	}
	return res, err
}

type Photo struct {
	Id      int64
	AlbumId int64

	BlobId int64

	Caption string

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

func PhotoCreate(albumId int64, blobId int64) (*Photo, error) {
	photo := new(Photo)
	photo.AlbumId = albumId
	photo.BlobId = blobId

	e := Save(photo)

	return photo, e
}
