// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"log"

	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
)

func Migrate() error {
	x.Sync2(&Version{})
	v, err := GetVersion()
	if err != nil {
		return errors.Wrap(err, "get version")
	}
	log.Printf("database_version=%d app_version=%d", v.Version, len(migrations)+1)
	if v.Version < len(migrations)+1 {
		for _, m := range migrations[v.Version-1:] {
			log.Printf("start migration: %s", m.Description)
			err := m.Migration(x)
			if err != nil {
				return errors.Wrap(err, "migrate")
			}
			err = IncVersion()
			if err != nil {
				return errors.Wrap(err, "inc version")
			}
		}
	}

	return nil
}

var migrations = []struct {
	Description string
	Migration   func(db *xorm.Engine) error
}{
	{"drop bad bucket table", dropBucket},
}

func dropBucket(db *xorm.Engine) error {
	// omit error
	_, _ = db.Exec(`drop table bucket`)
	return nil
}
