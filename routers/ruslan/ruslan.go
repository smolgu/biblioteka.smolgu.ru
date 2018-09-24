package ruslan

//
// import (
// 	"github.com/fatih/color"
// 	"github.com/sisteamnik/ruslanparser"
// 	"github.com/smolgu/lib/models"
// 	"github.com/smolgu/lib/modules/middleware"
// 	"github.com/smolgu/lib/modules/ruslan"
// )
//
// func Search(c *middleware.Context) {
// 	if c.Query("TERM_1") != "" {
// 		e := c.Req.ParseForm()
// 		if e != nil {
// 			color.Red("%s", e)
// 		}
// 		vals := c.Req.Form
// 		records, e := ruslan.Search(vals)
// 		if e != nil {
// 			color.Red("%s", e)
// 		}
// 		c.Data["records"] = records
// 	}
//
// 	c.HTML(200, "ruslan/search_form")
// }
//
// func Book(c *middleware.Context) {
//
// 	var (
// 		id = c.Query("id")
// 		b  ruslanparser.Book
// 		e  error
// 	)
//
// 	k, e := models.KatalogGetByOldId(id)
// 	if e != nil {
// 		color.Red("%s", e)
// 	}
//
// 	if k == nil {
// 		b, e = ruslan.GetById(id)
// 		if e != nil {
// 			color.Red("%s", e)
// 		}
//
// 		k, e = models.NewKatalogFromRuslanBook(b)
// 		if e != nil {
// 			color.Red("%s", e)
// 		}
//
// 		e = models.Save(k)
// 		if e != nil {
// 			color.Red("%s", e)
// 		}
// 	}
//
// 	c.Data["book"] = k
// 	c.HTML(200, "ruslan/book")
// }
