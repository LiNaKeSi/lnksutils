package cachekv

import (
	"os"
	"testing"

	"github.com/linakesi/lnksutils/fskv"
	"github.com/stretchr/testify/require"
)

func TestCacheKv(t *testing.T) {
	tmpPath, err := os.MkdirTemp("", "cachekv.*")
	require.Nil(t, err)

	backend, err := fskv.New(tmpPath)
	require.Nil(t, err)
	defer os.RemoveAll(tmpPath)

	cacheKv := New(backend)

	// set then get
	err = cacheKv.Set("k1", "v1")
	require.Nil(t, err)

	var value string
	found, err := cacheKv.Get("k1", &value)
	require.Nil(t, err)
	require.True(t, found)
	require.Equal(t, "v1", value)

	// update then get
	err = cacheKv.Set("k1", "v2")
	require.Nil(t, err)

	found, err = cacheKv.Get("k1", &value)
	require.Nil(t, err)
	require.True(t, found)
	require.Equal(t, "v2", value)
}
