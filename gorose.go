package gorose

import (
	"fmt"
)

func Open(conf ...interface{}) (e IEngin, err error) {
	// 驱动engin
	conn := NewEngin()
	if len(conf) == 0 {
		return
	}
	switch conf[0].(type) {
	// 传入的是单个配置
	case *Config:
		err = conn.bootSingle(conf[0].(*Config))
	// 传入的是集群配置
	case *ConfigCluster:
		conn.config = conf[0].(*ConfigCluster)
		err = conn.bootCluster()
	default:
		panic(fmt.Sprint("Open() need *gorose.Config or *gorose.ConfigCluster param, also can empty for build sql string only, but ",
			conf, " given"))
	}

	return conn, err
}
