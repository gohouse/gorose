package gorose

import (
	"fmt"
	"strings"
)

var (
	regex = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like",
		"in", "not in", "between", "not between"}
)

type Orm struct {
	ISession
	*OrmApi
	regex  []string
	driver string
}

var _ IOrm = &Orm{}

func NewOrm(s ISession) *Orm {
	return &Orm{ISession: s, OrmApi: &OrmApi{}, regex: regex}
}

func (dba *Orm) Hello() {
	fmt.Println("hello gorose orm struct")
}

func (dba *Orm) GetRegex() []string {
	return dba.regex
}

func (dba *Orm) GetDriver() string {
	return dba.driver
}

// Fields : select fields
func (dba *Orm) Table(tab interface{}) IOrm {
	dba.Bind(tab)
	//dba.table = dba.ISession.GetTableName()
	return dba
}

// Fields : select fields
func (dba *Orm) Fields(fields ...string) IOrm {
	dba.fields = fields
	return dba
}

// AddFields : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddFields method:
func (dba *Orm) AddFields(fields ...string) IOrm {
	dba.fields = append(dba.fields, fields...)
	return dba
}

// Distinct : select distinct
func (dba *Orm) Distinct() IOrm {
	dba.distinct = true

	return dba
}

// Data : insert or update data
func (dba *Orm) Data(data interface{}) IOrm {
	dba.data = data
	return dba
}

// Group : select group by
func (dba *Orm) Group(group string) IOrm {
	dba.group = group
	return dba
}

// GroupBy : equals Group()
func (dba *Orm) GroupBy(group string) IOrm {
	return dba.Group(group)
}

// Having : select having
func (dba *Orm) Having(having string) IOrm {
	dba.having = having
	return dba
}

// Order : select order by
func (dba *Orm) Order(order string) IOrm {
	dba.order = order
	return dba
}

// OrderBy : equal order
func (dba *Orm) OrderBy(order string) IOrm {
	return dba.Order(order)
}

// Limit : select limit
func (dba *Orm) Limit(limit int) IOrm {
	dba.limit = limit
	return dba
}

// Offset : select offset
func (dba *Orm) Offset(offset int) IOrm {
	dba.offset = offset
	return dba
}

// Page : select page
func (dba *Orm) Page(page int) IOrm {
	dba.offset = (page - 1) * dba.GetLimit()
	return dba
}

// Where : query or execute where condition, the relation is and
func (dba *Orm) Where(args ...interface{}) IOrm {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组
	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"and", args}

	dba.where = append(dba.where, w)

	return dba
}

// Where : query or execute where condition, the relation is and
func (dba *Orm) OrWhere(args ...interface{}) IOrm {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"or", args}

	dba.where = append(dba.where, w)

	return dba
}

// Join : select join query
func (dba *Orm) Join(args ...interface{}) IOrm {
	dba._joinBuilder("INNER", args)
	return dba
}
func (dba *Orm) LeftJoin(args ...interface{}) IOrm {
	dba._joinBuilder("LEFT", args)
	return dba
}
func (dba *Orm) RightJoin(args ...interface{}) IOrm {
	dba._joinBuilder("RIGHT", args)
	return dba
}
func (dba *Orm) CrossJoin(args ...interface{}) IOrm {
	dba._joinBuilder("CROSS", args)
	return dba
}

// _joinBuilder
func (dba *Orm) _joinBuilder(joinType string, args []interface{}) {
	dba.join = append(dba.join, []interface{}{joinType, args})
}

// BuildSql
// operType(select, insert, update, delete)
func (dba *Orm) BuildSql(operType ...string) (a string, b []interface{}, err error) {
	// 解析table
	dba.table, err = dba.ISession.GetTableName()
	//dba.table = dba.GetBindName()
	if err != nil {
		return
	}
	if len(operType) == 0 || (len(operType) > 0 && strings.ToLower(operType[0]) == "select") {
		// 根据传入的struct, 设置limit, 有效的节约空间
		if dba.union==""{
			var bindType = NewBinder().GetBindType()
			if bindType==OBJECT_MAP || bindType==OBJECT_STRUCT {
				dba.Limit(1)
			}
		}
		return NewBuilder(dba.GetSlaveDriver()).BuildQuery(dba)
	} else {
		return NewBuilder(dba.GetMasterDriver()).BuildExecute(dba, strings.ToLower(operType[0]))
	}
}
