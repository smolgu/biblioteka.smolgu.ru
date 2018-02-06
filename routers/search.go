package routers

import (
	"fmt"
	"io"
	"os"

	"github.com/smolgu/lib/modules/middleware"
	"github.com/smolgu/lib/modules/ruslan"
)

func SearchHistory(c *middleware.Context) {
	c.Resp.Header().Set("Content-type", "text/plain; charset=utf8")
	c.Resp.WriteHeader(200)

	f, e := os.OpenFile("data/db/search.log", os.O_RDONLY, 0777)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer f.Close()
	_, e = io.Copy(c.Resp, f)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func SearchPerson(c *middleware.Context) {
	var (
		ruslanIdKey = "ruslan_ssid"
		ssidFace    = c.Session.Get(ruslanIdKey)
		ssid        string
		ok          bool

		baseFmt = "http://92.241.99.100/Scripts/zgate.exe?ACTION=SEARCH&DBNAME=БЕЛ_СМОЛЯНЕ&ESNAME=B&HOST=192.168.0.1&LANG=rus&MAXRECORDS=20&PORT=210&SESSION_ID=%s&SHOW_HOLDINGS=on&TERM_1=%s&USE_1=1009"
	)
	if ssid, ok = ssidFace.(string); !ok {
		newSsid, err := ruslan.NewSsid()
		if err != nil {
			return
		}
		ssid = newSsid
		c.Session.Set(ruslanIdKey, newSsid)
	}

	c.Redirect(fmt.Sprintf(baseFmt, ssid, c.Query("q")), 302)
}
