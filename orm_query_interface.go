package gorose

// IOrmQuery ...
type IOrmQuery interface {
	// 获取数据, 依据传入的绑定对象, 选择查询一条或多条数据并绑定到传入对象上
	// 当绑定对象传入的是string类型时, 返回多条结果集, 需要使用 Get() 来获取最终结果
	Select() error
	// 获取一条结果并返回, 只有当传入的table对象是字符串时生效
	First() (Data, error)
	// 获取多条结果并返回, 只有当传入的table对象是字符串时生效
	Get() ([]Data, error)
	// 如果你不需要完整的一行，可以使用 value 方法从结果中获取单个值，该方法会直接返回指定列的值：
	Value(field string) (v interface{}, err error)
	// 如果想要获取包含单个列值的数组，可以使用 pluck 方法
	// 还可以在返回数组中为列值指定自定义键（该自定义键必须是该表的其它字段列名，否则会报错）
	Pluck(field string, fieldKey ...string) (v interface{}, err error)
	// 查询构建器还提供了多个聚合方法，如count, max, min, avg 和 sum，你可以在构造查询之后调用这些方法：
	Count(args ...string) (int64, error)
	Sum(sum string) (interface{}, error)
	Avg(avg string) (interface{}, error)
	Max(max string) (interface{}, error)
	Min(min string) (interface{}, error)
	// 分页, 返回分页需要的基本数据
	Paginate(page ...int) (res Data, err error)
	// 组块结果集
	// 如果你需要处理成千上万或者更多条数据库记录，可以考虑使用 chunk 方法，该方法一次获取结果集的一小块，
	// 然后传递每一小块数据到闭包函数进行处理，该方法在编写处理大量数据库记录的 Artisan 命令的时候非常有用。
	// 例如，我们可以将处理全部 users 表数据分割成一次处理 100 条记录的小组块
	// 你可以通过从闭包函数中返回 err 来终止组块的运行
	Chunk(limit int, callback func([]Data) error) (err error)
	// 跟Chunk类似,只不过callback的是传入的结构体
	ChunkStruct(limit int, callback func() error) (err error)
	Loop(limit int, callback func([]Data) error) (err error)
}
