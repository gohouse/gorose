package gorose

import "math"

// LimitOffsetClause 存储LIMIT和OFFSET信息。
type LimitOffsetClause struct {
	Limit  int
	Offset int
	Page   int
}

// Pagination 是用于分页查询结果的结构体，包含当前页数据及分页信息。
type Pagination struct {
	Limit       int              `json:"limit"`
	Pages       int              `json:"pages"`
	CurrentPage int              `json:"currentPage"`
	PrevPage    int              `json:"prevPage"`
	NextPage    int              `json:"nextPage"`
	Total       int64            `json:"total"`
	Data        []map[string]any `json:"data"`
}

func (db *Database) Paginate(obj ...any) (result Pagination, err error) {
	if len(obj) > 0 {
		db.Table(obj[0])
	}
	var count int64
	count, err = db.Count()
	if err != nil || count == 0 {
		return
	}
	if db.Context.LimitOffsetClause.Limit == 0 {
		db.Limit(15)
	}
	if db.Context.LimitOffsetClause.Page == 0 {
		db.Page(1)
	}

	res, err := db.Get()
	if err != nil {
		return result, err
	}

	result.Total = count
	result.Data = res
	result.Limit = db.Context.LimitOffsetClause.Limit
	result.Pages = int(math.Ceil(float64(count) / float64(db.Context.LimitOffsetClause.Limit)))
	result.CurrentPage = db.Context.LimitOffsetClause.Page
	result.PrevPage = db.Context.LimitOffsetClause.Page - 1
	result.NextPage = db.Context.LimitOffsetClause.Page + 1
	if db.Context.LimitOffsetClause.Page == 1 {
		result.PrevPage = 1
	}
	if db.Context.LimitOffsetClause.Page == result.Pages {
		result.NextPage = result.Pages
	}
	return
}
