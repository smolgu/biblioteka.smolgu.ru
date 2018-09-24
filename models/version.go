// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

type Version struct {
	Id      int64
	Version int
}

func GetVersion() (*Version, error) {
	v := new(Version)
	_, err := x.Id(1).Get(v)
	if err != nil || v.Version == 0 {
		v = new(Version)
		v.Version = 1
		_, err = x.Insert(v)
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

func IncVersion() error {
	v, err := GetVersion()
	if err != nil {
		return err
	}
	v.Version++
	_, err = x.Update(v)
	return err
}
