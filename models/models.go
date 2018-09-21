package models

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/smolgu/lib/modules/setting"
	"github.com/ungerik/go-dry"
)

var (
	x      *xorm.Engine
	Trains []string

	NewCategoryId int64

	tables = []interface{}{new(Page),
		new(Category),
		new(Book),
		new(Link),
		new(Katalog),
		new(TrainingDirection),
		new(Vkr),
		new(Banner),
		new(Menu),
		new(MenuItem),
		new(Album),
		new(Photo),
		new(File),
		new(Router),
		new(Bucket),
	}
)

// Engine represents a xorm engine or session.
type Engine interface {
	Delete(interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Id(interface{}) *xorm.Session
	Sql(string, ...interface{}) *xorm.Session
	Where(string, ...interface{}) *xorm.Session
}

func NewEngine() {
	var err error
	log.Printf("open database file driver=%s setting=%s\n", setting.DbDriver, setting.DbSetting)
	x, err = xorm.NewEngine(setting.DbDriver, setting.DbSetting)
	if err != nil {
		err = errors.Wrapf(err, "open database file driver=%s setting=%s", setting.DbDriver, setting.DbSetting)
		log.Fatalln(err)
	}

	// logPath := filepath.Join(setting.LogDir, "sql.log")
	// f, e := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0666)
	// if e != nil {
	// 	panic(e)
	// }
	//logger := xorm.NewSimpleLogger(f)
	//x.SetLogger(logger)
	x.ShowSQL(true)

	x.Sync2(tables...)

	bts, err := dry.FileGetBytes(filepath.Join(setting.DataDir, "conf/trains.json"))
	if err != nil {
		log.Fatalf("open trains json: %s", err)
	}
	err = ffjson.Unmarshal(bts, &Trains)
	if err != nil {
		log.Fatalf("unmarshal trains json: %s", err)
	}

	var c Category
	_, err = x.Where("name = 'Новости'").Get(&c)
	if err != nil {
		log.Fatalf("get category: %s", err)
	}
	NewCategoryId = c.Id
}

func Delete(bean interface{}) (e error) {
	_, e = x.Delete(bean)
	return
}

func Save(bean interface{}) (e error) {
	val := reflect.ValueOf(bean)
	val = val.Elem()
	idVal := val.FieldByName("Id")
	if !idVal.IsValid() {
		return fmt.Errorf("id value is nil")
	}
	id := idVal.Interface().(int64)

	if id == 0 {
		_, e = x.InsertOne(bean)
		if e != nil {
			return
		}
		e = SaveTags(bean)
		return
	} else {
		_, e = x.Id(id).Update(bean)
		if e != nil {
			return
		}
		e = SaveTags(bean)
		return
	}
	return
}
