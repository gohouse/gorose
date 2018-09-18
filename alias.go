package gorose

import "github.com/gohouse/gorose/across"

type Configure struct {
	across.DbConfigCluster
	across.DbConfigSingle
}
