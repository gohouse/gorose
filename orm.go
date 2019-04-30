package gorose

import (
	"errors"
	"fmt"
	"strings"
)

var (
	regex = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like", "in", "not in", "between", "not between"}
)

type OrmArgs struct {
	fields   []string        // fields
	where    [][]interface{} // where
	order    string          // order
	limit    int             // limit
	offset   int             // offset
	join     [][]interface{} // join
	distinct bool            // distinct
	union    string          // sum/count/avg/max/min
	group    string          // group
	having   string          // having
	data     interface{}     // data
}

type Orm struct {
	ISession
	*Binder
	*OrmArgs
}

func NewOrm(b *Binder) IOrm {
	var o = new(Orm)
	o.Binder = b

	return o
}

func (dba Orm) Hello() {
	fmt.Println("hello gorose orm struct")
}

// Fields : select fields
func (dba Orm) Fields(fields ...string) Orm {
	dba.fields = fields
	return dba
}

// AddFields : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddFields method:
func (dba Orm) AddFields(fields ...string) Orm {
	dba.fields = append(dba.fields, fields...)
	return dba
}

// Select : equals Fields()
func (dba Orm) Select(fields ...string) Orm {
	return dba.Fields(fields...)
}

// AddSelect : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddSelect method:
func (dba Orm) AddSelect(fields ...string) Orm {
	dba.fields = append(dba.fields, fields...)
	return dba
}

// Table : select table
func (dba Orm) Table(table string) Orm {
	dba.ISession.Table(table)
	return dba
}

// Data : insert or update data
func (dba Orm) Data(data interface{}) Orm {
	dba.data = data
	return dba
}

// Group : select group by
func (dba Orm) Group(group string) Orm {
	dba.group = group
	return dba
}

// GroupBy : equals Group()
func (dba Orm) GroupBy(group string) Orm {
	return dba.Group(group)
}

// Having : select having
func (dba Orm) Having(having string) Orm {
	dba.having = having
	return dba
}

// Order : select order by
func (dba Orm) Order(order string) Orm {
	dba.order = order
	return dba
}

// OrderBy : equal order
func (dba Orm) OrderBy(order string) Orm {
	return dba.Order(order)
}

// Limit : select limit
func (dba Orm) Limit(limit int) Orm {
	dba.limit = limit
	return dba
}

// Offset : select offset
func (dba Orm) Offset(offset int) Orm {
	dba.offset = offset
	return dba
}

// Page : select page
func (dba Orm) Page(page int) Orm {
	dba.offset = (page - 1) * dba.limit
	return dba
}

// Where : query or execute where condition, the relation is and
func (dba Orm) Where(args ...interface{}) Orm {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"and", args}

	dba.where = append(dba.where, w)

	return dba
}

// Join : select join query
func (dba Orm) Join(args ...interface{}) Orm {
	//dba.parseJoin(args, "INNER")
	dba.join = append(dba.join, []interface{}{"INNER", args})

	return dba
}

// Distinct : select distinct
func (dba Orm) Distinct() Orm {
	dba.distinct = true

	return dba
}

// Get : select more rows , relation limit set
func (dba Orm) Get() (error) {
	// 构建sql
	sqlStr,args,err := dba.BuildSql()
	if err != nil {
		return err
	}

	// 执行查询
	return dba.ISession.Query(sqlStr, args...)
}

// Insert : insert data and get affected rows
func (dba Orm) Insert() (int64, error) {
	// 构建sql
	sqlStr,args,err := dba.BuildSql("insert")
	if err != nil {
		return 0,err
	}

	return dba.ISession.Execute(sqlStr, args...)
}

// insertGetId : insert data and get id
func (dba Orm) InsertGetId() (int64, error) {
	_, err := dba.Insert()
	if err != nil {
		return 0, err
	}
	return dba.LastInsertId(),nil
}

// Update : update data
func (dba Orm) Update() (int64, error) {
	// 构建sql
	sqlStr,args,err := dba.BuildSql("update")
	if err != nil {
		return 0,err
	}

	return dba.ISession.Execute(sqlStr, args...)
}

// Delete : delete data
func (dba Orm) Delete() (int64, error) {
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
func (dba Orm) Increment(args ...interface{}) (int64, error) {
	argLen := len(args)
	var field string
	var value string = "1"
	var mode string = "+"
	switch argLen {
	case 1:
		field = args[0].(string)
	case 2:
		field = args[0].(string)
		switch args[1].(type) {
		case int:
			value = parseStr(args[1])
		case int64:
			value = parseStr(args[1])
		case float32:
			value = parseStr(args[1])
		case float64:
			value = parseStr(args[1])
		case string:
			value = args[1].(string)
		default:
			return 0, errors.New("第二个参数类型错误")
		}
	case 3:
		field = args[0].(string)
		switch args[1].(type) {
		case int:
			value = parseStr(args[1])
		case int64:
			value = parseStr(args[1])
		case float32:
			value = parseStr(args[1])
		case float64:
			value = parseStr(args[1])
		case string:
			value = args[1].(string)
		default:
			return 0, errors.New("第二个参数类型错误")
		}
		mode = args[2].(string)
	default:
		return 0, errors.New("参数数量只允许1个,2个或3个")
	}
	dba.Data(field + "=" + field + mode + value)
	return dba.Update()
}

// Decrement : auto Decrement -1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
func (dba Orm) Decrement(args ...interface{}) (int64, error) {
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
func (dba Orm) BuildSql(operType ...string) (string, []interface{}, error) {
	if len(operType)==0 || (len(operType)>0 && strings.ToLower(operType[0])=="select") {
		return NewBuilder(dba.ISession.GetDriver()).BuildQuery(dba)
	} else {
		return NewBuilder(dba.ISession.GetDriver()).BuildExecute(dba, strings.ToLower(operType[0]))
	}
}
