package gorose


type IOrmExecute interface {
	GetForce() bool
	// insert,insertGetId
	Insert(data ...interface{}) (int64, error)
	InsertGetId() (int64, error)
	Update() (int64, error)
	// updateOrInsert

	// increment,decrement
	// 在操作过程中你还可以指定额外的列进行更新：
	Increment(args ...interface{}) (int64, error)
	Decrement(args ...interface{}) (int64, error)
	// delete
	Delete() (int64, error)
	LastInsertId() int64
	//GetTableName() (string, error)
	//GetBindName() string
}
