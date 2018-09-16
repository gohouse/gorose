package parser

import "github.com/gohouse/gorose/config"

type IParser interface {
	Parse(d string) (conf *config.DbConfig, err error)
}
