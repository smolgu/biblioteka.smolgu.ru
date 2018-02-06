package book

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
)

type Book struct {
	Id int64

	Title    string
	PdfPath  string
	NumPages int
}

func ConvertPdfToImages(fpath, outpath string, page int) error {
	color.Cyan("Convert %s to %s page %d", fpath, outpath, page)
	c := exec.Command("./script/pdf2png.sh", fpath, outpath, fmt.Sprint(page))
	var b bytes.Buffer

	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	color.Cyan("%s", string(b.Bytes()))
	return c.Run()
}

func ExtractFirstPages(fpath, outpath string) error {
	color.Cyan("Extract %s to %s", fpath, outpath)
	c := exec.Command("./script/first_pages.sh", outpath, fpath)
	var b bytes.Buffer

	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	color.Cyan("%s", string(b.Bytes()))
	return c.Run()
}
