package gorose

import (
	"sync"
)

type Bindings struct {
	sql4prepare string
	err         error
	bindings    []any
}

type IDriver interface {
	//ToSqlIncDec(c *Context, symbol string, data map[string]any) (sql4prepare string, values []any, err error)

	ToSql(c *Context) (sql4prepare string, binds []any, err error)
	ToSqlSelect(c *Context) (sql4prepare string, binds []any)
	ToSqlTable(c *Context) (sql4prepare string, values []any, err error)
	ToSqlJoin(c *Context) (sql4prepare string, binds []any, err error)
	ToSqlWhere(c *Context) (sql4prepare string, values []any, err error)
	ToSqlOrderBy(c *Context) (sql4prepare string)
	ToSqlLimitOffset(c *Context) (sqlSegment string, binds []any)
	ToSqlInsert(c *Context, obj any, args ...TypeToSqlInsertCase) (sqlSegment string, binds []any, err error)
	ToSqlUpdate(c *Context, arg any) (sqlSegment string, binds []any, err error)
	ToSqlDelete(c *Context, obj any) (sqlSegment string, binds []any, err error)
}

type Driver struct {
	IDriver
}

func NewDriver(d IDriver) *Driver {
	return &Driver{IDriver: d}
}

var driverMap = map[string]IDriver{}
var driverLock sync.RWMutex

func Register(driver string, parser IDriver) {
	driverLock.Lock()
	defer driverLock.Unlock()
	driverMap[driver] = parser
}

func GetDriver(driver string) IDriver {
	driverLock.RLock()
	defer driverLock.RUnlock()
	return driverMap[driver]
}

func DriverList() (dr []string) {
	driverLock.RLock()
	defer driverLock.RUnlock()
	for d := range driverMap {
		dr = append(dr, d)
	}
	return
}
