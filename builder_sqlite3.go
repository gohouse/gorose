package gorose

type BuilderSqlite3 struct {
	IOrm
}

func init()  {
	NewBuilderDriver().Register("sqlite3", &BuilderSqlite3{})
}

func (bs *BuilderSqlite3) BuildQuery(o IOrm) (sqlStr string, args []interface{}, err error) {
	sqlStr = "selct xxx test sql"
	return
}

func (bs *BuilderSqlite3) BuildExecute(o IOrm,operType string) (sqlStr string, args []interface{}, err error) {
	sqlStr = "execute xxx test sql"
	return
}
