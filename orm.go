package gorose

import "fmt"

type Orm struct {
	IBuilder
	*bind
	where                [][]interface{} // where
	order                string          // order
	offset               int             // offset
	join                 [][]interface{} // join
	distinct             bool            // distinct
	union                string          // sum/count/avg/max/min
	group                string          // group
	having               string          // having
	data                 interface{}     // data
	//beforeParseWhereData [][]interface{}
}

func NewOrm(b *bind) IOrm {
	var o = new(Orm)
	o.bind = b

	return o
}

func (o Orm) Hello()  {
	fmt.Println("hello gorose orm struct")
}

func (o Orm) BuildSql(operType string) (string, []interface{}, error) {
	if operType == "select" {
		return o.BuildQuery(o)
	} else {
		return o.BuildExecute(o,operType)
	}
}
