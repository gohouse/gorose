package gorose

import (
	"sync"
)

var (
	regex = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like",
		"in", "not in", "between", "not between"}
)

type Builder struct {
	IOrm
	driver   string
	builders map[string]IBuilder
	regex    []string
}

var onceBuilder sync.Once
var builder *Builder

func NewBuilder(driver string) *Builder {
	onceBuilder.Do(func() {
		builder = &Builder{builders: make(map[string]IBuilder), regex: regex}
	})
	builder.driver = driver
	return builder
}
func NewDriver() *Builder {
	return NewBuilder("")
}
func (b *Builder) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	return b.Getter(b.driver).BuildQuery(o)
}
func (b *Builder) BuildExecute(o IOrm, operType string) (sqlStr string, args []interface{}, err error) {
	return b.Getter(b.driver).BuildExecute(o, operType)
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

func (b *Builder) GetRegex() []string {
	return b.regex
}
