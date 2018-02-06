package models

import (
	"sort"
	"time"
)

type Menu struct {
	Id    int64
	Title string
	Slug  string

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`

	Items MenuItems `xorm:"-"`
}

func MenuGetAll() ([]*Menu, error) {
	var menus []*Menu
	e := x.Find(&menus)
	if e != nil {
		return nil, e
	}
	return menus, e
}

type MenuItem struct {
	Id    int64
	Title string
	Link  string

	MenuId   int64
	ParentId int64

	Position int

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

type MenuItems []MenuItem

func (a MenuItems) Len() int           { return len(a) }
func (a MenuItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MenuItems) Less(i, j int) bool { return a[i].Position < a[j].Position }

func (m MenuItems) FirstLevel() (res MenuItems) {
	for _, v := range m {
		if v.ParentId == 0 {
			res = append(res, v)
		}
	}
	sort.Sort(res)
	return
}

func (m MenuItems) Childs(id int64) (res MenuItems) {
	for _, v := range m {
		if v.ParentId == id {
			res = append(res, v)
		}
	}
	sort.Sort(res)
	return
}

func (m MenuItems) Menu(id int64) (items MenuItems) {
	for _, v := range m {
		if v.MenuId == id {
			items = append(items, v)
		}
	}
	return
}

func MenuCreate(title string, slug string) (*Menu, error) {
	m := &Menu{
		Title: title,
		Slug:  slug,
	}
	e := Save(m)
	return m, e
}

func MenuItemCreate(title, link string, menuId int64, parent int64) (*MenuItem, error) {
	m := &MenuItem{
		Title:    title,
		Link:     link,
		MenuId:   menuId,
		ParentId: parent,
	}
	e := Save(m)
	return m, e
}

func MenuItemGet(id int64) (*MenuItem, error) {
	mi := new(MenuItem)
	_, e := x.Id(id).Get(mi)
	return mi, e
}

func MenuSetPosition(positions []int64) error {
	for position, itemId := range positions {
		item := MenuItem{Position: position + 1}
		_, e := x.Id(itemId).Update(&item)
		if e != nil {
			return e
		}
	}
	return nil
}

func MenuList() ([]*Menu, error) {
	var (
		res []*Menu
	)
	e := x.Find(&res)
	return res, e
}

func MenuItemsGet(menuId int64) (MenuItems, error) {
	var res []MenuItem
	e := x.Where("menu_id = ?", menuId).Find(&res)
	return res, e
}

func MenuGet(id int64) (*Menu, error) {
	var (
		res = new(Menu)
	)
	_, e := x.Id(id).Get(res)
	if e != nil {
		return res, e
	}
	res.Items, e = MenuItemsGet(id)
	return res, e
}

type Menus struct {
	Menus []*Menu
	Items MenuItems
}

func (m *Menus) Menu(slug string) *Menu {
	for _, v := range m.Menus {
		if v.Slug == slug {
			if v.Items != nil {
				return v
			} else {
				v.Items = m.Items.Menu(v.Id)
				return v
			}
		}
	}
	return &Menu{Items: MenuItems{}}
}

func MenusGet() (*Menus, error) {
	var (
		items = MenuItems{}
		menus = new(Menus)
	)
	e := x.Find(&items)
	if e != nil {
		return menus, e
	}
	menus.Items = items

	allmenu, e := MenuGetAll()
	if e != nil {
		return menus, e
	}
	menus.Menus = allmenu

	return menus, nil
}
