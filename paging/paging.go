package paging

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/linakesi/lnksutils/paging/internal/paginator"
)

type Pageable interface {
	PrepareGORMPaging(ctx interface{}, db *gorm.DB, setup PagingSetup) (*gorm.DB, error)
}

type PagingSetup struct {
	//limit records of per page. 使用ParseSetup时，若Limit=0则自动转换为最大值，即不处理分页直接返回所有内容
	Limit int `json:"limit" form:"limit"`
	Page  int `json:"page" form:"page"` //index start at one

	Keyword string `form:"keyword"`
	Order   string `form:"order"` //SQL Order statment

	Meta map[string]interface{} `form:"meta"` //任何其他元数据，比如"分类ID"等额外信息用来传递给PrepareGORMPaging做高级自定义搜索
}

func (p PagingSetup) GetInt(key string) (int, bool) {
	v, ok := p.Meta[key]
	if !ok {
		return 0, false
	}
	_v, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return int(_v), true
}

func (p PagingSetup) GetString(key string) (string, bool) {
	v, ok := p.Meta[key]
	if !ok {
		return "", false
	}
	_v, ok := v.(string)
	return _v, ok
}

func (p PagingSetup) GetBool(key string) (bool, bool) {
	v, ok := p.Meta[key]
	if !ok {
		return false, false
	}
	_v, ok := v.(bool)
	return _v, ok
}

type PagingResult struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
	Err   string      `json:"err"`
}

func WithPageable(ctx interface{}, setup PagingSetup, db *gorm.DB, pa Pageable) PagingResult {
	raw, err := pa.PrepareGORMPaging(ctx, db, setup)
	if err != nil {
		return PagingResult{
			Err: err.Error(),
		}
	}
	result := reflect.New(reflect.SliceOf(reflect.TypeOf(pa))).Interface()
	if gorm.IsRecordNotFoundError(raw.Error) {
		return PagingResult{
			Total: 0,
			Data:  reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(pa)), 0, 0).Interface(),
		}
	}
	return doQuery(paginator.New(paginator.NewGORMAdapter(raw), setup.Limit), setup.Page, result)
}

func ParseSetup(ctx *gin.Context) PagingSetup {
	p := PagingSetup{}
	ctx.BindQuery(&p)
	if p.Limit == 0 {
		p.Limit = int(^(uint(0)) >> 1)
	}
	return p
}

func WithSlice(ctx *gin.Context, all interface{}, result interface{}) PagingResult {
	setup := ParseSetup(ctx)
	pn := paginator.New(paginator.NewSliceAdapter(all), setup.Limit)
	return doQuery(pn, setup.Page, result)
}

func WithGorm(ctx *gin.Context, db *gorm.DB, result interface{}) PagingResult {
	setup := ParseSetup(ctx)
	pn := paginator.New(paginator.NewGORMAdapter(db), setup.Limit)
	return doQuery(pn, setup.Page, result)
}

func doQuery(p paginator.Paginator, page int, result interface{}) PagingResult {
	p.SetPage(page)
	p.Results(result)

	var errStr string
	n, err := p.Nums()
	if err != nil {
		errStr = err.Error()
	}
	return PagingResult{
		Total: n,
		Data:  reflect.ValueOf(result).Elem().Interface(),
		Err:   errStr,
	}
}
