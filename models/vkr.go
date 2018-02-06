package models

import (
	"github.com/Unknwon/com"
	"github.com/dchest/uniuri"
	"github.com/smolgu/lib/modules/book"
	"os"
	"time"
)

type TrainingDirection struct {
	Id int64

	Title    string
	SubTitle string

	FacultyId int64

	Code string

	CreatedBy int64

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted time.Time `xorm:"deleted"`
}

func TrainingDirectionAllList() ([]*TrainingDirection, error) {
	var (
		res []*TrainingDirection
	)
	e := x.Find(&res)
	return res, e
}

func TrainingDirectionFacultyList(facId int64) ([]*TrainingDirection, error) {
	var (
		res []*TrainingDirection
	)
	e := x.Where("faculty_id = ?", facId).Find(&res)
	return res, e
}

func TrainingDirectionById(id int64) (*TrainingDirection, error) {
	var (
		td = new(TrainingDirection)
	)

	_, e := x.Id(id).Get(td)

	return td, e
}

type Vkr struct {
	Id int64

	Title string

	FirstName  string
	LastName   string
	Patronymic string

	TrainingDirectionId int64

	Level     string
	FormStudy string

	BookId int64

	UploaderId int64

	Year int

	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
	Deleted time.Time `xorm:"deleted"`
}

func TrainingDirectionYears(tdId int64) ([]int, error) {
	var (
		r []int
	)
	res, e := x.Query("select distinct year from vkr where training_direction_id = ?", tdId)
	if e != nil {
		return nil, e
	}
	for _, v := range res {
		r = append(r, com.StrTo(string(v["year"])).MustInt())
	}
	return r, nil
}

func DirectionYearList(tdId int64, year int) ([]*Vkr, error) {
	var (
		res []*Vkr
	)

	e := x.Where("training_direction_id = ? and year = ?", tdId, year).Find(&res)
	if e != nil {
		return nil, e
	}
	return res, nil
}

func VkrGet(id int64) (*Vkr, error) {
	var (
		v = new(Vkr)
	)
	_, e := x.Id(id).Get(v)
	return v, e
}

func VkrDownload(id int64) (string, error) {
	vkr, e := VkrGet(id)
	if e != nil {
		return "", e
	}

	boook, e := BookGet(vkr.BookId)
	if e != nil {
		return "", e
	}

	name := "./tmp/" + uniuri.New() + ".pdf"

	f, e := os.Create(name)
	if e != nil {
		return "", e
	}
	f.Close()

	e = book.ExtractFirstPages(boook.PdfPath, name)
	if e != nil {
		return "", e
	}

	return name, nil

	//dry.FileCopy(name, dest)

	/*_, e = io.Copy(wr, f)

	e = f.Close()
	if e != nil {
		return e
	}*/

	/*e = os.Remove(name)
	if e != nil {
		return e
	}*/
}
