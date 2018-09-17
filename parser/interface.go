package parser

import (
	"github.com/gohouse/gorose/across"
)

type IParser interface {
	Parse(d string) (conf *across.DbConfigCluster, err error)
}
