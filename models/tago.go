package models

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	"github.com/smolgu/lib/modules/setting"
	"github.com/zhuharev/tago"
)

var (
	tags *tago.Tago
)

func NewTagsContext() {
	bdb, e := bolt.Open(filepath.Join(setting.DataDir, "bolt.db"), 0666, nil)
	if e != nil {
		log.Fatalf("tag db open err: %s", e)
	}
	tags, e = tago.NewWithBoltDb(bdb)
	if e != nil {
		log.Fatalf("tag db create err: %s", e)
	}
}

func objPrefix(bean interface{}) []byte {
	color.Green("obj prefix is %s", reflect.TypeOf(bean).Elem().Name())
	return []byte(reflect.TypeOf(bean).Elem().Name())
}

func tagsId(bean interface{}) ([]string, int64, error) {
	val := reflect.ValueOf(bean)
	val = val.Elem()
	idVal := val.FieldByName("Id")
	if !idVal.IsValid() {
		return nil, 0, fmt.Errorf("id value is nil")
	}
	id := idVal.Interface().(int64)
	tagsVal := val.FieldByName("Tags")
	if !tagsVal.IsValid() || tagsVal.IsNil() {
		return nil, 0, fmt.Errorf("id value is nil")
	}
	objTags := tagsVal.Interface().([]string)

	return objTags, id, nil
}

func SaveTags(bean interface{}) error {
	objTags, id, e := tagsId(bean)
	if e != nil {
		// ignore if not tags field in struct
		return nil
	}
	if objTags == nil {
		return RemoveAllTags(bean)
	}
	for _, tag := range objTags {
		e := tags.SetTag(tag, objPrefix(bean), id)
		if e != nil {
			return e
		}
	}
	return nil
}

func GetTagItems(tagName string, bean interface{}) ([]int64, error) {
	return tags.GetTagItems(tagName, objPrefix(bean))
}

func RemoveTag(tagName string, bean interface{}) error {
	_, id, _ := tagsId(bean)
	e := tags.RemoveTag(tagName, objPrefix(bean), id)
	if e != nil {
		return e
	}
	return nil
}

func RemoveAllTags(bean interface{}) error {
	objTags, id, _ := tagsId(bean)
	for _, tag := range objTags {
		e := tags.RemoveTag(tag, objPrefix(bean), id)
		if e != nil {
			return e
		}
	}
	return nil
}
