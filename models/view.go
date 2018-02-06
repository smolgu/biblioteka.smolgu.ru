package models

import (
	"bytes"
	"fmt"

	"github.com/willf/bloom"
)

var (
	kvMax  uint = 100
	kvKeys uint = 3
)

func ViewPage(id int64, sessionId string) (e error) {

	var (
		kvKey = fmt.Sprintf("bloom_%d", id)
	)

	data, e := Get(kvKey)
	if e != nil {
		return e
	}
	bf := bloom.New(kvMax, kvKeys)
	if len(data) != 0 {
		_, e = bf.ReadFrom(bytes.NewReader(data))
		if e != nil {
			return e
		}
	}
	if voted := bf.TestAndAddString(sessionId); !voted {
		_, e = x.Exec("update page set views=(select ifnull(views,0) from page where id = $1)+1 where id = $1;", id)
		var buf = bytes.NewBuffer(nil)
		_, e := bf.WriteTo(buf)
		if e != nil {
			return e
		}
		e = Set(kvKey, buf.Bytes())
		if e != nil {
			return e
		}
	}
	return
}
