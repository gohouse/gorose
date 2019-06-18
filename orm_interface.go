package gorose

type IOrm interface {
	IOrmApi
	IOrmQuery
	IOrmExecute
	IOrmSession
	//ISession
	Hello()
	BuildSql(operType ...string) (string, []interface{}, error)
	Table(tab interface{}) IOrm
	// fields=select
	Fields(fields ...string) IOrm
	AddFields(fields ...string) IOrm
	// distinct 方法允许你强制查询返回不重复的结果集：
	Distinct() IOrm
	Data(data interface{}) IOrm
	// groupBy, orderBy, having
	Group(group string) IOrm
	GroupBy(group string) IOrm
	Having(having string) IOrm
	Order(order string) IOrm
	OrderBy(order string) IOrm
	Limit(limit int) IOrm
	Offset(offset int) IOrm
	Page(page int) IOrm
	Where(args ...interface{}) IOrm
	OrWhere(args ...interface{}) IOrm
	// join(=innerJoin),leftJoin,rightJoin,crossJoin
	Join(args ...interface{}) IOrm
	LeftJoin(args ...interface{}) IOrm
	RightJoin(args ...interface{}) IOrm
	CrossJoin(args ...interface{}) IOrm
	//// truncate
	//Truncate()
	//// 悲观锁使用
	//// sharedLock(lock in share mode)
	//SharedLock()
	//// 此外你还可以使用 lockForUpdate 方法。“for update”锁避免选择行被其它共享锁修改或删除：
	//LockForUpdate()
	//GetRegex() []string
	GetDriver() string
	//GetBinder() IBinder
	SetBindValues(v interface{})
	GetBindValues() []interface{}
	ClearBindValues()
	//Transaction(closers ...func(db IOrm) error) (err error)
	Reset() IOrm
	ResetWhere() IOrm
	//GetBindName() string
}
