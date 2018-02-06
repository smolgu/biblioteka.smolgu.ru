package routers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cznic/mathutil"
	"github.com/disintegration/imaging"
	"github.com/fatih/color"
	"github.com/nfnt/resize"
	"github.com/ungerik/go-dry"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"

	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
	"github.com/smolgu/lib/modules/search"
	"github.com/smolgu/lib/modules/setting"
	"github.com/smolgu/lib/modules/vktracker"
	"gopkg.in/macaron.v1"
)

func InitGlobal() {

	models.NewEngine()

	e := models.NewUserService()
	if e != nil {
		log.Fatalln(e)
	}

	models.NewKVContext()

	models.NewTagsContext()
	search.NewContext()

	vktracker.NewContext()
}

func Index(c *middleware.Context) {
	page := mathutil.Max(c.QueryInt("p"), 1)
	ps, e := models.GetNews(page)
	if e != nil {
		// TODO
		_ = e
	}
	c.Data["pages"] = ps
	c.Data["nextPage"] = page + 1
	c.HTML(200, "index")
}

func GetPageByOldId(c *middleware.Context) {
	route := strings.TrimSuffix(c.Req.RequestURI, "index.shtml")
	route = strings.Trim(route, "/")
	color.Cyan("%s", route)
	rp, e := models.GetPageByOldId(route)
	if e != nil {
		fmt.Println(e)
	}

	c.Data["page"] = rp
	c.HTML(200, "page/get")
}

func EditPageByOldId(c *middleware.Context) {
	route := strings.TrimSuffix(c.Req.RequestURI, "index.shtml")
	route = strings.Trim(route, "/")
	route = strings.TrimPrefix(route, "pages/edit/")
	color.Cyan("%s", "pages/"+route)
	rp, e := models.GetPageByOldId("pages/" + route)
	if e != nil {
		fmt.Println("Can't get page, ", e)
	}

	if c.Req.Method == "POST" {
		p := new(models.Page)
		p.Id = rp.Id
		p.Title = c.QueryTrim("title")
		p.Body = c.QueryTrim("body")
		//p.OldId = uniuri.New()
		p.Category = models.NewCategoryId
		p.CreatedBy = c.User.Id
		p.Date = time.Now()
		if c.QueryTrim("date") != "" {
			t, e := time.Parse("2006-01-02 15:04", c.QueryTrim("date"))
			if e == nil {
				p.Date = t
			}
		}

		tags := strings.Split(c.Query("tags"), ",")
		for i, v := range tags {
			tags[i] = strings.TrimSpace(v)
		}
		p.Tags = tags

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

		var (
			images []string
		)
		for _, v := range p.Images {
			if v != "" {
				images = append(images, v)
			}
		}
		p.Images = images

		e := models.SavePage(p)
		if e != nil {
			color.Red("%s", e)
		}
		rp = p
		c.Redirect(c.Req.RequestURI)
	}

	c.Data["page"] = rp
	c.HTML(200, "pages/edit")
}

func Upload(c *macaron.Context) {
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
	defer file.Close()

	img, _, e := image.Decode(file)
	if e != nil {
		color.Red("%s", e)
		return
	}

	uid, e := saveUploadedImage(img)

	c.JSON(200, map[string]interface{}{
		"hash":     uid,
		"img_path": filepath.Join("/uploads/", uid[:1], uid[1:2], uid+".jpg"),

		"img_small": filepath.Join("/uploads/", uid[:1], uid[1:2],
			uid+"_small.jpg"),
	})
}

func RandomString() string {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}

	return hex.EncodeToString(buffer)
}

func saveUploadedImage(img image.Image) (uid string, e error) {
	dir := setting.DataDir
	startTime := time.Now()
	i := resize.Resize(0, 100, img, resize.Bilinear)
	color.Green("Resized %s", time.Since(startTime))
	uid = RandomString()
	if len(uid) < 3 {
		e = fmt.Errorf("random broken")
		return
	}
	fpath := filepath.Join(dir, "uploads", uid[:1], uid[1:2])
	fname := filepath.Join(fpath, uid)
	os.MkdirAll(fpath, 0777)
	e = imaging.Save(img, fname+".jpg")
	if e != nil {
		color.Red("%s", e)
	}

	f, err := os.OpenFile(fname+"_small.jpg", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	e = jpeg.Encode(f, i, &jpeg.Options{65})
	if e != nil {
		color.Red("%s", e)
	}
	return
}

func Resize(c *macaron.Context) {
	s := strings.Split(c.Params("*"), "/")[0]
	urlArr := strings.SplitN(c.Params("*"), "/", 2)
	if len(urlArr) != 2 {
		color.Red("Not 2 %v", urlArr)
		return
	}

	opts := ""

	size := strings.SplitN(s, "x", 2)
	if w := size[0]; w != "" {
		opts = "width=" + w
	}
	if h := size[1]; h != "" {
		if opts == "" {
			opts = "height=" + h
		} else {
			opts += "&height=" + h
		}
	}

	ch := func(w http.ResponseWriter, r *http.Response, header string) {
		key := http.CanonicalHeaderKey(header)
		if value, ok := r.Header[key]; ok {
			w.Header()[key] = value
		}
	}

	c.Resp.Header().Set("Cache-Control", "max-age=31536000, public")
	c.Resp.Header().Set("Expires", time.Now().AddDate(10, 0, 0).Format(http.TimeFormat))

	u := fmt.Sprintf("http://a.zhuharev.ru/resize?%s&url=%s",
		opts, strings.SplitN(c.Params("*"), "/", 2)[1])

	resp, e := http.DefaultClient.Get(u)
	if e != nil {
		color.Red("[d] %s", e)
		return
	}
	defer resp.Body.Close()

	ch(c.Resp, resp, "Content-Length")
	ch(c.Resp, resp, "Content-Type")

	_, e = io.Copy(c.Resp, resp.Body)
	if e != nil {
		return
	}

}

func IndexPages(c *macaron.Context) {
	pages, e := models.GetAllPages()
	if e != nil {
		color.Red("%s", e)
		return
	}
	e = search.IndexPages(pages)
	if e != nil {
		color.Red("%s", e)
		return
	}
	c.JSON(200, "indexed")
}

func Search(c *macaron.Context) {
	q := c.Query("q")
	ids, e := search.Search(q)
	if e != nil {
		color.Red("%s", e)
		return
	}
	pages, e := models.GetPagesByIds(ids)
	c.Data["pages"] = pages
	c.Data["SearchQuery"] = q
	c.Data["Title"] = fmt.Sprintf("Результаты поиска по запросу: %s", q)
	dry.FileAppendString("data/db/search.log", c.Query("q")+"\n")
	c.HTML(200, "search/results")
}
