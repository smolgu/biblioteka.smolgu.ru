package models

import (
	md5p "crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"github.com/smolgu/lib/modules/book"
	"github.com/ungerik/go-dry"
	"io"
	"os"
)

type Book struct {
	Id int64

	Title   string
	PdfPath string
	Pages   int
}

func CreateBook(title, pp string) (_ *Book, e error) {
	b := Book{}

	b.Title = title
	b.PdfPath = pp

	_, e = x.InsertOne(&b)

	return &b, e
}

func md5(in string) (out string) {
	return hex.EncodeToString(md5p.New().Sum([]byte(in)))
}

func CreateBookFromReader(title string, file io.Reader) (b *Book, e error) {
	b, e = CreateBook(title, "")
	if e != nil {
		return nil, e
	}

	var (
		hash       = md5(fmt.Sprint(b.Id))
		targetPath = "data/" + hash + ".pdf"
	)

	wr, e := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0777)
	if e != nil {
		return nil, e
	}
	_, e = io.Copy(wr, file)
	if e != nil {
		return nil, e
	}
	pagesNum, e := book.PdfPageNum(targetPath)
	if e != nil {
		return nil, e
	}
	b.Pages = pagesNum
	b.PdfPath = targetPath
	_, e = x.Id(b.Id).Update(b)
	if e != nil {
		return
	}
	return
}

func BookGet(id int64) (*Book, error) {
	b := new(Book)
	_, e := x.Id(id).Get(b)
	return b, e
}

func GetBookPage() {}

func WriteBookPage(w io.Writer, id int64, num int) (e error) {
	b := new(Book)

	_, e = x.Id(id).Get(b)
	if e != nil {
		return e
	}
	color.Green("%s", b.PdfPath)

	var (
		cachedPageDir  = "tmp/" + md5(fmt.Sprint(id)) + "/"
		cahcedPagePath = cachedPageDir + fmt.Sprint(num) + ".png"
		f              *os.File
	)

	if dry.FileExists(cahcedPagePath) {
		f, e = os.OpenFile(cahcedPagePath, os.O_RDONLY, 0777)
		if e != nil {
			return
		}
		defer f.Close()
	} else {
		os.MkdirAll(cachedPageDir, 0777)
		e = book.ConvertPdfToImages(b.PdfPath, cachedPageDir, num)
		if e != nil {
			return
		}
		f, e = os.OpenFile(cahcedPagePath, os.O_RDONLY, 0777)
		if e != nil {
			return
		}
		defer f.Close()
	}

	_, e = io.Copy(w, f)
	return
}
