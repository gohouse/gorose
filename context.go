package gorose

type Context struct {
	TableClause       *TableClause
	SelectClause      *SelectClause
	JoinClause        *JoinClause
	WhereClause       *WhereClause
	GroupClause       *GroupClause
	HavingClause      *HavingClause
	OrderByClause     *OrderByClause
	LimitOffsetClause *LimitOffsetClause

	PessimisticLocking string
	Prefix             string
}

func NewContext(prefix string) *Context {
	return &Context{Prefix: prefix}
}
