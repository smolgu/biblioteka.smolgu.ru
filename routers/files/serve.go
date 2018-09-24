// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package files

import (
	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/bloblog"
	"github.com/smolgu/lib/modules/middleware"
)

func Serve(c *middleware.Context) {
	f, err := models.GetFile(c.ParamsInt64(":id"))
	if err != nil {
		return
	}
	bts, err := bloblog.Get(f.BlobID)
	if err != nil {
		c.Error(500, err.Error())
		return
	}
	c.Resp.Header().Set("Content-Type", f.Mime)
	c.Resp.WriteHeader(200)
	c.Resp.Write(bts)
	c.Resp.Flush()
}
