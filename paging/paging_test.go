package paging

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MyObject struct {
	Content string
}

var tdb *gorm.DB

var testData []MyObject = func() []MyObject {
	var ret []MyObject
	for i := 1; i < 103; i++ {
		ret = append(ret, MyObject{
			Content: fmt.Sprintf("value of %d", i),
		})
	}
	return ret
}()

func TestMain(m *testing.M) {
	var err error
	tdb, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		panic(err)
	}
	tdb.AutoMigrate(&MyObject{})
	for _, o := range testData {
		tdb.Create(o)
	}
	os.Exit(m.Run())
}

// see why need create new session for https://gorm.io/docs/method_chaining.html#New-Session-Mode
func (MyObject) PrepareGORMPaging(ctx interface{}, db *gorm.DB, setup PagingSetup) (*gorm.DB, error) {
	if setup.Keyword == "" {
		return db.Model(&MyObject{}).Session(&gorm.Session{}), nil
	}
	return db.Model(MyObject{}).Where("content like ?", "%"+setup.Keyword+"%").Session(&gorm.Session{}), nil
}

func TestPageable(t *testing.T) {
	result := WithPageable(nil, PagingSetup{Limit: 1, Page: 1}, tdb, MyObject{})
	assert.Equal(t, "", result.Err)
	assert.Equal(t, len(testData), int(result.Total))
	assert.Equal(t, 1, len(result.Data.([]MyObject)))

	result = WithPageable(nil, PagingSetup{Limit: 50, Page: 1}, tdb, MyObject{})
	assert.Equal(t, "", result.Err)
	assert.Equal(t, 50, len(result.Data.([]MyObject)))
	assert.Equal(t, len(testData), int(result.Total))

	result = WithPageable(nil, PagingSetup{Keyword: "1", Limit: 10, Page: 1}, tdb, MyObject{})
	assert.Equal(t, "", result.Err)
	assert.Equal(t, 10, len(result.Data.([]MyObject)))
	assert.Equal(t, 22, int(result.Total))
	result = WithPageable(nil, PagingSetup{Keyword: "1", Limit: 10, Page: 2}, tdb, MyObject{})
	assert.Equal(t, "", result.Err)
	assert.Equal(t, 10, len(result.Data.([]MyObject)))
	assert.Equal(t, 22, int(result.Total))
	result = WithPageable(nil, PagingSetup{Keyword: "1", Limit: 10, Page: 3}, tdb, MyObject{})
	assert.Equal(t, "", result.Err)
	assert.Equal(t, 2, len(result.Data.([]MyObject)))
	assert.Equal(t, 22, int(result.Total))

	result = WithPageable(nil, PagingSetup{Keyword: "dotesists", Limit: 10, Page: 1}, tdb, MyObject{})
	assert.Equal(t, "", result.Err)
	assert.Equal(t, 0, len(result.Data.([]MyObject)))
	assert.Equal(t, 0, int(result.Total))
}
