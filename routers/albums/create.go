package albums

import (
	"fmt"
	"log"

	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Create(c *middleware.Context) {
	if c.Req.Method == "POST" {
		album, e := models.AlbumCreate(c.Query("title"), c.Query("text"), c.QueryInt64("cat"))
		if e != nil {
			fmt.Println(e)
			c.Error(500, e.Error())
			return
		}
		c.Redirect("/albums/" + fmt.Sprint(album.Id))
		return
	}
	c.HTML(200, "albums/create")

}

func Get(c *middleware.Context) {
	album, e := models.AlbumGet(c.ParamsInt64(":id"))
	if e != nil {
		c.Error(500, e.Error())
		return
	}
	c.Data["album"] = album
	c.Data["moreScripts"] = []string{"dropzone.js"}
	c.HTML(200, "albums/get")
}

func List(c *middleware.Context) {
	var (
		p           = c.QueryInt("p")
		itemsInPage = 10
	)
	if p != 0 {
		p--
	}
	albums, err := models.AlbumList(1, itemsInPage, p*itemsInPage)
	if err != nil {
		log.Println(err)
		c.Error(500, err.Error())
		return
	}
	c.Data["albums"] = albums
	c.HTML(200, "albums/list")
}
