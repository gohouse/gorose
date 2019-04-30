package gorose

import (
	"sync"
)

type BuilderDriver struct {
	builders map[string]IBuilder
}

var builderDriverOnce sync.Once
var  builderDriver *BuilderDriver
func NewBuilderDriver() *BuilderDriver {
	builderDriverOnce.Do(func() {
		builderDriver = &BuilderDriver{make(map[string]IBuilder)}
	})
	return builderDriver
}

func (b *BuilderDriver) Register(driver string, val IBuilder)  {
	//builders[driver] = val
	b.builders[driver] = val
}

func (b *BuilderDriver) Getter(driver string) IBuilder {
	return b.builders[driver]
}
