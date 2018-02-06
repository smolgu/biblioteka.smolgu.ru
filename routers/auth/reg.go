package auth

import (
	"fmt"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Reg(ctx *middleware.Context) {
	if ctx.IsSigned {
		ctx.Redirect("/account/dashboard")
		return
	}

	ctx.Data["Faculties"] = models.SmolGUFaculties

	if ctx.Req.Method == "POST" {

		// shortcut
		g := func(name string) string {
			return ctx.Req.FormValue(name)
		}

		email := g("email")
		fmt.Println(ctx.Req.Form)
		password := g("password")

		u, e := models.CreateUser(email, password)
		if e != nil {
			ctx.Flash.Error(e.Error())
			fmt.Println(e)
			ctx.Redirect("/account/reg")
			return
		}

		u.Email = email
		u.FirstName = g("first_name")
		u.LastName = g("last_name")
		u.Patronymic = g("patronymic")

		u.Data["fac"] = g("fac")

		e = models.SaveUser(u)
		if e != nil {
			ctx.Flash.Error(e.Error())
			fmt.Println(e)
			ctx.Redirect("/account/reg")
			return
		}

		ctx.Flash.Success("Вы зарегистрированы, можете войти.")
		ctx.Redirect("/account/auth")
		return
	}

	ctx.HTML(200, "reg")
}

/*func Restore(c *middleware.Context) {
	if ctx.IsSigned {
		ctx.Redirect("/account/dashboard")
		return
	}
	if c.Req.Method == "POST" {
		email := c.Query("email")

		u, e := models.GetByUserName(email)
		if e != nil {
			ctx.Flash.Error(e.Error())
			ctx.Redirect(ctx.Req.RequestURI)
			return
		}


	}

	c.HTML(200, "restore")
}
*/
