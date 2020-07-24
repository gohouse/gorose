package gorose

import (
	"github.com/gohouse/t"
	"strings"
)

//type TransactionHandlerFunc func(db IOrm) error
// Orm ...
type Orm struct {
	ISession
	//IBinder
	*OrmApi
	driver     string
	bindValues []interface{}
}

var _ IOrm = (*Orm)(nil)

// NewOrm ...
func NewOrm(e IEngin) *Orm {
	var orm = new(Orm)
	orm.SetISession(NewSession(e))
	//orm.IBinder = b
	orm.OrmApi = new(OrmApi)
	return orm
}

// Close ...
func (dba *Orm) Close() {
	dba.GetISession().Close()
}

// ExtraCols 额外的字段
func (dba *Orm) ExtraCols(args ...string) IOrm {
	dba.extraCols = append(dba.extraCols, args...)
	return dba
}

// ResetExtraCols ...
func (dba *Orm) ResetExtraCols() IOrm {
	dba.extraCols = []string{}
	return dba
}

// SetBindValues ...
func (dba *Orm) SetBindValues(v interface{}) {
	dba.bindValues = append(dba.bindValues, v)
}

// ClearBindValues ...
func (dba *Orm) ClearBindValues() {
	dba.bindValues = []interface{}{}
}

// GetBindValues ...
func (dba *Orm) GetBindValues() []interface{} {
	return dba.bindValues
}

// GetDriver ...
func (dba *Orm) GetDriver() string {
	return dba.driver
}

// SetISession ...
func (dba *Orm) SetISession(is ISession) {
	dba.ISession = is
}

// GetISession ...
func (dba *Orm) GetISession() ISession {
	return dba.ISession
}

// GetOrmApi ...
func (dba *Orm) GetOrmApi() *OrmApi {
	return dba.OrmApi
}

// Fields : select fields
func (dba *Orm) Table(tab interface{}) IOrm {
	dba.GetISession().Bind(tab)
	//dba.table = dba.GetISession().GetTableName()
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
	if len(args) == 0 ||
		t.New(args[0]).Bool() == false {
		return dba
	}
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

// WhereNull ...
func (dba *Orm) WhereNull(arg string) IOrm {
	return dba.Where(arg + " IS NULL")
}

// OrWhereNull ...
func (dba *Orm) OrWhereNull(arg string) IOrm {
	return dba.OrWhere(arg + " IS NULL")
}

// WhereNotNull ...
func (dba *Orm) WhereNotNull(arg string) IOrm {
	return dba.Where(arg + " IS NOT NULL")
}

// OrWhereNotNull ...
func (dba *Orm) OrWhereNotNull(arg string) IOrm {
	return dba.OrWhere(arg + " IS NOT NULL")
}

// WhereRegexp ...
func (dba *Orm) WhereRegexp(arg string, expstr string) IOrm {
	return dba.Where(arg, "REGEXP", expstr)
}

// OrWhereRegexp ...
func (dba *Orm) OrWhereRegexp(arg string, expstr string) IOrm {
	return dba.OrWhere(arg, "REGEXP", expstr)
}

// WhereNotRegexp ...
func (dba *Orm) WhereNotRegexp(arg string, expstr string) IOrm {
	return dba.Where(arg, "NOT REGEXP", expstr)
}

// OrWhereNotRegexp ...
func (dba *Orm) OrWhereNotRegexp(arg string, expstr string) IOrm {
	return dba.OrWhere(arg, "NOT REGEXP", expstr)
}

// WhereIn ...
func (dba *Orm) WhereIn(needle string, hystack []interface{}) IOrm {
	return dba.Where(needle, "IN", hystack)
}

// OrWhereIn ...
func (dba *Orm) OrWhereIn(needle string, hystack []interface{}) IOrm {
	return dba.OrWhere(needle, "IN", hystack)
}

// WhereNotIn ...
func (dba *Orm) WhereNotIn(needle string, hystack []interface{}) IOrm {
	return dba.Where(needle, "NOT IN", hystack)
}

// OrWhereNotIn ...
func (dba *Orm) OrWhereNotIn(needle string, hystack []interface{}) IOrm {
	return dba.OrWhere(needle, "NOT IN", hystack)
}

// WhereBetween ...
func (dba *Orm) WhereBetween(needle string, hystack []interface{}) IOrm {
	return dba.Where(needle, "BETWEEN", hystack)
}

// OrWhereBetween ...
func (dba *Orm) OrWhereBetween(needle string, hystack []interface{}) IOrm {
	return dba.OrWhere(needle, "BETWEEN", hystack)
}

// WhereNotBetween ...
func (dba *Orm) WhereNotBetween(needle string, hystack []interface{}) IOrm {
	return dba.Where(needle, "NOT BETWEEN", hystack)
}

// OrWhereNotBetween ...
func (dba *Orm) OrWhereNotBetween(needle string, hystack []interface{}) IOrm {
	return dba.OrWhere(needle, "NOT BETWEEN", hystack)
}

// Join : select join query
func (dba *Orm) Join(args ...interface{}) IOrm {
	dba._joinBuilder("INNER", args)
	return dba
}

// LeftJoin ...
func (dba *Orm) LeftJoin(args ...interface{}) IOrm {
	dba._joinBuilder("LEFT", args)
	return dba
}

// RightJoin ...
func (dba *Orm) RightJoin(args ...interface{}) IOrm {
	dba._joinBuilder("RIGHT", args)
	return dba
}

// CrossJoin ...
func (dba *Orm) CrossJoin(args ...interface{}) IOrm {
	dba._joinBuilder("CROSS", args)
	return dba
}

// _joinBuilder
func (dba *Orm) _joinBuilder(joinType string, args []interface{}) {
	dba.join = append(dba.join, []interface{}{joinType, args})
}

// Reset  orm api and bind values reset to init
func (dba *Orm) Reset() IOrm {
	dba.OrmApi = new(OrmApi)
	dba.ClearBindValues()
	dba.GetISession().SetUnion(nil)
	dba.ResetTable()
	return dba
}

// ResetTable ...
func (dba *Orm) ResetTable() IOrm {
	dba.GetISession().SetIBinder(NewBinder())
	return dba
}

// ResetWhere ...
func (dba *Orm) ResetWhere() IOrm {
	dba.where = [][]interface{}{}
	return dba
}

//// ResetUnion ...
//func (dba *Orm) ResetUnion() IOrm {
//	dba.union = ""
//	dba.GetISession().SetUnion(nil)
//	return dba
//}

// BuildSql
// operType(select, insert, update, delete)
func (dba *Orm) BuildSql(operType ...string) (a string, b []interface{}, err error) {
	// 解析table
	dba.table, err = dba.GetISession().GetTableName()
	if err != nil {
		dba.GetISession().GetIEngin().GetLogger().Error(err.Error())
		return
	}
	// 解析字段
	// 如果有union操作, 则不需要
	if inArray(dba.GetIBinder().GetBindType(), []interface{}{OBJECT_STRUCT, OBJECT_STRUCT_SLICE}) &&
		dba.GetUnion() == nil {
		dba.fields = getTagName(dba.GetIBinder().GetBindResult(), TAGNAME)
	}
	if len(operType) == 0 || (len(operType) > 0 && strings.ToLower(operType[0]) == "select") {
		//// 根据传入的struct, 设置limit, 有效的节约空间
		//if dba.union == "" {
		//	var bindType = dba.GetIBinder().GetBindType()
		//	if bindType == OBJECT_MAP || bindType == OBJECT_STRUCT {
		//		dba.Limit(1)
		//	}
		//}
		a, b, err = NewBuilder(dba.GetISession().GetIEngin().GetDriver()).BuildQuery(dba)
		if dba.GetISession().GetTransaction() {
			a = a + dba.GetPessimisticLock()
		}
	} else {
		a, b, err = NewBuilder(dba.GetISession().GetIEngin().GetDriver()).BuildExecute(dba, strings.ToLower(operType[0]))
		// 重置强制获取更新或插入的字段, 防止复用时感染
		dba.ResetExtraCols()
	}
	// 如果是事务, 因为需要复用单一对象, 故参数会产生感染
	// 所以, 在这里做一下数据绑定重置操作
	if dba.GetISession().GetTransaction() {
		dba.Reset()
	}
	// 这里统一清理一下绑定的数据吧, 万一要复用了, 造成绑定数据感染, 就尴尬了
	dba.ClearBindValues()
	return
}

// Transaction ...
func (dba *Orm) Transaction(closers ...func(db IOrm) error) (err error) {
	err = dba.ISession.Begin()
	if err != nil {
		dba.GetIEngin().GetLogger().Error(err.Error())
		return err
	}

	for _, closer := range closers {
		err = closer(dba)
		if err != nil {
			dba.GetIEngin().GetLogger().Error(err.Error())
			_ = dba.ISession.Rollback()
			return
		}
	}
	return dba.ISession.Commit()
}

// SharedLock 共享锁
// select * from xxx lock in share mode
func (dba *Orm) SharedLock() *Orm {
	dba.pessimisticLock = " lock in share mode"
	return dba
}

// LockForUpdate
// select * from xxx for update
func (dba *Orm) LockForUpdate() *Orm {
	dba.pessimisticLock = " for update"
	return dba
}
