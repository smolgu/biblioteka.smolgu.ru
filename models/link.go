package models

import (
	"strings"
)

type Link struct {
	Id     int64
	Title  string
	URL    string `form:"url"`
	Target LinkTarget
	Tags   []string
	Img    string
}

type LinkTarget string

const (
	LBlank LinkTarget = "_blank"
	LNone  LinkTarget = "_none"
	LTop   LinkTarget = "_top"
)

func NewLink(href, title string) *Link {
	target := LNone
	if strings.HasPrefix(href, "http") {
		target = LBlank
	}
	return &Link{
		Title:  title,
		URL:    href,
		Target: target,
	}
}

func LinkFind(limit, offset int) (Links, error) {
	var (
		l Links
	)
	e := x.Limit(limit, offset).OrderBy("id desc").Find(&l)
	return l, e
}

func LinkGet(id int64) (*Link, error) {
	var (
		l = new(Link)
	)
	_, e := x.Id(id).Get(l)
	return l, e
}

func LinkFindByTag(tagName string, offset, limit int) (Links, error) {
	var (
		res Links
	)
	ids, e := GetTagItems(tagName, new(Link))
	e = x.In("id", ids).OrderBy("id desc").Find(&res)
	if e != nil {
		return nil, e
	}
	return res, nil
}

func LinkSave(l *Link) (e error) {
	if l.Id == 0 {
		_, e = x.Insert(l)
		if e != nil {
			return
		}
		e = SaveTags(l)
		return
	} else {
		_, e = x.Id(l.Id).Update(l)
		if e != nil {
			return
		}
		e = SaveTags(l)
		return
	}
}

func DeleteLink(id int64) (e error) {
	_, e = x.Id(id).Delete(new(Link))
	return
}

type Links []Link

func (l Links) ByTags(otags ...string) (ls Links) {
	if otags == nil {
		return nil
	}
	for _, link := range l {
		ok := true
		for _, tag := range link.Tags {
			if !strInArr(tag, otags) {
				ok = false
			}
		}
		if ok {
			ls = append(ls, link)
		}
	}
	return
}

func strInArr(in string, arr []string) bool {
	for _, v := range arr {
		if in == v {
			return true
		}
	}
	return false
}
