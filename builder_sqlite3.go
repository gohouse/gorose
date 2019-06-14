package gorose

import (
	"fmt"
	"strings"
)

type BuilderSqlite3 struct {
	IOrm
	args []interface{}
}

// sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
//		distinct, fields, table, join, where, group, having, order, limit, offset)
// select {distinct} {fields} from {table} {join} {where} {group} {having} {order} {limit} {offset}
func init() {
	NewBuilderDriver().Register("sqlite3", &BuilderSqlite3{})
}

func (b *BuilderSqlite3) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	b.IOrm = o
	sqlStr = fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
			b.BuildDistinct(), b.BuildFields(), b.BuildTable(), b.BuildJoin(), b.BuildWhere(),
			b.BuildGroup(), b.BuildHaving(), b.BuildOrder(), b.BuildLimit(), b.BuildOffset())
	return
}

func (b *BuilderSqlite3) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	b.IOrm = o
	sqlStr = "execute xxx test sql"
	return
}

func (b *BuilderSqlite3) BuildDistinct() (dis string) {
	if b.IOrm.GetDistinct() {
		dis = "DISTINCT "
	}
	return
}

func (b *BuilderSqlite3) BuildFields() string {
	return strings.Join(b.IOrm.GetFields(), ",")
}

func (b *BuilderSqlite3) BuildTable() string {
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildJoin() string {
	join := b.IOrm.GetJoin()
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildWhere() string {
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildGroup() string {
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildHaving() string {
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildOrder() string {
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildLimit() string {
	return b.IOrm.GetTable()
}

func (b *BuilderSqlite3) BuildOffset() string {
	return b.IOrm.GetTable()
}
