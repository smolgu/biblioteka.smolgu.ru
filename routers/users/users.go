package users

import (
	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Users(c *middleware.Context) {
	var (
		offset = c.QueryInt("offset")
		limit  = c.QueryInt("limit")
	)

	if limit < 10 {
		limit = 10
	}

	users, e := models.GetUserList(limit, offset)
	if e != nil {
		color.Red("%s", e)
	}
	c.Data["users"] = users
	c.HTML(200, "users")
}

func Edit(c *middleware.Context) {
	var (
		id = c.ParamsInt64(":id")
		u  *models.User
		e  error
	)

	if id != 0 {
		u, e = models.GetUser(id)
		if e != nil {
			color.Red("%s", e)
		}
		c.Data["user"] = u
	}

	if c.Req.Method == "POST" {
		u.FirstName = c.Query("first_name")
		u.LastName = c.Query("last_name")
		u.Patronymic = c.Query("patronymic")

		u.Email = c.Query("email")
		color.Green("current status %s", u.Status)
		u.Status = models.Status(c.QueryInt64("user_status"))
		color.Green("status changed to %s", u.Status)

		/*data*/
		u.Data["Faculty"] = c.Query("fac")
		u.Data["TrainingDirection"] = c.Query("train")

		pass := c.Query("password")
		if pass != "" {
			e := u.SetPassword(pass)
			if e != nil {
				c.Flash.Error(e.Error())
				c.Redirect(c.Req.RequestURI)
				return
			}
		}

		e := models.SaveUser(u)
		if e != nil {
			c.Flash.Error(e.Error())
			c.Redirect(c.Req.RequestURI)
			return
		}

		c.Flash.Success("Данные сохранены")
		c.Redirect(c.Req.RequestURI)
	}
	c.Data["statuses"] = models.Statuses
	c.HTML(200, "users/edit")
}
