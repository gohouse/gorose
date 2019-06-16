package gorose

type ISession interface {
	Close()
	//Table(bind interface{}) IOrm
	Bind(bind interface{}) ISession
	Begin() (err error)
	Rollback() (err error)
	Commit() (err error)
	Transaction(closer ...func(session ISession) error) (err error)
	Query(sqlstring string, args ...interface{}) error
	Execute(sqlstring string, args ...interface{}) (int64, error)
	GetMasterDriver() string
	GetSlaveDriver() string
	//GetBinder() *Binder
	LastInsertId() int64
	LastSql() string
	GetTableName() (string, error)
	//GetBindName() string
	GetBinder() IBinder
}
