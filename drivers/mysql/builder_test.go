package mysql

import (
	"github.com/gohouse/gorose/v3"
	"testing"
)

type User struct {
	Id   int64  `db:"id,pk"`
	Name string `db:"name"`
}

var dbg = gorose.Open(nil)

func db() *gorose.Database {
	return dbg.NewDatabase()
}

func TestDatabase_ToSqlTo(t *testing.T) {
	var user = User{Id: 1}
	prepare, values, err := db().ToSqlTo(&user)
	assertsError(t, err)
	var expect = "SELECT `id`, `name` FROM `User` WHERE `id` = ? LIMIT ?"
	assertsEqual(t, expect, prepare)
	var expectValues = []int{1, 1}
	assertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlToSlice(t *testing.T) {
	var user []User
	prepare, values, err := db().Where("id", ">", 1).OrderBy("id").Limit(10).Page(2).ToSqlTo(&user)
	assertsError(t, err)
	var expect = "SELECT `id`, `name` FROM `User` WHERE `id` > ? ORDER BY `id` LIMIT ? OFFSET ?"
	assertsEqual(t, expect, prepare)
	var expectValues = []int{1, 10, 10}
	assertsEqual(t, expectValues, values)
}
func TestDatabase_ToSql(t *testing.T) {
	prepare, values, err := db().Table("users").Select("b").Where("c", 1).GroupBy("a").Having("a", 1).OrderBy("id").Limit(10).Page(2).ToSql()
	assertsError(t, err)
	var expect = "SELECT `b` FROM `users` WHERE `c` = ? GROUP BY `a` HAVING `a` = ? ORDER BY `id` LIMIT ? OFFSET ?"
	assertsEqual(t, expect, prepare)
	var expectValues = []int{1, 10, 10, 1}
	assertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlInsert(t *testing.T) {
	var user = User{Name: "john"}
	prepare, values, err := db().ToSqlInsert(&user)
	assertsError(t, err)
	var expect = "INSERT INTO `User` (`name`) VALUES (?)"
	assertsEqual(t, expect, prepare)
	var expectValues = []string{"john"}
	assertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlInserts(t *testing.T) {
	var user = []User{{Name: "John"}, {Name: "Alice"}}
	prepare, values, err := db().ToSqlInsert(&user)
	assertsError(t, err)
	var expect = "INSERT INTO `User` (`name`) VALUES (?),(?)"
	assertsEqual(t, expect, prepare)
	var expectValues = []string{"John", "Alice"}
	assertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlUpdate(t *testing.T) {
	var user = User{Id: 1, Name: "john"}
	prepare, values, err := db().ToSqlUpdate(&user)
	assertsError(t, err)
	var expect = "UPDATE `User` SET `name` = ? WHERE `id` = ?"
	assertsEqual(t, expect, prepare)
	var expectValues = []any{"john", 1}
	assertsEqual(t, expectValues, values)
}
func TestDatabase_ToSqlDelete(t *testing.T) {
	var user = User{Id: 1}
	prepare, values, err := db().ToSqlDelete(&user)
	assertsError(t, err)
	var expect = "DELETE FROM `User` WHERE `id` = ?"
	assertsEqual(t, expect, prepare)
	var expectValues = []any{1}
	assertsEqual(t, expectValues, values)
}
