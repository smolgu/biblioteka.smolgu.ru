package albums

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "image/gif"
	_ "image/png"

	"github.com/fatih/color"
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/bloblog"
	"github.com/smolgu/lib/modules/middleware"
)

func Upload(c *middleware.Context) {

	var (
		err     string
		albumId = c.ParamsInt64(":id")
	)

	c.Req.ParseMultipartForm(32 << 20)
	file, _, e := c.Req.FormFile("file")
	if e != nil {
		fmt.Println(e)
		return
	}
	defer file.Close()
	//fmt.Fprintf(w, "%v", handler.Header)

	bts, e := ioutil.ReadAll(file)
	if e != nil {
		fmt.Println(e)
		err = e.Error()
	}
	rdr := bytes.NewReader(bts)

	img, _, e := image.Decode(rdr)
	if e != nil {
		color.Red("%s", e)
		return
	}

	buf := bytes.NewBuffer(nil)
	e = jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	if e != nil {
		color.Red("%s", e)
		return
	}

	id, e := bloblog.Insert(buf.Bytes())
	if e != nil {
		fmt.Println(e)
		err = e.Error()
	}

	photo, e := models.PhotoCreate(albumId, id)

	c.JSON(200, map[string]interface{}{
		"success": fmt.Sprint(id),
		"error":   err,
		"err":     e,
		"photo":   photo,
	})
}
