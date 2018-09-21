package buckets

import (
	"fmt"
	"io/ioutil"

	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"

	// image packages for decoding
	_ "image/gif"
	_ "image/png"
	// another image decoders
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

func Upload(c *middleware.Context) {
	var (
		err      string
		bucketId = c.ParamsInt64(":id")
	)

	c.Req.ParseMultipartForm(32 << 20)
	file, info, e := c.Req.FormFile("file")
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

	uploaded, e := models.Upload(info.Filename, bts, models.FileCreateOptions{
		Title:    info.Filename,
		BucketID: bucketId,
	})
	if e != nil {
		err = e.Error()
		fmt.Println(err)
	}

	c.JSON(200, map[string]interface{}{
		"success": fmt.Sprint(uploaded.BlobID),
		"error":   err,
		"err":     e,
		"file":    uploaded,
	})
}
