package models

import (
	"path/filepath"
	"time"

	"github.com/smolgu/lib/modules/bloblog"
)

type File struct {
	Id       int64
	Type     int
	Title    string
	BlobID   int64 `xorm:"blob_id"`
	BucketID int64 `xorm:"bucket_id"`

	OriginalFileName string
	Mime             string

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted time.Time `xorm:"deleted"`
}

func GetFile(id int64) (*File, error) {
	f := new(File)
	_, err := x.Id(id).Get(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type FileCreateOptions struct {
	BucketID int64
	Title    string
}

func Upload(fileName string, data []byte, opts ...FileCreateOptions) (*File, error) {
	var (
		bucketID int64
		title    string
	)
	if len(opts) > 0 {
		bucketID = opts[0].BucketID
		title = opts[0].Title
	}

	id, err := bloblog.Insert(data)
	if err != nil {
		return nil, err
	}

	f := &File{
		Type:             1, //upload
		BlobID:           id,
		OriginalFileName: fileName,
		Mime:             GetMime(fileName),
		BucketID:         bucketID,
		Title:            title,
	}
	_, err = x.Insert(f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetMime(fName string) string {
	switch filepath.Ext(fName) {
	case ".pdf":
		return "application/pdf"
	case ".jpg":
		return "image/jpeg"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	return "text/plain"
}
