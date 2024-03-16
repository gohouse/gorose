package gorose

// TableClause table clause
type TableClause struct {
	Tables any // table name or struct(slice) or subQuery
	Alias  string
}

func As(table any, alias string) TableClause {
	return TableClause{
		Tables: table,
		Alias:  alias,
	}
}

// Table sets the table name for the query.
func (db *TableClause) Table(table any, alias ...string) {
	var as string
	if len(alias) > 0 {
		as = alias[0]
	}
	db.Tables = table
	db.Alias = as
}
