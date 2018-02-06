package bloblog

import (
	"path/filepath"

	"github.com/smolgu/lib/modules/setting"
	"github.com/zhuharev/bloblog"
	//"gopkg.in/macaron.v1"
	//"fmt"
	//"io/ioutil"
)

var (
	bl     *bloblog.BlobLog
	inited bool
)

func initBlobl() {
	if !inited {
		NewContext()
		inited = true
	}
}

func NewContext() {
	var e error
	bl, e = bloblog.Open(filepath.Join(setting.DataDir, "photos.bloblog"))
	if e != nil {
		panic(e)
	}
}

func Insert(data []byte) (int64, error) {
	initBlobl()
	return bl.Insert(data)
}

func Get(id int64) ([]byte, error) {
	initBlobl()
	return bl.Get(id)
}
