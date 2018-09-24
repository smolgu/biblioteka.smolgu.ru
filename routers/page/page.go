package page

import (
	"fmt"
	"log"
	"time"

	"github.com/dchest/uniuri"
	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
	"github.com/smolgu/lib/modules/vk"
	"github.com/zhuharev/vkutil"
)

func New(c *middleware.Context) {

	var (
		id       = c.ParamsInt64(":id")
		bucketID = c.QueryInt64("bucket_id")

		isNewPage = id == 0
	)

	if id != 0 {
		p, e := models.GetPage(id)
		if e != nil {
			color.Red("%s", e)
		}
		c.Data["page"] = p
	} else {
		c.Data["page"] = models.Page{Date: time.Now()}
	}

	buckets, err := models.BucketList(50)
	if err != nil {
		color.Red("%s", err)
	}
	c.Data["buckets"] = buckets

	if c.Req.Method == "POST" {
		p := new(models.Page)
		p.Id = id
		p.Title = c.QueryTrim("title")
		p.Body = c.QueryTrim("body")
		p.OldId = uniuri.New()
		p.BucketID = bucketID
		if c.QueryTrim("old_id") != "" {
			p.OldId = c.QueryTrim("old_id")
		}
		p.Category = c.QueryInt64("category")
		p.CreatedBy = c.User.Id
		p.Date = time.Now()
		if c.QueryTrim("date") != "" {
			t, e := time.Parse("2006-01-02 15:04", c.QueryTrim("date"))
			if e == nil {
				p.Date = t
			}
		}
		p.Slug = c.QueryTrim("slug")

		for k, v := range c.Req.Form {
			if k == "images[]" {
				p.Images = v
				break
			}
		}
		if len(p.Images) > 0 {
			p.Image = p.Images[0]
			p.Images = p.Images[1:]
		}

		e := models.SavePage(p)
		if e != nil {
			color.Red("%s", e)
		}

		if isNewPage && p.Slug != "" {
			err := models.RouterSave(p.Id, models.ObjectPage, p.Slug)
			if e != nil {
				color.Red("%s", err)
				c.Error(200, err.Error())
				return
			}
		}

		if c.QueryBool("post_in_vk") && isNewPage {
			go func() {
				var (
					photo vkutil.Photo
					err   error
				)
				if p.Image != "" {
					photo, err = vk.Upload(p.Image)
					if err != nil {
						log.Println(err)
						return
					}
				}

				_, err = vk.Post(p.SanitazedBody(), photo)
				if err != nil {
					log.Println(err)
				}
			}()
		}
		c.Redirect("/news/" + fmt.Sprint(p.Id))
	}
	c.HTML(200, "page/new")
}

func Edit(c *middleware.Context) {
	New(c)
}

func Restore(c *middleware.Context) {
	id := c.ParamsInt64(":id")
	e := models.RestorePage(id)
	if e != nil {
		color.Red("%s", e)
		c.Flash.Error("Ошибка")
	} else {
		c.Flash.Info("Страница востановлена")
	}
	c.Redirect("/")
}

func Delete(c *middleware.Context) {
	id := c.ParamsInt64(":id")
	p, e := models.GetPage(id)
	if e != nil {
		color.Red("%s", e)
		return
	}
	e = models.DeletePage(id)
	if e != nil {
		color.Red("%s", e)
		return
	}
	c.Flash.Info(fmt.Sprintf("Страница <b>%s</b> удалена.  <b><a href='/news/%d/restore'>Отменить</a></b>", p.Title, p.Id))
	c.Redirect("/")
}

func Get(c *middleware.Context) {
	id := c.ParamsInt64(":id")
	p, e := models.GetPage(id)
	if e != nil {
		color.Red("%s", e)
		return
	}
	c.Data["page"] = p
	e = models.ViewPage(id, c.Session.ID())
	if e != nil {
		color.Red("%s", e)
	}
	if p.BucketID != 0 {
		bucket, err := models.BucketGet(p.BucketID)
		if err != nil {
			color.Red("%s", err)
		}
		c.Data["bucket"] = bucket
	}

	c.HTML(200, "page/get")
}
