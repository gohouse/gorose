package gorose

type IOrm interface {
	Hello()
	BuildSql(operType ...string) (string, []interface{}, error)
	//// distinct 方法允许你强制查询返回不重复的结果集：
	//Distinct()
	//// groupBy, orderBy, having
	//GroupBy()
	//OrderBy()
	//Having()
	//// join=innerJoin,leftJoin,rightJoin,crossJoin
	//Join()
	//LeftJoin()
	//RightJoin()
	//CrossJoin()
	//// fields=select
	//Fields()
	//// truncate
	//Truncate()
	//// 悲观锁使用
	//// sharedLock(lock in share mode)
	//SharedLock()
	//// 此外你还可以使用 lockForUpdate 方法。“for update”锁避免选择行被其它共享锁修改或删除：
	//LockForUpdate()
}

type IOrmQuery interface {
	BuildSql(operType string) (string, []interface{}, error)
	// 如果你只是想要从数据表中获取一行数据，可以使用 first 方法
	First()
	Get()
	// 如果你不需要完整的一行，可以使用 value 方法从结果中获取单个值，该方法会直接返回指定列的值：
	Value()
	// 如果想要获取包含单个列值的数组，可以使用 pluck 方法
	// 还可以在返回数组中为列值指定自定义键（该自定义键必须是该表的其它字段列名，否则会报错）
	Pluck()
	// 组块结果集
	// 如果你需要处理成千上万或者更多条数据库记录，可以考虑使用 chunk 方法，该方法一次获取结果集的一小块，
	// 然后传递每一小块数据到闭包函数进行处理，该方法在编写处理大量数据库记录的 Artisan 命令的时候非常有用。
	// 例如，我们可以将处理全部 users 表数据分割成一次处理 100 条记录的小组块
	// 你可以通过从闭包函数中返回 false 来终止组块的运行
	Chunk()
	// 查询构建器还提供了多个聚合方法，如count, max, min, avg 和 sum，你可以在构造查询之后调用这些方法：
	Count()
	Max()
	Min()
	Avg()
	Sum()
}

type IOrmExecute interface {
	// insert,insertGetId
	Insert()
	InsertGetId()
	// update
	Update()
	// updateOrInsert
	// increment,decrement
	// 在操作过程中你还可以指定额外的列进行更新：
	Increment()
	Decrement()
	// delete
	Delete()
}
