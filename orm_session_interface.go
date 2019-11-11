package gorose

// IOrmSession ...
type IOrmSession interface {
	//Close()
	//Table(bind interface{}) IOrm
	//Bind(bind interface{}) ISession
	Begin() (err error)
	Rollback() (err error)
	Commit() (err error)
	//Transaction(closer ...func(session ISession) error) (err error)
	Query(sqlstring string, args ...interface{}) ([]Data, error)
	Execute(sqlstring string, args ...interface{}) (int64, error)
	//GetMasterDriver() string
	//GetSlaveDriver() string
	LastInsertId() int64
	LastSql() string
	//SetIBinder(b IBinder)
	//GetTableName() (string, error)
	GetIBinder() IBinder
	SetUnion(u interface{})
	GetUnion() interface{}
}
