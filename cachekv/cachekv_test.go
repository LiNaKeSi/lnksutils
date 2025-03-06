package cachekv

import (
	"os"
	"sync"
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

// go test -race
func TestCacheKvDataRace(t *testing.T) {
	tmpPath, err := os.MkdirTemp("", "cachekv.*")
	require.Nil(t, err)

	backend, err := fskv.New(tmpPath)
	require.Nil(t, err)
	defer os.RemoveAll(tmpPath)

	cacheKv := New(backend)

	cacheKv.Set("k1", []byte("a1"))

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var s []byte
			cacheKv.Get("k1", &s)
			s[0] = 'b'
		}()
	}
	wg.Wait()
}
