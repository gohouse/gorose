package gorose

// IOrm ...
type IOrm interface {
	IOrmApi
	IOrmQuery
	IOrmExecute
	IOrmSession
	//ISession
	Close()
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
	// join(=innerJoin),leftJoin,rightJoin,crossJoin
	Join(args ...interface{}) IOrm
	LeftJoin(args ...interface{}) IOrm
	RightJoin(args ...interface{}) IOrm
	CrossJoin(args ...interface{}) IOrm
	// `Where`,`OrWhere`,`WhereNull / WhereNotNull`,`WhereIn / WhereNotIn / OrWhereIn / OrWhereNotIn`,`WhereBetween / WhereBetwee / OrWhereBetween / OrWhereNotBetween`
	Where(args ...interface{}) IOrm
	OrWhere(args ...interface{}) IOrm
	WhereNull(arg string) IOrm
	OrWhereNull(arg string) IOrm
	WhereNotNull(arg string) IOrm
	OrWhereNotNull(arg string) IOrm
	WhereRegexp(arg string, expstr string) IOrm
	OrWhereRegexp(arg string, expstr string) IOrm
	WhereNotRegexp(arg string, expstr string) IOrm
	OrWhereNotRegexp(arg string, expstr string) IOrm
	WhereIn(needle string, hystack []interface{}) IOrm
	OrWhereIn(needle string, hystack []interface{}) IOrm
	WhereNotIn(needle string, hystack []interface{}) IOrm
	OrWhereNotIn(needle string, hystack []interface{}) IOrm
	WhereBetween(needle string, hystack []interface{}) IOrm
	OrWhereBetween(needle string, hystack []interface{}) IOrm
	WhereNotBetween(needle string, hystack []interface{}) IOrm
	OrWhereNotBetween(needle string, hystack []interface{}) IOrm
	// truncate
	//Truncate()
	GetDriver() string
	//GetIBinder() IBinder
	SetBindValues(v interface{})
	GetBindValues() []interface{}
	ClearBindValues()
	Transaction(closers ...func(db IOrm) error) (err error)
	Reset() IOrm
	ResetTable() IOrm
	ResetWhere() IOrm
	GetISession() ISession
	GetOrmApi() *OrmApi
	// 悲观锁使用
	// sharedLock(lock in share mode) 不会阻塞其它事务读取被锁定行记录的值
	SharedLock() *Orm
	// 此外你还可以使用 lockForUpdate 方法。“for update”锁避免选择行被其它共享锁修改或删除：
	// 会阻塞其他锁定性读对锁定行的读取（非锁定性读仍然可以读取这些记录，lock in share mode 和 for update 都是锁定性读）
	LockForUpdate() *Orm
	//ResetUnion() IOrm
}
