package gorose

//func (db *Database) MustFirst(columns ...string) (res map[string]any) {
//	res, db.Context.Err = db.First(columns...)
//	return
//}
//func (db *Database) MustGet(columns ...string) (res []map[string]any) {
//	res, db.Context.Err = db.Get(columns...)
//	return
//}

func (db *Database) WhereSub(column string, operation string, sub WhereSubHandler) *Database {
	db.Context.WhereClause.WhereSub(column, operation, sub)
	return db
}
func (db *Database) OrWhereSub(column string, operation string, sub WhereSubHandler) *Database {
	db.Context.WhereClause.OrWhereSub(column, operation, sub)
	return db
}
func (db *Database) WhereBuilder(column string, operation string, sub IBuilder) *Database {
	db.Context.WhereClause.WhereBuilder(column, operation, sub)
	return db
}
func (db *Database) OrWhereBuilder(column string, operation string, sub IBuilder) *Database {
	db.Context.WhereClause.OrWhereBuilder(column, operation, sub)
	return db
}
func (db *Database) WhereNested(handler WhereNestedHandler) *Database {
	db.Context.WhereClause.WhereNested(handler)
	return db
}
func (db *Database) OrWhereNested(handler WhereNestedHandler) *Database {
	db.Context.WhereClause.OrWhereNested(handler)
	return db
}
func (db *Database) WhereIn(column string, value any) *Database {
	db.Context.WhereClause.whereIn("AND", column, value)
	return db
}
func (db *Database) WhereNotIn(column string, value any) *Database {
	db.Context.WhereClause.whereIn("AND", column, value, true)
	return db
}
func (db *Database) OrWhereIn(column string, value any) *Database {
	db.Context.WhereClause.whereIn("OR", column, value)
	return db
}
func (db *Database) OrWhereNotIn(column string, value any) *Database {
	db.Context.WhereClause.whereIn("OR", column, value, true)
	return db
}
func (db *Database) WhereNull(column string) *Database {
	db.Context.WhereClause.whereNull("AND", column)
	return db
}
func (db *Database) WhereNotNull(column string) *Database {
	db.Context.WhereClause.whereNull("AND", column, true)
	return db
}
func (db *Database) OrWhereNull(column string) *Database {
	db.Context.WhereClause.whereNull("OR", column)
	return db
}
func (db *Database) OrWhereNotNull(column string) *Database {
	db.Context.WhereClause.whereNull("OR", column, true)
	return db
}
func (db *Database) WhereBetween(column string, value any) *Database {
	db.Context.WhereClause.whereBetween("AND", column, value)
	return db
}
func (db *Database) WhereNotBetween(column string, value any) *Database {
	db.Context.WhereClause.whereBetween("AND", column, value, true)
	return db
}
func (db *Database) OrWhereBetween(column string, value any) *Database {
	db.Context.WhereClause.whereBetween("OR", column, value)
	return db
}
func (db *Database) OrWhereNotBetween(column string, value any) *Database {
	db.Context.WhereClause.whereBetween("OR", column, value, true)
	return db
}
func (db *Database) WhereExists(clause IBuilder) {
	db.Context.WhereClause.WhereExists(clause)
}
func (db *Database) WhereNotExists(clause IBuilder) {
	db.Context.WhereClause.WhereNotExists(clause)
}
func (db *Database) WhereLike(column, value string) *Database {
	db.Context.WhereClause.whereLike("AND", column, value)
	return db
}
func (db *Database) WhereNotLike(column, value string) *Database {
	db.Context.WhereClause.whereLike("AND", column, value, true)
	return db
}
func (db *Database) OrWhereLike(column, value string) *Database {
	db.Context.WhereClause.whereLike("OR", column, value)
	return db
}
func (db *Database) OrWhereNotLike(column, value string) *Database {
	db.Context.WhereClause.whereLike("OR", column, value, true)
	return db
}
func (db *Database) WhereNot(column any, args ...any) *Database {
	db.Context.WhereClause.WhereNot(column, args...)
	return db
}

func (db *Database) OrderByAsc(column string) *Database {
	return db.OrderBy(column, "ASC")
}

func (db *Database) OrderByDesc(column string) *Database {
	return db.OrderBy(column, "DESC")
}
