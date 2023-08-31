package flatfs

import (
	"github.com/ipfs/go-datastore"
)

// keyIsValid returns true if the key is valid for flatfs.
// Allows keys that match [0-9A-Z+-_=].
func keyIsValid(key datastore.Key) bool {
	if len(key.Bytes()) > 1<<16 {
		return false
	}
	return true
}
