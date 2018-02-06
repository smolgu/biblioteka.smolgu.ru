package book

// Copyright 2013 The Agostle Authors. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

import (
	"bytes"
	"github.com/fatih/color"
	"os/exec"
	"rsc.io/pdf"
	"strconv"
	"sync"
)

func PdfPageNum(srcfn string) (numberofpages int, err error) {
	r, e := pdf.Open(srcfn)
	if e != nil {
		return 0, e
	}
	return r.NumPage(), nil
}

// PdfPageNum returns the number of pages
/*func PdfPageNum(srcfn string) (numberofpages int, err error) {
	if numberofpages, _, err = pdfPageNum(srcfn); err == nil {
		return
	}
	numberofpages, _, err = pdfPageNum(srcfn)
	return
}*/

func pdfPageNum(srcfn string) (numberofpages int, encrypted bool, err error) {
	numberofpages = -1

	pdfinfo := false
	var cmd *exec.Cmd
	cmd = exec.Command("pdfinfo", srcfn)
	pdfinfo = true
	color.Cyan("pdfPageNum calling %v", cmd)
	out, e := cmd.CombinedOutput()
	err = e
	if 0 == len(out) {
		return
	}

	getLine := func(hay []byte, needle string) (ret string) {
		i := bytes.Index(hay, []byte("\n"+needle))
		if i >= 0 {
			line := hay[i+1+len(needle):]
			j := bytes.IndexByte(line, '\n')
			if j >= 0 {
				return string(bytes.Trim(line[:j], " \t\r\n"))
			}
		}
		return ""
	}

	if pdfinfo {
		encrypted = getLine(out, "Encrypted:") == "yes"
		numberofpages, err = strconv.Atoi(getLine(out, "Pages:"))
	}
	return
}

var (
	alreadyCleaned = make(map[string]bool, 16)
	cleanMtx       = sync.Mutex{}
	pdfCleanStatus = int(0)
)
