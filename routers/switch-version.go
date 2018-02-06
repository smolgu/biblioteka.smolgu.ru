package routers

import (
	"gopkg.in/macaron.v1"
)

func SwitchToMobile(c *macaron.Context) {
	c.SetCookie("mobile_disabled", "")
	c.Redirect(c.Query("next"), 302)
}

func SwitchToFull(c *macaron.Context) {
	c.SetCookie("mobile_disabled", "1")
	c.Redirect(c.Query("next"), 302)
}
