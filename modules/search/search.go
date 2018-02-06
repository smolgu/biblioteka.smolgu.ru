package search

import (
	"log"

	"github.com/blevesearch/bleve"

	"github.com/Unknwon/com"
	//"github.com/go-xorm/xorm"

	"github.com/smolgu/lib/models"

	"fmt"
)

var (
	indexPath = "data/db/pages.index"

	index bleve.Index
)

func NewContext() {
	var e error
	if com.IsExist(indexPath) {
		index, e = bleve.Open(indexPath)
		if e != nil {
			log.Fatalf("err open index: %s", e)
		}
	} else {
		mapping := bleve.NewIndexMapping()
		mapping.DefaultAnalyzer = "en"
		index, e = bleve.New(indexPath, mapping)
		if e != nil {
			panic(e)
		}
	}
}

func Search(q string) (res []int64, e error) {
	query := bleve.NewMatchQuery(q)
	search := bleve.NewSearchRequest(query)
	searchResults, e := index.Search(search)
	if e != nil {
		return
	}

	for _, v := range searchResults.Hits {
		res = append(res, com.StrTo(v.ID).MustInt64())
	}
	return
}

func Index(id int64, text string) error {
	return index.Index(fmt.Sprint(id), text)
}
func IndexPages(pages []*models.Page) error {
	for _, v := range pages {
		fmt.Println("Index", v.Title)
		e := Index(v.Id, v.Title+" "+v.Body)
		if e != nil {
			return e
		}
	}
	return nil
}
