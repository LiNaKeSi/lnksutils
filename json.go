package utils

import (
	"encoding/json"
	"os"
)

func FileToJSON(p string, obj interface{}) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	return json.NewDecoder(f).Decode(obj)
}

func JSONToFile(dst string, obj interface{}) error {
	err := EnsureBaseDir(dst)
	if err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(obj)
}
