package tag

import (
	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Show(c *middleware.Context) {
	pages, e := models.GetPagesByTag(c.Params(":tagName"))
	if e != nil {
		color.Red("%s", e)
	}
	c.Data["pages"] = pages
	c.HTML(200, "tags/show")
}
