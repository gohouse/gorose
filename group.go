package gorose

type TypeGroupItem struct {
	Column string
	IsRaw  bool
}
type GroupClause struct {
	Groups []TypeGroupItem
}
// GroupBy 添加 GROUP BY 子句
func (db *GroupClause) GroupBy(columns ...string) {
	for _, col := range columns {
		db.Groups = append(db.Groups, TypeGroupItem{
			Column: col,
		})
	}
}
func (db *GroupClause) GroupByRaw(columns ...string) {
	for _, col := range columns {
		db.Groups = append(db.Groups, TypeGroupItem{
			Column: col,
			IsRaw:  true,
		})
	}
}
