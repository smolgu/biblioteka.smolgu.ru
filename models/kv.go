package models

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/smolgu/lib/modules/setting"
)

var (
	boltDb *bolt.DB

	kvBucketName = []byte("kv")
	ErrNotFound  = fmt.Errorf("not found")
)

func NewKVContext() {
	if err := initKv(); err != nil {
		log.Fatalf("init kv: %s", err)
	}
}

func kvCreateDefaultBuckets(tx *bolt.Tx) error {
	_, e := tx.CreateBucketIfNotExists(kvBucketName)
	return e
}

func initKv() error {
	var (
		dbPath = filepath.Join(setting.DataDir, "db/kv.bolt")
	)

	var (
		e error
	)

	boltDb, e = bolt.Open(dbPath, 0777, nil)
	if e != nil {
		return e
	}
	e = boltDb.Update(kvCreateDefaultBuckets)
	if e != nil {
		return e
	}

	return nil
}

func makeSetFunc(key string, value []byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		buck := tx.Bucket(kvBucketName)
		if buck == nil {
			return bolt.ErrBucketNotFound
		}
		return buck.Put([]byte(key), value)
	}
}

func Set(key string, value []byte) error {
	boltDb.Update(makeSetFunc(key, value))
	return nil
}

func makeGetFunc(key string, value *[]byte) func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		buck := tx.Bucket(kvBucketName)
		if buck == nil {
			return bolt.ErrBucketNotFound
		}
		bts := buck.Get([]byte(key))
		if value == nil {
			return ErrNotFound
		}
		*value = bts
		return nil
	}
}

func Get(key string) ([]byte, error) {
	var value []byte
	e := boltDb.View(makeGetFunc(key, &value))
	return value, e
}
