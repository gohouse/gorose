package gorose

import (
	"github.com/gohouse/gorose/across"
	"github.com/gohouse/gorose/builder"
	"github.com/gohouse/gorose/parser"
)

// 单一数据库配置
type DbConfigSingle struct {
	Driver          string        // 驱动: mysql/sqlite3/oracle/mssql/postgres
	EnableQueryLog  bool          // 是否开启sql日志
	SetMaxOpenConns int           // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns int           // (连接池)闲置的连接数
	Prefix          string        // 表前缀
	Dsn             string        // 数据库链接
}

// 数据库集群配置
// 如果不启用集群, 则直接使用 DbConfig 即可
// 如果仍然使用此配置为非集群, 则 Slave 配置置空即可, 等同于使用 DbConfig
type DbConfigCluster struct {
	Slave  []*DbConfigSingle // 多台读服务器, 如果启用则需要放入对应的多台从服务器配置
	Master *DbConfigSingle   // 一台主服务器负责写数据
}

func NewDbConfigCluster() *DbConfigCluster {
	return &DbConfigCluster{}
}

//// 来自于 struct 和 bytes 的互转大法
//func (s *DbConfigCluster) structToBytes(structTmp *across.DbConfigCluster) []byte {
//	type sliceMock struct {
//		addr uintptr
//		len  int
//		cap  int
//	}
//	Len := unsafe.Sizeof(*structTmp)
//	bytesTmp := &sliceMock{
//		addr: uintptr(unsafe.Pointer(structTmp)),
//		cap:  int(Len),
//		len:  int(Len),
//	}
//	return *(*[]byte)(unsafe.Pointer(bytesTmp))
//}
//// 来自于 struct 和 bytes 的互转大法
//func (s *DbConfigCluster) bytesToStruct(b []byte) *DbConfigCluster {
//	return *(**DbConfigCluster)(unsafe.Pointer(&b))
//}
//// struct 赋值自 另一个 struct
//func (s *DbConfigCluster) StructFrom(from *across.DbConfigCluster) *DbConfigCluster {
//	return s.bytesToStruct(s.structToBytes(from))
//}

// NewBuilder : sql构造器
func NewBuilder(ormApi across.OrmApi, operType ...string) (string, error) {
	return builder.NewBuilder(ormApi, operType...)
}

// NewFileParser : 配置解析器
func NewFileParser(fileOrDriverType, dsnOrFile string) (*DbConfigCluster, error) {
	var dbConf *DbConfigCluster
	var err error
	//var c *across.DbConfigCluster

	err = parser.NewFileParser(fileOrDriverType, dsnOrFile, &dbConf)
	//// 看我移形换位大法
	//dbConf = dbConf.StructFrom(c)

	return dbConf, err
}
