package banner

import (
	"image"
	_ "image/gif"
	_ "image/png"
	"path/filepath"

	"github.com/dchest/uniuri"
	"github.com/disintegration/imaging"
	"github.com/fatih/color"
	"github.com/nfnt/resize"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
	"github.com/smolgu/lib/modules/setting"
)

var (
	bannerWidth = 172
)

func Add(c *middleware.Context) {
	if c.Req.Method == "POST" {
		name, e := uploadBanner(c)
		if e != nil {
			color.Red("%s", e)
			c.Flash.ErrorMsg = e.Error()
			c.Redirect("/admin/banners")
			return
		}
		banner := new(models.Banner)
		banner.Title = c.Query("title")
		banner.ImgPath = name
		banner.Url = c.Query("url")
		banner.Enabled = true
		banner.Color = c.Query("color")

		e = models.Save(banner)
		if e != nil {
			c.InternalServerError(e)
			return
		}
		c.Redirect("/admin/banners")
		return
	}

	c.HTML(200, "banner/add")
}

func Hide(c *middleware.Context) {
	var (
		id = c.ParamsInt64(":id")
	)
	banner, e := models.BannerGet(id)
	if e != nil {
		color.Red("%s", e)
		c.InternalServerError(e)
		return
	}
	models.Delete(banner)
	c.Redirect("/admin/banners")
}

func uploadBanner(c *middleware.Context) (string, error) {
	e := c.Req.ParseMultipartForm(32 << 20)
	if e != nil {
		return "", e
	}

	file, _, e := c.Req.FormFile("file")
	if e != nil {
		return "", e
	}
	defer file.Close()

	img, _, e := image.Decode(file)
	if e != nil {
		return "", e
	}

	rimg := resize.Resize(172, 0, img, resize.Bilinear)

	fname := uniuri.New() + ".jpg"
	name := filepath.Join(setting.DataDir, "uploads/banners", fname)

	e = imaging.Save(rimg, name)
	return "/uploads/banners/" + fname, e
}

func Img(c *middleware.Context) {

}

func Upload(c *middleware.Context) {

}

func Show(c *middleware.Context) {

}

func List(c *middleware.Context) {
	c.HTML(200, "banner/list")
}
