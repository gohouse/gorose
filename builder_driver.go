package gorose

import (
	"sync"
)

// BuilderDriver ...
type BuilderDriver struct {
	builders map[string]IBuilder
	b        *sync.Map
}

var builderDriverOnce sync.Once
var builderDriver *BuilderDriver

// NewBuilderDriver ...
func NewBuilderDriver() *BuilderDriver {
	builderDriverOnce.Do(func() {
		//builderDriver = &BuilderDriver{builders:make(map[string]IBuilder)}
		builderDriver = &BuilderDriver{b: &sync.Map{}}
	})
	return builderDriver
}

// Register ...
func (b *BuilderDriver) Register(driver string, val IBuilder) {
	//withLockContext(func() {
	//	b.builders[driver] = val
	//})
	b.b.Store(driver, val)
}

// Getter ...
func (b *BuilderDriver) Getter(driver string) IBuilder {

	//return b.builders[driver]
	if v, ok := (b.b.Load(driver)); ok {
		return v.(IBuilder)
	}
	return nil
}
