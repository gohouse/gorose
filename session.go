package gorose

func (dba *Database) Table(arg interface{}) *Database {
	dba.STable = arg
	return dba
}

func (dba *Database) Get() (result []map[string]interface{}, err error) {
	var sqlStr string
	sqlStr, err = dba.BuildSql()
	if err != nil {
		return
	}
	result, err = dba.Query(sqlStr)

	return
}

func (dba *Database) First() (result map[string]interface{}, err error) {
	dba.SLimit = 1
	var resultSlice []map[string]interface{}
	if resultSlice, err = dba.Get(); err != nil {
		return
	}
	if len(resultSlice) > 0 {
		result = resultSlice[0]
	}
	return
}

func (dba *Database) Select() (err error) {
	_, err = dba.Get()
	return
}
