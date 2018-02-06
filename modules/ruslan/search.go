package ruslan

import (
	//"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sisteamnik/ruslanparser"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mu sync.Mutex

	RHost = "192.168.0.1"
	RPort = "210"
)

/*func Search(vals url.Values) (res []ruslanparser.Book, e error) {
	mu.Lock()
	defer mu.Unlock()

	mr, e := PrepareRequest(vals)
	if e != nil {
		return nil, e
	}
	if mr > 20 {
		mr = 20
	}

	recs, e := ruslanparser.GetResult(ssid, 1, mr)
	if e != nil {
		return nil, e
	}

	for _, v := range recs.Records {
		res = append(res, v.ToBook())
	}
	return
}*/

func PrepareRequest(vals url.Values) (mr int, e error) {
	vals.Set("HOST", RHost)
	vals.Set("PORT", RPort)
	ssid, e := Ssid()
	if e != nil {
		return 0, e
	}
	vals.Set("SESSION_ID", ssid)
	vals.Set("MAXRECORDS", "1")
	vals.Set("LANG", "rus")
	vals.Set("ACTION", "SEARCH")
	vals.Set("ESNAME", "B")
	resp, e := http.PostForm("http://92.241.99.100/Scripts/zgate.exe", vals)
	if e != nil {
		return 0, e
	}
	defer resp.Body.Close()
	bts, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return 0, e
	}
	return resultMax(bts)
}

func resultMax(r []byte) (maxRecords int, e error) {
	var (
		reg = regexp.MustCompile("Записи с [\\d]* по [\\d]* из ([\\d]*)")
		d   *goquery.Document
	)
	red := strings.NewReader(string(r))
	d, e = goquery.NewDocumentFromReader(red)
	if e != nil {
		return
	}
	t := d.Find("span.succ").Text()
	bts := reg.FindAllStringSubmatch(t, -1)
	if len(bts) > 0 {
		if len(bts[0]) > 1 {
			maxRecords, e = strconv.Atoi(bts[0][1])
			return
		}
	}
	return
}

var (
	smu  sync.Mutex
	ssid string
	last time.Time
)

func Ssid() (string, error) {
	smu.Lock()
	defer smu.Unlock()

	if ssid != "" && time.Since(last) < 8*time.Minute {
		return ssid, nil
	} else {
		var e error
		ssid, e = getSsid()
		if e != nil {
			return "", e
		}
		last = time.Now()
		return ssid, nil
	}
}

func getSsid() (string, error) {
	return ruslanparser.NewSessionSearchUrl("http://92.241.99.100/Scripts/zgate.exe?Init+test.xml,simple.xsl+rus")
}

func NewSsid() (string, error) {
	return ruslanparser.NewSessionSearchUrl("http://92.241.99.100/Scripts/zgate.exe?Init+test.xml,simple.xsl+rus")
}

/*func GetById(id string) (ruslanparser.Book, error) {
	var (
		res ruslanparser.Book
	)
	vals := url.Values{"DBNAME": {"BOOK", "REF", "STAT", "ЭЛЕКТР_РЕСУРСЫ", "ФРК", "ЭБ", "БЕЛ_СМОЛЯНЕ", "ПРЖЕВАЛЬСКИЙ"}}
	vals.Set("TERM_1", id)
	vals.Set("USE_1", "12")
	rec, e := Search(vals)
	if e != nil {
		return res, e
	}
	if len(rec) == 1 {
		res = rec[0]
		return res, nil
	} else {
		return res, fmt.Errorf("book with id %s not found", id)
	}
}*/
