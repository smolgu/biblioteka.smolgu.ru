package models

import (
	"log"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/smolgu/lib/modules/setting"
)

type Page struct {
	Id       int64
	OldId    string `xorm:"unique"`
	Title    string
	Body     string `xorm:"TEXT"`
	Category int64

	BucketID int64 `xorm:"bucket_id"`

	Slug string

	Image  string
	Images []string

	Published bool

	Date      time.Time
	CreatedBy int64

	Views int64

	Tags []string

	Deleted time.Time `xorm:"deleted"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

// Summary return cutted Body by 255 symbols
func (p Page) Summary() (summary string) {
	out, err := sanitize.HTMLAllowing(p.Body, []string{"p"})
	if err != nil {
		log.Println(err)
		return ""
	}
	if strings.Index(out, "\n") > 10 && strings.Index(out, "\n") < 255 {
		return strings.SplitN(out, "\n", 2)[0]
	}
	arr := strings.Split(out, " ")
	if len(arr) < 2 {
		return p.Body
	}
	var cnt = 0
	for i, v := range arr {
		cnt += len(v)
		if cnt > 255 {
			summary = strings.Join(arr[:i-1], " ")
			break
		}
	}
	return
}

// SanitazedBody return body text without html tags.
// Useful for external posting.
func (p Page) SanitazedBody() string {
	return sanitize.HTML(p.Body)
}

type Category struct {
	Id   int64
	Name string
}

func GetNews(pages ...int) ([]Page, error) {
	var (
		page = 1
	)

	if len(pages) > 0 {
		page = pages[0]
	}
	var ps []Page
	e := x.Where("category = ?", NewCategoryId).Limit(setting.ItemsInPage, (page-1)*setting.ItemsInPage).OrderBy("datetime(date) desc").Find(&ps)
	return ps, e
}

func GetPageByOldId(id string) (p *Page, e error) {
	p = new(Page)
	_, e = x.Where("? = old_id ", id).Get(p)
	return p, e
}

func GetPage(id int64) (p *Page, e error) {
	p = new(Page)
	_, e = x.Id(id).Get(p)
	return p, e
}

func DeletePage(id int64) error {
	p, e := GetPage(id)
	if e != nil {
		return e
	}
	e = RemoveAllTags(p)
	if e != nil {
		return e
	}
	_, e = x.Delete(p)
	return e
}

func RestorePage(id int64) error {
	_, e := x.Exec("update page set deleted = NULL where id = ?", id)
	return e
}

func SavePage(p *Page) (e error) {
	if p.Id == 0 {
		_, e = x.InsertOne(p)
		if e != nil {
			return
		}
		e = SaveTags(p)
		return
	} else {
		_, e = x.Id(p.Id).Cols("title", "body", "image", "images", "date", "slug").Update(p)
		if e != nil {
			return
		}
		e = SaveTags(p)
		return
	}
	return
}

func GetPagesByTag(tagName string, pages ...int) ([]*Page, error) {
	ids, e := GetTagItems(tagName, new(Page))
	if e != nil {
		return nil, e
	}
	return GetPagesByIds(ids, pages...)
}

func GetPageBySlug(slug string) (*Page, error) {
	p := new(Page)
	has, err := x.Where("slug = ?", slug).Get(p)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return p, nil
}

func GetPagesByIds(ids []int64, pages ...int) ([]*Page, error) {
	var (
		res  []*Page
		page = 1
	)
	if len(pages) > 0 {
		page = pages[0]
	}
	e := x.In("id", ids).Limit(setting.ItemsInPage, (page-1)*setting.ItemsInPage).OrderBy("datetime(date) desc").Find(&res)
	if e != nil {
		return nil, e
	}
	return res, nil
}

func GetAllPages() ([]*Page, error) {
	var (
		pages []*Page
	)
	e := x.Find(&pages)
	if e != nil {
		return pages, e
	}
	return pages, e
}
