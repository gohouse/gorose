package gorose

import (
	"sync"
)

type Builder struct {
}

func NewBuilder(driver string) IBuilder {
	return NewDriver().Getter(driver)
}

type Driver struct {
	builders map[string]IBuilder
}

var driverOnce sync.Once
var driver *Driver

func NewDriver() *Driver {
	driverOnce.Do(func() {
		driver = &Driver{
			builders: make(map[string]IBuilder),
		}
	})
	return driver
}

func (b *Driver) Register(driver string, val IBuilder) {
	withLockContext(func() {
		b.builders[driver] = val
	})
}

func (b *Driver) Getter(driver string) (ib IBuilder) {
	return b.builders[driver]
}
