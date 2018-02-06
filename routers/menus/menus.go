package menus

import (
	"fmt"
	"strings"

	"github.com/Unknwon/com"

	"github.com/smolgu/lib/models"
	"github.com/smolgu/lib/modules/middleware"
)

func List(c *middleware.Context) {
	// todo handle error
	list, _ := models.MenuList()
	c.Data["list"] = list
	c.HTML(200, "menus/list")
}

func Create(c *middleware.Context) {
	if title := c.Query("title"); title != "" {
		models.MenuCreate(title, c.Query("slug"))
	}

	c.Redirect("/menus")
}

func Edit(c *middleware.Context) {
	// todo error
	menu, _ := models.MenuGet(c.QueryInt64("id"))
	c.Data["menu_for_edit"] = menu
	c.HTML(200, "menus/edit")
}

func ItemCreate(c *middleware.Context) {
	var (
		title  = c.Query("title")
		link   = c.Query("link")
		menuId = c.QueryInt64("id")
		parent = c.QueryInt64("parent")
	)
	if link != "" && title != "" && menuId != 0 {
		models.MenuItemCreate(title, link, menuId, parent)
	}

	c.Redirect("/menus/edit?id=" + fmt.Sprint(menuId))
}

func SetPosition(c *middleware.Context) {
	var (
		pos    []int64
		posStr = strings.Split(c.Query("positions"), ";")
	)

	fmt.Println(posStr)
	for _, v := range posStr {
		pos = append(pos, com.StrTo(v).MustInt64())
	}
	e := models.MenuSetPosition(pos)
	if e != nil {
		fmt.Println(e)
	}
	c.Redirect("/menus/edit?id=" + c.Query("id"))
}

func ItemEdit(c *middleware.Context) {
	var (
		id = c.QueryInt64("id")
	)
	mi, e := models.MenuItemGet(id)
	if e != nil {
		fmt.Println(e)
		c.Redirect("/menus")
		return
	}
	if c.Req.Method == "POST" {
		mi.Id = id
		mi.Title = c.Query("title")
		mi.Link = c.Query("link")
		e = models.Save(mi)
	}
	if e != nil {
		fmt.Println(e)
		c.Redirect("/menus")
		return
	}
	c.Data["menuitem"] = mi
	c.HTML(200, "menus/itemedit")

}
