package link

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-macaron/i18n"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
	"strings"
)

func List(c *middleware.Context) {
	var (
		limit  = c.QueryInt("limit")
		offset = c.QueryInt("offset")
	)
	links, e := models.LinkFind(limit, offset)
	if e != nil {
		color.Red("%s", e)
	}
	c.Data["links"] = links
	c.HTML(200, "link/list")
}

func Tags(c *middleware.Context, loc i18n.Locale) {
	var (
		tagName = c.Query("tag")
		limit   = c.QueryInt("limit")
		offset  = c.QueryInt("offset")
	)
	links, e := models.LinkFindByTag(tagName, offset, limit)
	if e != nil {
		color.Green("%s", e)
	}

	c.Data["links"] = links
	c.Data["Title"] = loc.Tr(tagName)
	c.HTML(200, "link/list")
}

func Edit(c *middleware.Context, l models.Link) {
	if c.Req.Method == "POST" {
		l.Id = c.ParamsInt64(":id")

		tags := strings.Split(c.Query("link_tags"), ",")
		for i, v := range tags {
			tags[i] = strings.TrimSpace(v)
		}
		l.Tags = tags

		if l.Title != "" && l.URL != "" {
			e := models.LinkSave(&l)
			if e != nil {
				color.Red("%s", e)
			}
			c.Redirect("/links/edit/" + fmt.Sprint(l.Id))
		}
	}
	lk, e := models.LinkGet(c.ParamsInt64(":id"))
	if e != nil {
		color.Red("%s", e)
	}
	c.Data["link"] = lk
	c.HTML(200, "link/edit")
}

func Delete(c *middleware.Context) {
	models.DeleteLink(c.ParamsInt64(":id"))
	c.Redirect("/links")
}

func Batch(c *middleware.Context) {
	if c.Req.Method == "POST" {
		body := c.Query("body")
		lines := strings.Split(body, "\n")

		for _, v := range lines {
			if arr := strings.Fields(v); len(arr) < 3 {
				continue
			} else {
				l := models.NewLink(arr[len(arr)-2], strings.Join(arr[0:len(arr)-2], " "))
				l.Tags = strings.Split(arr[len(arr)-1], ",")
				e := models.LinkSave(l)
				if e != nil {
					color.Red("%s", e)
				}
			}
		}

		c.Redirect("/links/batch")
	}

	c.HTML(200, "link/batch")

}
