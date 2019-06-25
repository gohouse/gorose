package gorose

const TAGNAME = "gorose"

func Open(conf ...interface{}) (engin *Engin, err error) {
	// 驱动engin
	engin,err = NewEngin(conf...)
	if err!=nil {
		if engin.GetLogger().EnableErrorLog() {
			engin.GetLogger().Error(err.Error())
		}
		return
	}

	// 使用默认的log, 如果自定义了logger, 则只需要调用 Use() 方法即可覆盖
	engin.Use(DefaultLogger())
	return
}
