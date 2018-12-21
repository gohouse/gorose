package builder

import (
	"errors"
	"fmt"
	"github.com/gohouse/gorose/across"
	"github.com/gohouse/gorose/utils"
	"strconv"
	"strings"
)

type MysqlBuilder struct {
}

func init() {
	// 检查解析器是否实现了接口
	var builderTmp IBuilder = &MysqlBuilder{}

	// 注册驱动
	Register("mysql", builderTmp)
}

func (m *MysqlBuilder) BuildQuery(ormApi across.OrmApi) (sql string, err error) {
	// distinct
	distinct := utils.If(ormApi.Sdistinct, "DISTINCT ", "")
	// fields
	fields := strings.Join(ormApi.Sfields, ", ")
	if ormApi.Sunion != "" {
		fields = ormApi.Sunion
	} else if fields == "" {
		fields = "*"
	}
	// table
	table := ormApi.Prefix + ormApi.TableName
	// join
	parseJoin, err := parseJoin(ormApi)
	if err != nil {
		return "", err
	}
	join := parseJoin
	// where
	ormApi.SbeforeParseWhereData = ormApi.Swhere
	parseWhere, err := parseWhere(ormApi)
	if err != nil {
		return "", err
	}
	where := utils.If(parseWhere == "", "", " WHERE "+parseWhere).(string)
	// group
	group := utils.If(ormApi.Sgroup == "", "", " GROUP BY "+ormApi.Sgroup).(string)
	// having
	having := utils.If(ormApi.Shaving == "", "", " HAVING "+ormApi.Shaving).(string)
	// order
	order := utils.If(ormApi.Sorder == "", "", " ORDER BY "+ormApi.Sorder).(string)
	// limit
	limit := utils.If(ormApi.Slimit == 0, "", " LIMIT "+strconv.Itoa(ormApi.Slimit))
	// offset
	offset := utils.If(ormApi.Soffset == 0, "", " OFFSET "+strconv.Itoa(ormApi.Soffset))
	//sqlstr := "select " + distinct + fields + " from " + table + where + group + having + order + limit + offset
	sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
		distinct, fields, table, join, where, group, having, order, limit, offset)

	return sqlstr, nil
}

// BuildExecut : build execute query string
func (m *MysqlBuilder) BuildExecute(ormApi across.OrmApi, operType string) (sql string, err error) {
	// insert : {"name":"fizz, "website":"fizzday.net"} or {{"name":"fizz2", "website":"www.fizzday.net"}, {"name":"fizz", "website":"fizzday.net"}}}
	// update : {"name":"fizz", "website":"fizzday.net"}
	// delete : ...
	var update, insertkey, insertval, sqlstr string
	if operType != "delete" {
		update, insertkey, insertval = buildData(ormApi)
	}

	ormApi.SbeforeParseWhereData = ormApi.Swhere
	res, err := parseWhere(ormApi)
	if err != nil {
		return res, err
	}
	where := utils.If(res == "", "", " WHERE "+res).(string)

	tableName := ormApi.Prefix + ormApi.TableName
	switch operType {
	case "insert":
		sqlstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, insertkey, insertval)
	case "update":
		if res=="" && ormApi.Sforce==false{
			return sqlstr, errors.New("出于安全考虑, update时where条件不能为空, 如果真的不需要where条件, 请使用force(如: db.xxx.Force().Update())")
		}
		sqlstr = fmt.Sprintf("UPDATE %s SET %s%s", tableName, update, where)
	case "delete":
		if res=="" && ormApi.Sforce==false{
			return sqlstr, errors.New("出于安全考虑, delete时where条件不能为空, 如果真的不需要where条件, 请使用force(如: db.xxx.Force().Delete())")
		}
		sqlstr = fmt.Sprintf("DELETE FROM %s%s", tableName, where)
	}
	//fmt.Println(sqlstr)
	//dba.Reset()

	return sqlstr, nil
}
