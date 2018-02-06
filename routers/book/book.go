package book

import (
	"github.com/Unknwon/paginater"
	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Upload(c *middleware.Context) {
	// 32MiB
	e := c.Req.ParseMultipartForm(32 << 20)
	if e != nil {
		color.Red("%s", e)
		return
	}

	file, _, e := c.Req.FormFile("file")
	if e != nil {
		color.Red("%s", e)
		return
	}
	defer file.Close()

	book, e := models.CreateBookFromReader(c.Query("title"), file)
	if e != nil {
		color.Red("%s", e)
	}

	c.JSON(200, map[string]interface{}{
		"book": book,
	})
}

func Add(c *middleware.Context) {
	c.HTML(200, "book/add")
}

func Show(c *middleware.Context) {
	b, e := models.BookGet(c.ParamsInt64(":id"))
	if e != nil {
		color.Red("%s", e)
	}
	var (
		p = c.QueryInt("p")
	)
	if p < 1 {
		p = 1
	}
	c.Data["paginater"] = paginater.New(b.Pages, 1, p, 1)
	c.Data["book"] = b
	c.Data["readMode"] = c.Query("read")
	c.HTML(200, "book/show")
}

func Page(c *middleware.Context) {
	c.Resp.Header().Set("Content-type", "image/png")
	var (
		id = c.QueryInt64("id")
		p  = c.QueryInt("p")
	)

	e := models.WriteBookPage(c.Resp, id, p)
	if e != nil {
		color.Red("%s", e)
	}

}
