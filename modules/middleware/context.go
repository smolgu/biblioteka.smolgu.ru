package middleware

import (
	"github.com/fatih/color"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	"github.com/hiteshmodha/goDevice"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/setting"
	"gopkg.in/macaron.v1"
	"time"
)

type Context struct {
	*macaron.Context
	Cache   cache.Cache
	Flash   *session.Flash
	Session session.Store

	User        *models.User
	IsSigned    bool
	IsBasicAuth bool
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, cache cache.Cache, sess session.Store, f *session.Flash) {
		ctx := &Context{
			Context: c,
			Cache:   cache,
			Flash:   f,
			Session: sess,
		}
		// Compute current URL for real-time change language.
		//ctx.Data["Link"] = setting.AppSubUrl + ctx.Req.URL.Path

		ctx.Data["PageStartTime"] = time.Now()

		// Get user from session if logined.
		userInterface := ctx.Session.Get("user")
		if u, ok := userInterface.(*models.User); ok {
			ctx.IsSigned = true
			ctx.User = u
			ctx.Data["IsSigned"] = ctx.IsSigned
			ctx.Data["SignedUser"] = u
			ctx.Data["SignedUserName"] = u.Name
			ctx.Data["IsAdmin"] = u.IsAdmin()
		} else {
			ctx.Data["SignedUserName"] = ""
		}

		banners, e := models.GetEnabledBanners()
		if e != nil {
			color.Red("%s", e)
		}
		c.Data["banners"] = banners

		// If request sends files, parse them here otherwise the Query() can't be parsed and the CsrfToken will be invalid.
		/*if ctx.Req.Method == "POST" && strings.Contains(ctx.Req.Header.Get("Content-Type"), "multipart/form-data") {
			if err := ctx.Req.ParseMultipartForm(setting.AttachmentMaxSize << 20); err != nil && !strings.Contains(err.Error(), "EOF") { // 32MB max size
				ctx.Handle(500, "ParseMultipartForm", err)
				return
			}
		}*/

		ctx.Data["MainMenu"] = setting.MainMenu
		menus, e := models.MenusGet()
		if e != nil {
			color.Red("%s", e)
		}
		ctx.Data["Menus"] = menus
		ctx.Data["Faculties"] = models.SmolGUFaculties
		ctx.Data["TrainingDirections"] = models.Trains
		ctx.Data["CurrentUrl"] = ctx.Req.URL.String()

		c.Map(ctx)
	}
}

func (c *Context) HTML(status int, name string, data ...interface{}) {
	dt := goDevice.GetType(c.Req.Request)
	if dt == goDevice.MOBILE {
		c.Data["IsMobile"] = true
		mobileDisabled := c.GetCookie("mobile_disabled")
		if mobileDisabled == "" {
			c.HTMLSet(status, "mobile", name, data...)
			return
		}
	}
	c.Context.HTML(status, name, data...)
}
