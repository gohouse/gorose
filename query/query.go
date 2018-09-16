package query

import (
	"errors"
	"github.com/gohouse/gorose"
)

// BuildSql : build sql string , but not execute sql really
// operType : select/insert/update/delete
func BuildSql(dba *gorose.Database, operType ...string) (string, error) {
	switch len(operType) {
	case 0:
		return dba.BuildQuery()
	case 1:
		if operType[0] == "select" {
			return
			return dba.BuildQuery()
		} else {
			return dba.BuildExecut(operType[0])
		}
	default:
		return "", errors.New("参数有误")
	}
}
