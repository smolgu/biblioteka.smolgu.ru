package buckets

import (
	"fmt"
	"log"

	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

// Create handler
func Create(c *middleware.Context) {
	if c.Req.Method == "POST" {
		bucket, e := models.BucketCreate(c.Query("title"))
		if e != nil {
			fmt.Println(e)
			c.Error(500, e.Error())
			return
		}
		c.Redirect("/buckets/" + fmt.Sprint(bucket.Id))
		return
	}
	c.HTML(200, "buckets/create")
}

func Get(c *middleware.Context) {
	id := c.ParamsInt64(":id")
	bucket, e := models.BucketGet(id)
	if e != nil {
		log.Printf("err get bucket by id bucket=%d", id)
		c.Error(500, e.Error())
		return
	}
	c.Data["bucket"] = bucket
	c.Data["moreScripts"] = []string{"dropzone.js"}
	c.HTML(200, "buckets/get")
}

//
// func List(c *middleware.Context) {
// 	var (
// 		p           = c.QueryInt("p")
// 		itemsInPage = 10
// 	)
// 	if p != 0 {
// 		p--
// 	}
// 	albums, err := models.AlbumList(1, itemsInPage, p*itemsInPage)
// 	if err != nil {
// 		log.Println(err)
// 		c.Error(500, err.Error())
// 		return
// 	}
// 	c.Data["albums"] = albums
// 	c.HTML(200, "albums/list")
// }
