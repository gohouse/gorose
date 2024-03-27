package gorose

type Context struct {
	TableClause       TableClause
	SelectClause      SelectClause
	JoinClause        JoinClause
	WhereClause       WhereClause
	GroupClause       GroupClause
	HavingClause      HavingClause
	OrderByClause     OrderByClause
	LimitOffsetClause LimitOffsetClause

	PessimisticLocking string
	Prefix             string
}

func NewContext(prefix string) *Context {
	return &Context{Prefix: prefix}
}

func (db *Context) Table(table any, alias ...string) *Context {
	db.TableClause.Table(table, alias...)
	return db
}
func (db *Context) Select(columns ...string) *Context {
	db.SelectClause.Select(columns...)
	return db
}
func (db *Context) AddSelect(columns ...string) *Context {
	db.SelectClause.AddSelect(columns...)
	return db
}
func (db *Context) SelectRaw(raw string, binds ...any) *Context {
	db.SelectClause.SelectRaw(raw, binds...)
	return db
}
func (db *Context) Join(table any, argOrFn ...any) *Context {
	db.JoinClause.Join(table, argOrFn...)
	return db
}
func (db *Context) LeftJoin(table any, argOrFn ...any) *Context {
	db.JoinClause.LeftJoin(table, argOrFn...)
	return db
}
func (db *Context) RightJoin(table any, argOrFn ...any) *Context {
	db.JoinClause.RightJoin(table, argOrFn...)
	return db
}
func (db *Context) CrossJoin(table any, argOrFn ...any) *Context {
	db.JoinClause.CrossJoin(table, argOrFn...)
	return db
}
func (db *Context) Where(column any, argsOrclosure ...any) *Context {
	db.WhereClause.Where(column, argsOrclosure...)
	return db
}
func (db *Context) OrWhere(column any, argsOrclosure ...any) *Context {
	db.WhereClause.OrWhere(column, argsOrclosure...)
	return db
}
func (db *Context) WhereRaw(raw string, bindings ...any) *Context {
	db.WhereClause.WhereRaw(raw, bindings...)
	return db
}
func (db *Context) OrWhereRaw(raw string, bindings ...any) *Context {
	db.WhereClause.OrWhereRaw(raw, bindings...)
	return db
}
func (db *Context) GroupBy(columns ...string) *Context {
	db.GroupClause.GroupBy(columns...)
	return db
}
func (db *Context) GroupByRaw(columns ...string) *Context {
	db.GroupClause.GroupByRaw(columns...)
	return db
}
func (db *Context) Having(column any, argsOrClosure ...any) *Context {
	db.HavingClause.Where(column, argsOrClosure...)
	return db
}
func (db *Context) OrHaving(column any, argsOrClosure ...any) *Context {
	db.HavingClause.OrWhere(column, argsOrClosure...)
	return db
}
func (db *Context) HavingRaw(raw string, argsOrClosure ...any) *Context {
	db.HavingClause.WhereRaw(raw, argsOrClosure...)
	return db
}
func (db *Context) OrHavingRaw(raw string, argsOrClosure ...any) *Context {
	db.HavingClause.OrWhereRaw(raw, argsOrClosure...)
	return db
}
func (db *Context) OrderBy(column string, directions ...string) *Context {
	db.OrderByClause.OrderBy(column, directions...)
	return db
}
func (db *Context) OrderByRaw(column string) *Context {
	db.OrderByClause.OrderByRaw(column)
	return db
}
func (db *Context) Limit(limit int) *Context {
	db.LimitOffsetClause.Limit = limit
	return db
}
func (db *Context) Offset(offset int) *Context {
	db.LimitOffsetClause.Offset = offset
	return db
}
func (db *Context) Page(num int) *Context {
	db.LimitOffsetClause.Page = num
	return db
}
func (db *Context) SharedLock() *Context {
	db.PessimisticLocking = "LOCK IN SHARE MODE"
	return db
}
func (db *Context) LockForUpdate() *Context {
	db.PessimisticLocking = "FOR UPDATE"
	return db
}
