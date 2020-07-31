package gorose

import (
	"fmt"
	"regexp"
)

// IBuilder ...
type IBuilder interface {
	IFieldQuotes
	BuildQuery(orm IOrm) (sqlStr string, args []interface{}, err error)
	BuildExecute(orm IOrm, operType string) (sqlStr string, args []interface{}, err error)
	Clone() IBuilder
	//GetIOrm() IOrm
}

// IFieldQuotes 给系统关键词冲突的字段加引号,如: mysql是反引号, pg是双引号
type IFieldQuotes interface {
	AddFieldQuotes(field string) string
}

type FieldQuotesDefault struct {

}

func (FieldQuotesDefault) AddFieldQuotes(field string) string {
	reg := regexp.MustCompile(`^\w+$`)
	if reg.MatchString(field) {
		return fmt.Sprintf("`%s`", field)
	}
	return field
}
