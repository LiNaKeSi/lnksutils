package fskv

import (
	"context"
	"errors"
	"strings"

	datastore "github.com/linakesi/lnksutils/fskv/internal/go-datastore"
	query "github.com/linakesi/lnksutils/fskv/internal/go-datastore/query"
	flatfs "github.com/linakesi/lnksutils/fskv/internal/go-ds-flatfs"

	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

type Store struct {
	db    *flatfs.Datastore
	Codec encoding.Codec
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

	data, err := s.Codec.Marshal(v)
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

	return true, s.Codec.Unmarshal(data, v)
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

func (s Store) ListKeys() (keys []string, err error) {
	res, err := s.db.Query(context.Background(), query.Query{})
	if err != nil {
		return nil, err
	}
	queryEntries, err := res.Rest()
	if err != nil {
		return nil, err
	}
	for _, e := range queryEntries {
		keys = append(keys, strings.TrimLeft(e.Key, "/"))
	}
	return keys, nil
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
	result.Codec = encoding.JSON
	return result, nil
}
