package utils

import (
	"encoding/json"
	"os"
)

//TODO: MOVE TO linakesi.com/utils

func FileToJSON(p string, obj interface{}) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	return json.NewDecoder(f).Decode(obj)
}
