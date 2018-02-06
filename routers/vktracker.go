package routers

import "github.com/smolgu/lib/modules/middleware"
import "github.com/smolgu/lib/modules/vktracker"

func VkTrackerLog(c *middleware.Context) {
	recs, err := vktracker.Reports(10)
	if err != nil {
		c.Error(500, err.Error())
		return
	}
	c.Data["recs"] = recs
	c.HTML(200, "vktracker/log")
}
