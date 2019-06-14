package gorose

import (
	"errors"
	"fmt"
	"github.com/gohouse/t"
	"strings"
)

type OrmApi struct {
	table    string
	fields   []string
	where    [][]interface{}
	order    string
	limit    int
	offset   int
	join     [][]interface{}
	distinct bool
	union    string
	group    string
	having   string
	data     interface{}
}
type Orm struct {
	ISession
	*OrmApi
}

var _ IOrm = &Orm{}

func NewOrm(s ISession) *Orm {
	return &Orm{s, NewOrmApi()}
}

func (dba *Orm) Hello() {
	fmt.Println("hello gorose orm struct")
}

// Fields : select fields
func (dba *Orm) Table(tab interface{}) IOrm {
	dba.Bind(tab)
	//dba.dba = dba.IBinder.GetBindName())
	return dba
}

// Fields : select fields
func (dba *Orm) Fields(fields ...string) IOrm {
	dba.fields = fields
	return dba
}

// AddFields : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddFields method:
func (dba *Orm) AddFields(fields ...string) IOrm {
	_fields := dba.GetFields()
	_fields = append(_fields, fields...)
	dba.SetFields(_fields)
	return dba
}

// Distinct : select distinct
func (dba *Orm) Distinct() IOrm {
	dba.true = true

	return dba
}

// Data : insert or update data
func (dba *Orm) Data(data interface{}) IOrm {
	dba.SetData(data)
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

// _joinBuilder
func (dba *Orm) _joinBuilder(joinType string, args []interface{}) {
	dba.join = append(dba.join, []interface{}{joinType, args})
}

// Get : select more rows , relation limit set
func (dba *Orm) Get() (error) {
	// 构建sql
	sqlStr, args, err := dba.BuildSql()
	if err != nil {
		return err
	}

	// 执行查询
	return dba.ISession.Query(sqlStr, args...)
}

// Insert : insert data and get affected rows
func (dba *Orm) Insert(data ...interface{}) (int64, error) {
	if dba.GetData() == nil && len(data) > 0 {
		dba.data = data[0]
	}
	// 构建sql
	sqlStr, args, err := dba.BuildSql("insert")
	if err != nil {
		return 0, err
	}

	return dba.ISession.Execute(sqlStr, args...)
}

// insertGetId : insert data and get id
func (dba *Orm) InsertGetId() (int64, error) {
	_, err := dba.Insert()
	if err != nil {
		return 0, err
	}
	return dba.ISession.LastInsertId(),nil
}

// Update : update data
func (dba *Orm) Update() (int64, error) {
	// 构建sql
	sqlStr,args,err := dba.BuildSql("update")
	if err != nil {
		return 0,err
	}

	return dba.ISession.Execute(sqlStr, args...)
}

// Delete : delete data
func (dba *Orm) Delete() (int64, error) {
	// 构建sql
	sqlStr,args,err := dba.BuildSql("delete")
	if err != nil {
		return 0,err
	}

	return dba.ISession.Execute(sqlStr, args...)
}

// Increment : auto Increment +1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
// we can use this method as decrement with the third param as "-"
// orm.Increment("top") , orm.Increment("top", 2, "-")=orm.Decrement("top",2)
func (dba *Orm) Increment(args ...interface{}) (int64, error) {
	argLen := len(args)
	var field string
	var mode string = "+"
	var value string = "1"
	switch argLen {
	case 1:
		field = t.New(args[0]).String()
	case 2:
		field = t.New(args[0]).String()
		value = t.New(args[1]).String()
	case 3:
		field = t.New(args[0]).String()
		value = t.New(args[1]).String()
		mode = t.New(args[2]).String()
	default:
		return 0, errors.New("参数数量只允许1个,2个或3个")
	}
	dba.Data(field + "=" + field + mode + value)
	return dba.Update()
}

// Decrement : auto Decrement -1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
func (dba *Orm) Decrement(args ...interface{}) (int64, error) {
	arglen := len(args)
	switch arglen {
	case 1:
		args = append(args, 1)
		args = append(args, "-")
	case 2:
		args = append(args, "-")
	default:
		return 0, errors.New("Decrement参数个数有误")
	}
	return dba.Increment(args...)
}

// BuildSql
// operType(select, insert, update, delete)
func (dba *Orm) BuildSql(operType ...string) (string, []interface{}, error) {
	if len(operType) == 0 || (len(operType) > 0 && strings.ToLower(operType[0]) == "select") {
		return NewBuilder(dba.GetSlaveDriver()).BuildQuery(dba)
	} else {
		return NewBuilder(dba.GetMasterDriver()).BuildExecute(dba, strings.ToLower(operType[0]))
	}
}
