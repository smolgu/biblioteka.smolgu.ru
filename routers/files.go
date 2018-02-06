package routers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/bloblog"
	"github.com/smolgu/lib/modules/middleware"
)

// UploadFile upload an file
func UploadFile(c *middleware.Context) {
	if c.Req.Method == "POST" {
		// 32MiB
		e := c.Req.ParseMultipartForm(32 << 30)
		if e != nil {
			color.Red("%s", e)
			return
		}

		file, info, e := c.Req.FormFile("file")
		if e != nil {
			color.Red("%s", e)
			return
		}
		defer file.Close()

		bts, err := ioutil.ReadAll(file)
		if err != nil {
			c.Error(200, err.Error())
			return
		}

		f, err := models.Upload(info.Filename, bts)
		if err != nil {
			c.Error(200, err.Error())
			return
		}

		c.Data["uploaded"] = f
		c.Redirect("/files/upload/new?uploaded=" + fmt.Sprint(f.BlobID))
		return
	}

	if id := c.QueryInt("uploaded"); id != 0 {
		c.Data["uploaded"] = id

	}
	c.HTML(200, "files/new")
}

func Blob(c *middleware.Context) {
	id := c.ParamsInt64(":id")
	bts, err := bloblog.Get(id)
	if err != nil {
		c.Error(500, err.Error())
		return
	}
	ct := http.DetectContentType(bts[:1024])
	c.Resp.Header().Set("Content-Type", ct)
	c.Resp.WriteHeader(200)
	c.Resp.Write(bts)
	c.Resp.Flush()
}
