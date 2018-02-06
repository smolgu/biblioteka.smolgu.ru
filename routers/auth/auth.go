package auth

import (
	"fmt"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func MustSigned(ctx *middleware.Context) {
	if !ctx.IsSigned {
		ctx.Flash.Error("Вы должны войти")
		ctx.Redirect("/account/auth")
	}
}

func MustAdmin(ctx *middleware.Context) {
	if !ctx.IsSigned {
		ctx.Flash.Error("Вы должны войти")
		ctx.Redirect("/account/auth")
		return
	}
	if !ctx.User.IsAdmin() {
		ctx.Flash.Error("Не хватает прав")
		ctx.Redirect("/account/dashboard")
	}
}

func Auth(ctx *middleware.Context) {
	if ctx.IsSigned {
		ctx.Redirect("/account/dashboard")
		return
	}
	if ctx.Req.Method == "POST" {

		// shortcut
		g := func(name string) string {
			return ctx.Req.FormValue(name)
		}

		email := g("email")
		password := g("password")

		u, e := models.GetByUserName(email)
		if e != nil {
			ctx.Flash.Error(e.Error())
			ctx.Redirect(ctx.Req.RequestURI)
			return
		}

		if !u.ValidatePassword(password) {
			e = fmt.Errorf("Password or login not found")
			ctx.Flash.Error(e.Error())
			ctx.Redirect(ctx.Req.RequestURI)
			return
		}

		e = ctx.Session.Set("user", u)
		if e != nil {
			ctx.Flash.Error(e.Error())
			ctx.Redirect(ctx.Req.RequestURI)
			return
		}

		fmt.Println("Authorised")
		ctx.Redirect("/account/dashboard")
		return
	}
	ctx.HTML(200, "auth")
}

func Logout(c *middleware.Context) {
	c.Session.Delete("user")
	c.SetCookie("s", "")
	c.Redirect("/account/auth")
}
