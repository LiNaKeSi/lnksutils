package fskv

import (
	"context"
	"errors"
	"strings"

	flatfs "github.com/linakesi/lnksutils/fskv/go-ds-flatfs"

	"github.com/ipfs/go-datastore"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

type Store struct {
	db    *flatfs.Datastore
	codec encoding.Codec
}

func checkKey(k string) error {
	if strings.Contains(k, "/") {
		return errors.New("The passed key has slash character, which is invalid")
	}
	return nil
}

func checkKeyAndValue(k string, v interface{}) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}
	return checkKey(k)
}

func (s Store) Set(k string, v interface{}) error {
	if err := checkKeyAndValue(k, v); err != nil {
		return err
	}

	data, err := s.codec.Marshal(v)
	if err != nil {
		return err
	}

	err = s.db.Put(context.Background(), datastore.NewKey(k), data)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) Get(k string, v interface{}) (found bool, err error) {
	if err := checkKeyAndValue(k, v); err != nil {
		return false, err
	}

	var data []byte
	switch data, err = s.db.Get(context.Background(), datastore.NewKey(k)); {
	case err == datastore.ErrNotFound:
		return false, nil
	case err != nil:
		return false, err
	}

	return true, s.codec.Unmarshal(data, v)
}

func (s Store) Delete(k string) error {
	if err := checkKey(k); err != nil {
		return err
	}

	return s.db.Delete(context.Background(), datastore.NewKey(k))
}

func (s Store) Close() error {
	return s.db.Close()
}

// 使用文件系统存储key-value数据
func New(storeDir string) (Store, error) {
	result := Store{}

	// Open DB
	db, err := flatfs.CreateOrOpen(storeDir, true)
	if err != nil {
		return result, err
	}
	result.db = db
	result.codec = encoding.JSON
	return result, nil
}
