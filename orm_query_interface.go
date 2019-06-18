package gorose

import "github.com/gohouse/t"

type IOrmQuery interface {
	// 如果你只是想要从数据表中获取一行数据，可以使用 first 方法
	//First()
	Get() error
	GetAll() ([]Map,error)
	// 如果你不需要完整的一行，可以使用 value 方法从结果中获取单个值，该方法会直接返回指定列的值：
	Value(field string) (v t.T, err error)
	// 如果想要获取包含单个列值的数组，可以使用 pluck 方法
	// 还可以在返回数组中为列值指定自定义键（该自定义键必须是该表的其它字段列名，否则会报错）
	Pluck(field string, fieldKey ...string) (v t.T, err error)
	//// 组块结果集
	//// 如果你需要处理成千上万或者更多条数据库记录，可以考虑使用 chunk 方法，该方法一次获取结果集的一小块，
	//// 然后传递每一小块数据到闭包函数进行处理，该方法在编写处理大量数据库记录的 Artisan 命令的时候非常有用。
	//// 例如，我们可以将处理全部 users 表数据分割成一次处理 100 条记录的小组块
	//// 你可以通过从闭包函数中返回 false 来终止组块的运行
	//Chunk()
	// 查询构建器还提供了多个聚合方法，如count, max, min, avg 和 sum，你可以在构造查询之后调用这些方法：
	Count(args ...string) (int64, error)
	Sum(sum string) (interface{}, error)
	Avg(avg string) (interface{}, error)
	Max(max string) (interface{}, error)
	Min(min string) (interface{}, error)
	Paginate(limit,current_page int) (res Data, err error)
	Chunk(limit int, callback func([]Map) error) (err error)
	Loop(limit int, callback func([]Map) error) (err error)
}
