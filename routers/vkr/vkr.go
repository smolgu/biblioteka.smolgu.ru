package vkr

import (
	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func Index(c *middleware.Context) {
	c.Data["facs"] = models.SmolGUFaculties
	c.HTML(200, "vkr/index")
}

func AddVkr(c *middleware.Context) {
	if c.Req.Method == "POST" {
		// 32MiB
		e := c.Req.ParseMultipartForm(32 << 20)
		if e != nil {
			color.Red("%s", e)
			return
		}

		file, _, e := c.Req.FormFile("file")
		if e != nil {
			color.Red("%s", e)
			return
		}

		book, e := models.CreateBookFromReader(c.Query("title"), file)
		if e != nil {
			color.Red("%s", e)
			return
		}
		e = file.Close()
		if e != nil {
			color.Red("%s", e)
			return
		}

		vkr := new(models.Vkr)
		vkr.Title = c.Query("title")
		vkr.FirstName = c.Query("first_name")
		vkr.LastName = c.Query("last_name")
		vkr.Patronymic = c.Query("patronymic")
		vkr.BookId = book.Id
		vkr.TrainingDirectionId = c.QueryInt64("training_direction")
		vkr.Level = c.Query("level")
		vkr.UploaderId = c.User.Id
		vkr.Year = c.QueryInt("year")

		e = models.Save(vkr)
		if e != nil {
			color.Red("%s", e)
		}
	}
	dirs, e := models.TrainingDirectionAllList()
	if e != nil {
		color.Red("%s", e)
	}
	c.Data["dirs"] = dirs
	c.Data["faculties"] = models.SmolGUFaculties
	c.HTML(200, "vkr/add")
}

func DirectionYears(c *middleware.Context) {
	var (
		tdId = c.ParamsInt64(":tdId")
	)

	years, e := models.TrainingDirectionYears(tdId)
	if e != nil {
		color.Red("%s", e)
	}

	c.Data["tdId"] = tdId
	c.Data["fac"] = c.ParamsInt64(":fac")
	c.Data["years"] = years
	c.HTML(200, "vkr/years")
}

func DirectionYearList(c *middleware.Context) {
	var (
		tdId = c.ParamsInt64(":tdId")
		year = c.ParamsInt(":year")
	)

	vkrs, e := models.DirectionYearList(tdId, year)
	if e != nil {
		color.Red("%s", e)
	}

	c.Data["vkrs"] = vkrs

	c.HTML(200, "vkr/year")
}

func Download(c *middleware.Context) {
	var (
		id = c.ParamsInt64(":id")
	)
	//c.Resp.Header().Set("Content-type", "image/png")
	fname, e := models.VkrDownload(id)
	if e != nil {
		color.Red("%s", e)
	}
	c.ServeFile(fname)
}
