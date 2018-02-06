package routers

import (
	"strings"

	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Route(c *middleware.Context) {
	color.Green("%s", c.Req.URL.Path)
	p := strings.TrimPrefix(c.Req.URL.Path, "/")
	objID, typ, err := models.RouterResolve(p)
	if err != nil {
		color.Yellow("Object not found %s ", err)
		c.NotFound()
		return
	}

	switch typ {
	case models.ObjectPage:
		routePage(c, objID)
		return
	default:
		c.NotFound()
		return
	}

}

func routePage(c *middleware.Context, pageID int64) {
	page, err := models.GetPage(pageID)
	if err != nil {
		c.Error(200, err.Error())
		return
	}
	c.Data["page"] = page
	err = models.ViewPage(pageID, c.Session.ID())
	if err != nil {
		color.Red("%s", err)
		c.Error(200, err.Error())
		return
	}
	c.HTML(200, "page/get")
}
