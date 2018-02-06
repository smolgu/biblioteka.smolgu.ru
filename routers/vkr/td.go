package vkr

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func AddDirection(c *middleware.Context) {

	if c.Req.Method == "POST" {
		td := new(models.TrainingDirection)
		td.Title = c.Query("title")
		td.SubTitle = c.Query("subtitle")
		td.Code = c.Query("code")
		td.FacultyId = c.QueryInt64("faculty")
		td.CreatedBy = c.User.Id

		e := models.Save(td)
		if e != nil {
			color.Red("%s", e)
		}
		c.Redirect("/td/" + fmt.Sprint(td.FacultyId))
		return
	}

	c.Data["faculties"] = models.SmolGUFaculties
	c.HTML(200, "td/add")
}

func FacultyDirections(c *middleware.Context) {
	var (
		facId = c.ParamsInt64(":fac")
	)

	dirs, e := models.TrainingDirectionFacultyList(facId)
	if e != nil {
		color.Red("%s", e)
	}

	c.Data["dirs"] = dirs
	c.HTML(200, "td/list")
}

func DeleteDirection(c *middleware.Context) {
	var (
		tdId = c.ParamsInt64(":id")
	)

	td, e := models.TrainingDirectionById(tdId)
	if e != nil {
		color.Red("%s", e)
		return
	}

	e = models.Delete(td)
	if e != nil {
		color.Red("%s", e)
		return
	}
	c.Redirect("/td/" + c.Params(":fac"))
}
