package builder

import (
	"github.com/gohouse/gorose/across"
)

type IBuilder interface {
	BuildQuery(api across.OrmApi) (string, error)
	BuildExecute(api across.OrmApi, operType string) (string, error)
}
