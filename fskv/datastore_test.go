package fskv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFSKV(t *testing.T) {
	tmpPath, err := os.MkdirTemp("", "fskv.*")
	require.Nil(t, err)

	db, err := New(tmpPath)
	require.Nil(t, err)
	defer os.RemoveAll(tmpPath)

	// set then get
	err = db.Set("k1", "v1")
	require.Nil(t, err)

	var value string
	found, err := db.Get("k1", &value)
	require.Nil(t, err)
	require.True(t, found)
	require.Equal(t, "v1", value)

	// update then get
	err = db.Set("k1", "v2")
	require.Nil(t, err)

	found, err = db.Get("k1", &value)
	require.Nil(t, err)
	require.True(t, found)
	require.Equal(t, "v2", value)

	// list
	err = db.Set("k2", "v2")
	require.Nil(t, err)

	keys, err := db.ListKeys()
	require.Nil(t, err)
	require.Contains(t, keys, "k1")
	require.Contains(t, keys, "k2")
	require.Equal(t, len(keys), 2)

	// delete then get
	err = db.Delete("k1")
	require.Nil(t, err)

	found, err = db.Get("k1", &value)
	require.Nil(t, err)
	require.False(t, found)
}
