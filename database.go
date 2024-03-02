package gorose

type Database struct {
	Driver  *Driver
	Engin   *Engin
	Context *Context
}

func NewDatabase(g *GoRose) *Database {
	return &Database{
		Driver:  NewDriver(GetDriver(g.driver)),
		Engin:   NewEngin(g),
		Context: NewContext(g.prefix),
	}
}
func (db *Database) Table(table any, alias ...string) *Database {
	db.Context.TableClause.Table(table, alias...)
	return db
}
// Select specifies the columns to retrieve.
// Select("a","b")
// Select("a.id as aid","b.id bid")
// Select("id,nickname name")
func (db *Database) Select(columns ...string) *Database {
	db.Context.SelectClause.Select(columns...)
	return db
}
// AddSelect 添加选择列
func (db *Database) AddSelect(columns ...string) *Database {
	db.Context.SelectClause.AddSelect(columns...)
	return db
}
// SelectRaw 允许直接在查询中插入原始SQL片段作为选择列。
func (db *Database) SelectRaw(raw string, binds ...any) *Database {
	db.Context.SelectClause.SelectRaw(raw, binds...)
	return db
}

// Join clause
func (db *Database) Join(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.Join(table, argOrFn...)
	return db
}
// LeftJoin clause
func (db *Database) LeftJoin(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.LeftJoin(table, argOrFn...)
	return db
}
// RightJoin clause
func (db *Database) RightJoin(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.RightJoin(table, argOrFn...)
	return db
}
// CrossJoin clause
func (db *Database) CrossJoin(table any, argOrFn ...any) *Database {
	db.Context.JoinClause.CrossJoin(table, argOrFn...)
	return db
}
func (db *Database) Where(column any, argsOrclosure ...any) *Database {
	db.Context.WhereClause.Where(column, argsOrclosure...)
	return db
}
func (db *Database) OrWhere(column any, argsOrclosure ...any) *Database {
	db.Context.WhereClause.OrWhere(column, argsOrclosure...)
	return db
}

// WhereRaw 在查询中添加一个原生SQL“where”条件。
//
// sql: 原生SQL条件字符串。
// bindings: SQL绑定参数数组。
func (db *Database) WhereRaw(raw string, bindings ...any) *Database {
	db.Context.WhereClause.WhereRaw(raw, bindings...)
	return db
}
func (db *Database) OrWhereRaw(raw string, bindings ...any) *Database {
	db.Context.WhereClause.OrWhereRaw(raw, bindings...)
	return db
}

// GroupBy 添加 GROUP BY 子句
func (db *Database) GroupBy(columns ...string) *Database {
	db.Context.GroupClause.GroupBy(columns...)
	return db
}
func (db *Database) GroupByRaw(columns ...string) *Database {
	db.Context.GroupClause.GroupByRaw(columns...)
	return db
}

// Having 添加 HAVING 子句, 同where
func (db *Database) Having(column any, argsOrClosure ...any) *Database {
	db.Context.HavingClause.Where(column, argsOrClosure...)
	return db
}
func (db *Database) OrHaving(column any, argsOrClosure ...any) *Database {
	db.Context.HavingClause.OrWhere(column, argsOrClosure...)
	return db
}

// HavingRaw 添加 HAVING 子句, 同where
func (db *Database) HavingRaw(raw string, argsOrClosure ...any) *Database {
	db.Context.HavingClause.WhereRaw(raw, argsOrClosure...)
	return db
}
func (db *Database) OrHavingRaw(raw string, argsOrClosure ...any) *Database {
	db.Context.HavingClause.OrWhereRaw(raw, argsOrClosure...)
	return db
}

// OrderBy adds an ORDER BY clause to the query.
func (db *Database) OrderBy(column string, directions ...string) *Database {
	db.Context.OrderByClause.OrderBy(column, directions...)
	return db
}
func (db *Database) OrderByRaw(column string) *Database {
	db.Context.OrderByClause.OrderByRaw(column)
	return db
}

// Limit 设置查询结果的限制数量。
func (db *Database) Limit(limit int) *Database {
	db.Context.LimitOffsetClause.Limit = limit
	return db
}

// Offset 设置查询结果的偏移量。
func (db *Database) Offset(offset int) *Database {
	db.Context.LimitOffsetClause.Offset = offset
	return db
}

// Page 页数,根据limit确定
func (db *Database) Page(num int) *Database {
	db.Context.LimitOffsetClause.Page = num
	return db
}
