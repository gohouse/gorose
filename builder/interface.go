package builder

import (
	"github.com/gohouse/gorose/api"
)

type IBuilder interface {
	BuildQuery(api api.OrmApi) (string, error)
	BuildExecute(api api.OrmApi, operType string) (string, error)
}
