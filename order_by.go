package gorose


type OrderByItem struct {
	Column    string
	Direction string // "asc" 或 "desc"
	IsRaw     bool
}

// OrderByClause 存储排序信息。
type OrderByClause struct {
	Columns []OrderByItem
}

// OrderBy adds an ORDER BY clause to the query.
func (db *OrderByClause) OrderBy(column string, directions ...string) {
	var direction string
	if len(directions) > 0 {
		direction = directions[0]
	}
	db.Columns = append(db.Columns, OrderByItem{
		Column:    column,
		Direction: direction,
	})
}
func (db *OrderByClause) OrderByRaw(column string) {
	db.Columns = append(db.Columns, OrderByItem{
		Column: column,
		IsRaw:  true,
	})
}
