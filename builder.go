package gorose

import (
	"sync"
)

type Builder struct {
	IOrm
	driver   string
	builders map[string]IBuilder
}

var onceBuilder sync.Once
var builder *Builder

func NewBuilder(driver string) *Builder {
	onceBuilder.Do(func() {
		builder = &Builder{builders: make(map[string]IBuilder)}
	})
	builder.driver = driver
	//builder.IOrm = o
	return builder
	//return NewDriver().Getter(driver)
}
func NewDriver() *Builder {
	return NewBuilder("")
}
func (b *Builder) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return b.Getter(b.driver).BuildQuery(o)
}
func (b *Builder) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return
}

func (b *Builder) Register(driver string, val IBuilder) {
	withLockContext(func() {
		b.builders[driver] = val
	})
}

func (b *Builder) Getter(driver string) (ib IBuilder) {
	return b.builders[driver]
}

func (b *Builder) GetIOrm() IOrm {
	return b.IOrm
}
