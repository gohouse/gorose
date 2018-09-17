package builder

import (
	"github.com/gohouse/gorose/across"
)

type IBuilder interface {
	BuildQuery(ormApi across.OrmApi) (sql string, err error)
	BuildExecute(ormApi across.OrmApi, operType string) (sql string, err error)
}
