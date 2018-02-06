package account

import (
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Dashboard(ctx *middleware.Context) {
	if ctx.Req.Method == "POST" {
		// shortcut
		g := func(name string) string {
			return ctx.Req.FormValue(name)
		}

		u := ctx.User

		u.FirstName = g("first_name")
		u.LastName = g("last_name")
		u.Patronymic = g("patronymic")

		u.Email = g("email")

		/*data*/
		u.Data["Faculty"] = g("fac")
		u.Data["TrainingDirection"] = g("train")

		pass := g("password")
		if pass != "" {
			e := u.SetPassword(pass)
			if e != nil {
				ctx.Flash.Error(e.Error())
				ctx.Redirect(ctx.Req.RequestURI)
				return
			}
		}

		e := models.SaveUser(u)
		if e != nil {
			ctx.Flash.Error(e.Error())
			ctx.Redirect(ctx.Req.RequestURI)
			return
		}

		ctx.Flash.Success("Данные сохранены")
		ctx.Redirect(ctx.Req.RequestURI)
	}
	ctx.HTML(200, "dashboard")
}
