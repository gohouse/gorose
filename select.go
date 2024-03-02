package gorose

import "strings"

// Column 表示SELECT语句中的列信息。
type Column struct {
	Name  string
	Alias string // 可选别名
	IsRaw bool   // 是否是原生SQL片段
	Binds []any  // 绑定数据
}

// SelectClause 存储SELECT子句相关信息。
type SelectClause struct {
	Columns  []Column
	Distinct bool
}

// Select specifies the columns to retrieve.
// Select("a","b")
// Select("a.id as aid","b.id bid")
// Select("id,nickname name")
func (db *SelectClause) Select(columns ...string) *SelectClause {
	for _, column := range columns {
		splits := strings.Split(column, ",")
		for _, split := range splits {
			parts := strings.Split(strings.TrimSpace(split), " ")
			switch len(parts) {
			case 3:
				db.Columns = append(db.Columns, Column{
					Name:  strings.TrimSpace(parts[0]),
					Alias: strings.TrimSpace(parts[2]),
				})
			case 2:
				db.Columns = append(db.Columns, Column{
					Name:  strings.TrimSpace(parts[0]),
					Alias: strings.TrimSpace(parts[1]),
				})
			case 1:
				db.Columns = append(db.Columns, Column{
					Name: strings.TrimSpace(parts[0]),
				})
			}
		}
	}
	return db
}
func (db *SelectClause) AddSelect(columns ...string) *SelectClause { return db.Select(columns...) }

// SelectRaw 允许直接在查询中插入原始SQL片段作为选择列。
func (db *SelectClause) SelectRaw(raw string, binds ...any) *SelectClause {
	db.Columns = append(db.Columns, Column{
		Name:  raw,
		IsRaw: true,
		Binds: binds,
	})
	return db
}
