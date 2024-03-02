package gorose

//func (db *Database) MustFirst(columns ...string) (res map[string]any) {
//	res, db.Context.Err = db.First(columns...)
//	return
//}
//func (db *Database) MustGet(columns ...string) (res []map[string]any) {
//	res, db.Context.Err = db.Get(columns...)
//	return
//}

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
func (db *Database) WhereExists(clause IDriver) {
	db.Context.WhereClause.WhereExists(clause)
}
func (db *Database) WhereNotExists(clause IDriver) {
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

//func (db *Database) Paginate(obj ...any) (result Paginator, err error) {
//	if len(obj) > 0 {
//		db.Table(obj[0])
//	}
//	var count int64
//	count, err = db.Count()
//	if err != nil || count == 0 {
//		return
//	}
//	if db.Context.LimitOffsetClause.Limit == 0 {
//		db.Limit(15)
//	}
//	if db.Context.LimitOffsetClause.Page == 0 {
//		db.Page(1)
//	}
//
//	res, err := db.Get()
//	if err != nil {
//		return result, err
//	}
//
//	result.Total = count
//	result.Data = res
//	result.Limit = db.Context.LimitOffsetClause.Limit
//	result.Pages = int(math.Ceil(float64(count) / float64(db.Context.LimitOffsetClause.Limit)))
//	result.CurrentPage = db.Context.LimitOffsetClause.Page
//	result.PrevPage = db.Context.LimitOffsetClause.Page - 1
//	result.NextPage = db.Context.LimitOffsetClause.Page + 1
//	if db.Context.LimitOffsetClause.Page == 1 {
//		result.PrevPage = 1
//	}
//	if db.Context.LimitOffsetClause.Page == result.Pages {
//		result.NextPage = result.Pages
//	}
//	return
//}
