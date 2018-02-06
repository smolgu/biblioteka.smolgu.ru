package users

import (
	"github.com/pquerna/ffjson/ffjson"
	"github.com/syndtr/goleveldb/leveldb"
)

type KV struct {
	levelDb *leveldb.DB
}

func NewKV(path string) (*KV, error) {
	ldb, e := leveldb.OpenFile(path, nil)
	if e != nil {
		return nil, e
	}
	kv := new(KV)
	kv.levelDb = ldb
	return kv, nil
}

func (kv *KV) Set(key string, data interface{}) error {
	bts, e := ffjson.Marshal(data)
	if e != nil {
		return e
	}
	e = kv.levelDb.Put([]byte(key), bts, nil)
	if e != nil {
		return e
	}
	return nil
}

func (kv *KV) Get(key string) ([]byte, error) {
	return kv.levelDb.Get([]byte(key), nil)
}

func (kv *KV) GetString(key string) (string, error) {
	bts, e := kv.Get(key)
	if e != nil {
		return "", e
	}
	return string(bts), nil
}
