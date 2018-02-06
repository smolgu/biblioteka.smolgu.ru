package albums

import (
	"github.com/Unknwon/com"
	"github.com/smolgu/lib/modules/bloblog"
	"github.com/smolgu/lib/modules/middleware"
	"strings"
)

func GetBlob(c *middleware.Context) {
	strId := c.Params(":id")
	strId = strings.TrimSuffix(strId, ".jpg")
	id := com.StrTo(strId).MustInt64()
	//fmt.Println(id)
	bts, e := bloblog.Get(id)
	if e != nil {
		panic(e)
	}
	//fmt.Println(len(bts))
	c.Resp.Header().Set("Content-Type", "image/jpeg")
	c.Resp.WriteHeader(200)
	c.Resp.Write(bts)
}
