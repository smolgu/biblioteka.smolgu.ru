package models

import (
	"path/filepath"
	"time"

	"github.com/smolgu/lib/modules/bloblog"
)

type File struct {
	Id     int64
	Type   int
	BlobID int64

	OriginalFileName string
	Mime             string

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted time.Time `xorm:"deleted"`
}

func Upload(fileName string, data []byte) (*File, error) {
	id, err := bloblog.Insert(data)
	if err != nil {
		return nil, err
	}

	f := &File{
		Type:             1, //upload
		BlobID:           id,
		OriginalFileName: fileName,
		Mime:             getMime(fileName),
	}
	_, err = x.Insert(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func getMime(fName string) string {
	switch filepath.Ext(fName) {
	case ".pdf":
		return "application/pdf"
	}
	return "text/plain"
}
