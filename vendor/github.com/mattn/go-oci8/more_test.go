package oci8

// ( . oracle.sh ;DSN='user:pass@:0/(description=(address_list=(address=(protocol=tcp)(host=192.168.1.1)(port=1521)))(connect_data=(sid=SID)))?isolation=SERIALIZABLE'  go test )

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

type dbc interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

var (
	db        *sql.DB
	dsnStruct *DSN
)

func DB() *sql.DB {
	if db != nil {
		return db
	}

	os.Setenv("NLS_LANG", "American_America.AL32UTF8")

	var err error
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "scott/tiger@XE"
	}

	db, err := sql.Open("oci8", dsn)
	if err != nil {
		panic(err)
	}
	dsnStruct, err = ParseDSN(dsn)
	if err != nil {
		panic(err)
	}

	db.Exec("drop table foo")
	db.Exec(sql1)

	_, err = db.Exec("truncate table foo")
	if err != nil {
		panic(err)
	}
	return db
}

func TestTruncate(t *testing.T) {
	_, err := DB().Exec("truncate table foo")
	if err != nil {
		panic(err)
	}
}

var sql1 string = `create table foo(
	c1 varchar2(256),
	c2 nvarchar2(256),
	c3 number,
	c4 float,
	c6 date,
	c7 BINARY_FLOAT,
	c8 BINARY_DOUBLE,
	c9 TIMESTAMP,
	c10 TIMESTAMP WITH TIME ZONE,
	c11 TIMESTAMP WITH LOCAL TIME ZONE,
	c12 INTERVAL YEAR TO MONTH,
	c13 INTERVAL DAY TO SECOND,
	c14 RAW(80),
	c15 ROWID,
	c17 CHAR(15),
	c18 NCHAR(20),
	c19 CLOB,
	c21 BLOB,
	cend varchar2(12)
	)`

var sql12 string = `insert( c1,c2,c3,c4,c6,c7,c8,c9,c10,c11,c12,c13,c14,c17,c18,c19,c20,c21,cend) into foo values( 
:1,
:2,
:3,
:4,
:6, 
:7,
:8,
:9, 
:10, 
:11, 
NUMTOYMINTERVAL( :12, 'MONTH'),
NUMTODSINTERVAL( :13 / 1000000000, 'SECOND'),
:14,
:17, 
:18,
:19, 
:21,
'END'
)`

func sqlstest(d dbc, t *testing.T, sql string, p ...interface{}) map[string]interface{} {

	rows, err := NewS(d.Query(sql, p...))
	if err != nil {
		t.Fatal(err)
	}
	if !rows.Next() {
		rows.Close()
		t.Fatal("no row returned:", rows.Err())
	}
	err = rows.Scan()
	if err != nil {
		rows.Close()
		t.Fatal(err)
	}
	res := rows.Map()
	//res := rows.Row()
	rows.Print()
	err = rows.Close()
	if err != nil {
		rows.Close()
		t.Fatal(err)
	}
	return res
}

func sqlstestv(d dbc, t *testing.T, sql string, p ...interface{}) []interface{} {

	rows, err := NewS(d.Query(sql, p...))
	if err != nil {
		t.Fatal(err)
	}
	if !rows.Next() {
		rows.Close()
		t.Fatal("no row returned:", rows.Err())
	}
	err = rows.Scan()
	if err != nil {
		rows.Close()
		t.Fatal(err)
	}
	//res := rows.Map()
	res := rows.Row()
	rows.Print()
	err = rows.Close()
	if err != nil {
		rows.Close()
		t.Fatal(err)
	}
	return res
}

func TestSelect1(t *testing.T) {

	//rows, err := db.Query("select :1 as AA, :2 as BB, NUMTODSINTERVAL( :3, 'SECOND') as DD, NUMTOYMINTERVAL( :4, 'MONTH') as FF, :4 as nil from dual", time.Now(), 3.14, 3.004, 55, nil)
	//rows, err := db.Query("select :1 as AA, :2 as BB, :3 as CC from dual", time.Now(), time.Now().Add( 300000000000000000), time.Now().Add( 100000000100000000))
	//rows, err := db.Query("select sysdate from dual")

	fmt.Println("bind all go types:")

	sqlstest(DB(), t,
		"select :0 as nil, :1 as true, :2 as false, :3 as int64, :4 as time, :5 as string, :6 as bytes, :7 as float64 from dual",
		nil, true, false, int64(1234567890123456789), time.Now(), "bee     ", []byte{61, 62, 63, 64, 65, 66, 67, 68}, 3.14)
}

func TestInterval1(t *testing.T) {

	fmt.Println("test interval1:")
	n := time.Duration(1234567898123456789)
	r := sqlstest(DB(), t, "select NUMTODSINTERVAL( :0 / 1000000000, 'SECOND') as intervalds from dual", int64(n))
	if n != time.Duration(r["INTERVALDS"].(int64)) {
		t.Fatal(r, "!=", n)
	}
}

func TestInterval2(t *testing.T) {

	fmt.Println("test interval2:")
	n := time.Duration(-1234567898123456789)
	r := sqlstest(DB(), t, "select NUMTODSINTERVAL( :0 / 1000000000, 'SECOND') as intervalds from dual", int64(n))
	if n != time.Duration(r["INTERVALDS"].(int64)) {
		t.Fatal(r, "!=", n)
	}
}

func TestInterval3(t *testing.T) {

	fmt.Println("test interval3:")
	n := int64(1234567890)
	r := sqlstest(DB(), t, "select NUMTOYMINTERVAL( :0, 'MONTH') as intervalym from dual", n)
	if n != r["INTERVALYM"].(int64) {
		t.Fatal(r, "!=", n)
	}
}

func TestInterval4(t *testing.T) {

	fmt.Println("test interval4:")
	n := int64(-1234567890)
	r := sqlstest(DB(), t, "select NUMTOYMINTERVAL( :0, 'MONTH') as intervalym from dual", n)
	if n != r["INTERVALYM"].(int64) {
		t.Fatal(r, "!=", n)
	}
}

func TestIntervals5(t *testing.T) {

	fmt.Println("test interval5:")

	n1 := time.Duration(987)
	n2 := time.Duration(-65)
	n3 := int64(4332)
	n4 := int64(-1239872)
	r := sqlstest(DB(), t, "select NUMTODSINTERVAL( :0 / 1000000000, 'SECOND') as i1, NUMTODSINTERVAL( :1 / 1000000000, 'SECOND') as i2, NUMTOYMINTERVAL( :2, 'MONTH') as i3, NUMTOYMINTERVAL( :3, 'MONTH') as i4 from dual", n1, n2, n3, n4)
	if n1 != time.Duration(r["I1"].(int64)) {
		t.Fatal(r["I1"], "!=", n1)
	}
	if n2 != time.Duration(r["I2"].(int64)) {
		t.Fatal(r["I2"], "!=", n2)
	}
	if n3 != r["I3"].(int64) {
		t.Fatal(r["I3"], "!=", n3)
	}
	if n4 != r["I4"].(int64) {
		t.Fatal(r["I4"], "!=", n4)
	}
}

func TestTime1(t *testing.T) {

	fmt.Println("test time1:")
	n := time.Now()
	r := sqlstest(DB(), t, "select :0 as time from dual", n)
	if !n.Equal(r["TIME"].(time.Time)) {
		t.Fatal(r, "!=", n)
	}
}

func TestTime2(t *testing.T) {
	fmt.Println("test time 2:")

	const f = "2006-01-02 15:04:05.999999999 -07:00"
	in := []time.Time{}

	tm, err := time.Parse(f, "2015-01-23 12:34:56.123456789 +09:05")
	if err != nil {
		t.Fatal(err)
	}
	in = append(in, tm)

	tm, err = time.Parse(f, "1014-10-14 21:43:50.987654321 -08:50")
	if err != nil {
		t.Fatal(err)
	}
	in = append(in, tm)

	tm = time.Date(-4123, time.Month(12), 1, 2, 3, 4, 0, time.UTC)
	in = append(in, tm)

	tm = time.Date(9321, time.Month(11), 2, 3, 4, 5, 0, time.UTC)
	in = append(in, tm)

	r := sqlstestv(DB(), t, "select :0, :1, :2, :3  from dual", in[0], in[1], in[2], in[3])
	for i, v := range r {
		vt := v.(time.Time)
		if !vt.Equal(in[i]) {
			t.Fatal(vt, "!=", in[i])
		}
	}
}

func TestTime3(t *testing.T) {
	fmt.Println("test sysdate:")
	sqlstest(DB(), t, "select sysdate - 365*6500 from dual")
}

func TestBytes1(t *testing.T) {
	fmt.Println("test bytes1:")
	n := bytes.Repeat([]byte{'A'}, 4000)
	r := sqlstest(DB(), t, "select :0 as bytes from dual", n)
	if !bytes.Equal(n, r["BYTES"].([]byte)) {
		t.Fatal(r["BYTES"], "!=", n)
	}
}

func TestBytes2(t *testing.T) {
	fmt.Println("test bytes2:")
	n := []byte{7}
	r := sqlstest(DB(), t, "select :0 as bytes from dual", n)
	if !bytes.Equal(n, r["BYTES"].([]byte)) {
		t.Fatal(r["BYTES"], "!=", n)
	}
}

func TestQuestionMark(t *testing.T) {
	if !dsnStruct.enableQMPlaceholders {
		return
	}
	fmt.Println("test question mark placeholders")
	a, b := 4, 5
	c := "zz"
	r := sqlstest(DB(), t, "select ? as v1, ? as v2, ? as v3 from dual", a, b, c)
	if fmt.Sprintf("%v", r["V1"]) != fmt.Sprintf("%v", a) {
		t.Fatal(r["V1"], "!=", a)
	}
	if fmt.Sprintf("%v", r["V2"]) != fmt.Sprintf("%v", b) {
		t.Fatal(r["V2"], "!=", b)
	}
	if fmt.Sprintf("%v", r["V3"]) != fmt.Sprintf("%v", c) {
		t.Fatal(r["V3"], "!=", c)
	}
}

func TestString1(t *testing.T) {
	fmt.Println("test string1:")
	n := strings.Repeat("1234567890", 400)
	r := sqlstest(DB(), t, "select :0 as str from dual", n)
	if n != r["STR"].(string) {
		t.Fatal(r["STR"], "!=", n)
	}
}

func TestString2(t *testing.T) {

	fmt.Println("test string2:")
	n := "6"
	r := sqlstest(DB(), t, "select :0 as str from dual", n)
	if n != r["STR"].(string) {
		t.Fatal(r["STR"], "!=", n)
	}
}

func TestString3(t *testing.T) {
	fmt.Println("test string3:")
	//n := "こんにちは 世界 Καλημέρα κόσμε こんにちは안녕하세요góðan dagGrüßgotthyvää päivääyá'át'ééhΓεια σαςВiтаюგამარჯობაनमस्ते你好здравейсвят"
	//this test depends of database charset !!!!
	n := "здравейсвят"
	r := sqlstest(DB(), t, "select :0 as str from dual", n)
	if n != r["STR"].(string) {
		t.Fatal(r["STR"], "!=", n)
	}
}

func TestFooLargeBlob(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println("test foo blob:", runtime.FuncForPC(cn).Name())
	n := make([]byte, 600000)
	for i := 0; i < len(n); i++ {
		n[i] = byte(rand.Int31n(256))
	}

	id := "idlblob"
	db := DB()
	db.Exec("insert into foo( c21, cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select c21 from foo where cend= :1", id)
	if !bytes.Equal(n, r["C21"].([]byte)) {
		t.Fatal(r["C21"], "!=", n)
	}
}

func TestSmallBlob(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	n := make([]byte, 6)
	for i := 0; i < len(n); i++ {
		n[i] = byte(rand.Int31n(256))
	}

	id := "idsblob"
	db := DB()
	db.Exec("insert into foo( c21, cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select c21 from foo where cend=:1", id)
	if !bytes.Equal(n, r["C21"].([]byte)) {
		t.Fatal(r["C21"], "!=", n)
	}
}

func TestFooRowid(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	n := "Z"
	id := "idSmallClob"
	db := DB()
	_, e := db.Exec("insert into foo( c19, cend) values( :1, :2)", n, id)
	if e != nil {
		t.Fatal(e)
	}

	sqlstest(db, t, "select rowid from foo")
}

//this test fail if transactions are readonly
func TestTransaction1(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	db, e := DB().Begin()
	if e != nil {
		t.Fatal(e)
	}

	r, e := db.Exec("insert into foo( c1) values( :1)", "123abc")
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r.RowsAffected())

	r, e = db.Exec("update foo set c1='ertertetert'")
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(r.RowsAffected())

	e = db.Commit()
	//e = db.Rollback()
	if e != nil {
		t.Fatal(e)
	}
}

func TestBigClob(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	n := "Abc" + strings.Repeat("1234567890", 2000) + "xyZ"

	id := "idBigClob"
	db := DB()
	db.Exec("insert into foo( c19, cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select c19 from foo where cend= :1", id)
	if n != r["C19"].(string) {
		println(3)
		t.Fatal(r["C19"], "!=", n)
	}
}

func TestSmallClob(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	n := "Z"
	id := "idSmallClob"
	db := DB()
	db.Exec("insert into foo( c19, cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select c19 from foo where cend= :1", id)
	if n != r["C19"].(string) {
		t.Fatal(r["C19"], "!=", n)
	}
}

func TestNvarchar(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	n := "Zкирddd"
	id := "idNvarchar"
	db := DB()
	db.Exec("insert into foo( c2, cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select c2 from foo where cend= :1", id)
	if n != r["C2"].(string) {
		t.Fatal(r["C2"], "!=", n)
	}
}

func TestNumber1(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C3"
	n := "123456.55"
	id := "idNumc3"
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	if n != r[f].(string) {
		t.Fatal(r[f], "!=", n)
	}
}

func TestNumber2(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C3"
	n := 991236.5
	id := "idNum2c3"
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	if "991236.5" != r[f].(string) {
		t.Fatal(r[f], "!=", n)
	}
}

func TestFloat1(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C4"
	n := 991236.5
	id := "idFc4"
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	if "991236.5" != r[f].(string) {
		t.Fatal(r[f], "!=", n)
	}
}

func TestBinFloat1(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C7"
	n := 1.5
	id := "idbFc7"
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	if n != r[f].(float64) {
		t.Fatal(r[f], "!=", n)
	}
}

func TestBinFloat2(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C8"
	n := 9971236.757
	id := "idbdFc8"
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	if n != r[f].(float64) {
		t.Fatal(r[f], "!=", n)
	}
}

func TestNchar(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C18"
	n := "XXкирda"
	id := "idbdNC18"
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	if strings.TrimRight(n, " ") != strings.TrimRight(r[f].(string), " ") {
		t.Fatal(r[f], "!=", n)
	}
}

func TestChar(t *testing.T) {

	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C17"
	n := "XXкирda"
	id := "idbdC17"
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	if strings.TrimRight(n, " ") != strings.TrimRight(r[f].(string), " ") {
		t.Fatal(r[f], "!=", n)
	}
}

func TestDate(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)

	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C6"

	const fm = "2006-01-02 15:04:05.999999999 -07:00"

	n, err := time.Parse(fm, "2014-10-23 04:56:12.123456789 +09:06")
	if err != nil {
		t.Fatal(err)
	}

	id := "idbdate" + f
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	fmt.Println(n, r[f].(time.Time))
}

func TestTimestamp(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C9"

	const fm = "2006-01-02 15:04:05.999999999 -07:00"

	n, err := time.Parse(fm, "2014-10-23 04:56:12.123456789 +09:06")
	if err != nil {
		t.Fatal(err)
	}

	id := "idTstamp" + f
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	fmt.Println(n, r[f].(time.Time))
}

func TestTimestampTz(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C10"

	const fm = "2006-01-02 15:04:05.999999999 -07:00"

	n, err := time.Parse(fm, "2014-10-23 04:56:12.123456789 +09:06")
	if err != nil {
		t.Fatal(err)
	}

	id := "idTs" + f
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	fmt.Println(n, r[f].(time.Time))
}

func TestTimestampLtz(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())
	f := "C11"

	const fm = "2006-01-02 15:04:05.999999999 -07:00"

	n, err := time.Parse(fm, "2014-10-23 04:56:12.123456000 +09:06")
	if err != nil {
		t.Fatal(err)
	}

	id := "idTs" + f
	db := DB()
	db.Exec("insert into foo( "+f+", cend) values( :1, :2)", n, id)

	r := sqlstest(db, t, "select "+f+" from foo where cend= :1", id)
	fmt.Println(n, r[f].(time.Time), "equal ?", n.Equal(r[f].(time.Time)))
}

func TestQueryRowPrepared(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	sel, err := DB().Prepare("select :1 from dual")
	if err != nil {
		t.Fatal(err)
	}

	const val = 143
	ccc := val

	err = sel.QueryRow(ccc).Scan(&ccc)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ccc)
	if ccc != val {
		t.Fatal(err)
	}

	err = sel.QueryRow(ccc).Scan(&ccc)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ccc)
	sel.Close()
}

func TestTimeZones(t *testing.T) {
	zones := getZones()
	db := DB()
	seen := make(map[string]bool, len(zones)*2)
	for _, z0 := range zones {
		loc, err := time.LoadLocation(z0)
		z1 := loc.String() // for me z0 is always == z1

		if err != nil {
			continue
		}
		tt := time.Date(2015, 12, 31, 23, 59, 59, 123456789, loc)
		z2, _ := tt.Zone() // sometimes z1 != z2 e.g. EST5EDT v.s. EDT

		if _, exists := seen[z0]; !exists {
			seen[z0] = true
			handleZone(z0, &tt, db, t)
		}
		if _, exists := seen[z1]; !exists {
			seen[z1] = true
			handleZone(z1, &tt, db, t)
		}
		if _, exists := seen[z2]; !exists {
			seen[z2] = true
			handleZone(z2, &tt, db, t)
		}
	}
}

func handleZone(zone string, tt *time.Time, db *sql.DB, t *testing.T) {
	r := sqlstest(db, t, "select :0 as time from dual", *tt)
	if !tt.Equal(r["TIME"].(time.Time)) {
		t.Fatal(r, "!=", tt)
	}
}

func getZones() []string {
	// mkdir foo; cd foo; unzip $GOROOT/lib/time/zoneinfo.zip; find . -type f | sed 's/^../"/;s/$/",/' | grep -v '^"Factory",$' | sort -u | fmt -w 80 > ../zone.names; cd ..; rm -rf foo
	return []string{
		"Africa/Abidjan", "Africa/Accra", "Africa/Addis_Ababa", "Africa/Algiers",
		"Africa/Asmara", "Africa/Asmera", "Africa/Bamako", "Africa/Bangui",
		"Africa/Banjul", "Africa/Bissau", "Africa/Blantyre", "Africa/Brazzaville",
		"Africa/Bujumbura", "Africa/Cairo", "Africa/Casablanca", "Africa/Ceuta",
		"Africa/Conakry", "Africa/Dakar", "Africa/Dar_es_Salaam", "Africa/Djibouti",
		"Africa/Douala", "Africa/El_Aaiun", "Africa/Freetown", "Africa/Gaborone",
		"Africa/Harare", "Africa/Johannesburg", "Africa/Juba", "Africa/Kampala",
		"Africa/Khartoum", "Africa/Kigali", "Africa/Kinshasa", "Africa/Lagos",
		"Africa/Libreville", "Africa/Lome", "Africa/Luanda", "Africa/Lubumbashi",
		"Africa/Lusaka", "Africa/Malabo", "Africa/Maputo", "Africa/Maseru",
		"Africa/Mbabane", "Africa/Mogadishu", "Africa/Monrovia", "Africa/Nairobi",
		"Africa/Ndjamena", "Africa/Niamey", "Africa/Nouakchott", "Africa/Ouagadougou",
		"Africa/Porto-Novo", "Africa/Sao_Tome", "Africa/Timbuktu", "Africa/Tripoli",
		"Africa/Tunis", "Africa/Windhoek", "America/Adak", "America/Anchorage",
		"America/Anguilla", "America/Antigua", "America/Araguaina",
		"America/Argentina/Buenos_Aires", "America/Argentina/Catamarca",
		"America/Argentina/ComodRivadavia", "America/Argentina/Cordoba",
		"America/Argentina/Jujuy", "America/Argentina/La_Rioja",
		"America/Argentina/Mendoza", "America/Argentina/Rio_Gallegos",
		"America/Argentina/Salta", "America/Argentina/San_Juan",
		"America/Argentina/San_Luis", "America/Argentina/Tucuman",
		"America/Argentina/Ushuaia", "America/Aruba", "America/Asuncion",
		"America/Atikokan", "America/Atka", "America/Bahia", "America/Bahia_Banderas",
		"America/Barbados", "America/Belem", "America/Belize", "America/Blanc-Sablon",
		"America/Boa_Vista", "America/Bogota", "America/Boise", "America/Buenos_Aires",
		"America/Cambridge_Bay", "America/Campo_Grande", "America/Cancun",
		"America/Caracas", "America/Catamarca", "America/Cayenne", "America/Cayman",
		"America/Chicago", "America/Chihuahua", "America/Coral_Harbour",
		"America/Cordoba", "America/Costa_Rica", "America/Creston",
		"America/Cuiaba", "America/Curacao", "America/Danmarkshavn",
		"America/Dawson", "America/Dawson_Creek", "America/Denver",
		"America/Detroit", "America/Dominica", "America/Edmonton",
		"America/Eirunepe", "America/El_Salvador", "America/Ensenada",
		"America/Fortaleza", "America/Fort_Wayne", "America/Glace_Bay",
		"America/Godthab", "America/Goose_Bay", "America/Grand_Turk",
		"America/Grenada", "America/Guadeloupe", "America/Guatemala",
		"America/Guayaquil", "America/Guyana", "America/Halifax", "America/Havana",
		"America/Hermosillo", "America/Indiana/Indianapolis", "America/Indiana/Knox",
		"America/Indiana/Marengo", "America/Indiana/Petersburg",
		"America/Indianapolis", "America/Indiana/Tell_City", "America/Indiana/Vevay",
		"America/Indiana/Vincennes", "America/Indiana/Winamac", "America/Inuvik",
		"America/Iqaluit", "America/Jamaica", "America/Jujuy", "America/Juneau",
		"America/Kentucky/Louisville", "America/Kentucky/Monticello",
		"America/Knox_IN", "America/Kralendijk", "America/La_Paz", "America/Lima",
		"America/Los_Angeles", "America/Louisville", "America/Lower_Princes",
		"America/Maceio", "America/Managua", "America/Manaus", "America/Marigot",
		"America/Martinique", "America/Matamoros", "America/Mazatlan",
		"America/Mendoza", "America/Menominee", "America/Merida", "America/Metlakatla",
		"America/Mexico_City", "America/Miquelon", "America/Moncton",
		"America/Monterrey", "America/Montevideo", "America/Montreal",
		"America/Montserrat", "America/Nassau", "America/New_York", "America/Nipigon",
		"America/Nome", "America/Noronha", "America/North_Dakota/Beulah",
		"America/North_Dakota/Center", "America/North_Dakota/New_Salem",
		"America/Ojinaga", "America/Panama", "America/Pangnirtung",
		"America/Paramaribo", "America/Phoenix", "America/Port-au-Prince",
		"America/Porto_Acre", "America/Port_of_Spain", "America/Porto_Velho",
		"America/Puerto_Rico", "America/Rainy_River", "America/Rankin_Inlet",
		"America/Recife", "America/Regina", "America/Resolute", "America/Rio_Branco",
		"America/Rosario", "America/Santa_Isabel", "America/Santarem",
		"America/Santiago", "America/Santo_Domingo", "America/Sao_Paulo",
		"America/Scoresbysund", "America/Shiprock", "America/Sitka",
		"America/St_Barthelemy", "America/St_Johns", "America/St_Kitts",
		"America/St_Lucia", "America/St_Thomas", "America/St_Vincent",
		"America/Swift_Current", "America/Tegucigalpa", "America/Thule",
		"America/Thunder_Bay", "America/Tijuana", "America/Toronto",
		"America/Tortola", "America/Vancouver", "America/Virgin",
		"America/Whitehorse", "America/Winnipeg", "America/Yakutat",
		"America/Yellowknife", "Antarctica/Casey", "Antarctica/Davis",
		"Antarctica/DumontDUrville", "Antarctica/Macquarie", "Antarctica/Mawson",
		"Antarctica/McMurdo", "Antarctica/Palmer", "Antarctica/Rothera",
		"Antarctica/South_Pole", "Antarctica/Syowa", "Antarctica/Troll",
		"Antarctica/Vostok", "Arctic/Longyearbyen", "Asia/Aden", "Asia/Almaty",
		"Asia/Amman", "Asia/Anadyr", "Asia/Aqtau", "Asia/Aqtobe", "Asia/Ashgabat",
		"Asia/Ashkhabad", "Asia/Baghdad", "Asia/Bahrain", "Asia/Baku", "Asia/Bangkok",
		"Asia/Beirut", "Asia/Bishkek", "Asia/Brunei", "Asia/Calcutta", "Asia/Chita",
		"Asia/Choibalsan", "Asia/Chongqing", "Asia/Chungking", "Asia/Colombo",
		"Asia/Dacca", "Asia/Damascus", "Asia/Dhaka", "Asia/Dili", "Asia/Dubai",
		"Asia/Dushanbe", "Asia/Gaza", "Asia/Harbin", "Asia/Hebron", "Asia/Ho_Chi_Minh",
		"Asia/Hong_Kong", "Asia/Hovd", "Asia/Irkutsk", "Asia/Istanbul",
		"Asia/Jakarta", "Asia/Jayapura", "Asia/Jerusalem", "Asia/Kabul",
		"Asia/Kamchatka", "Asia/Karachi", "Asia/Kashgar", "Asia/Kathmandu",
		"Asia/Katmandu", "Asia/Khandyga", "Asia/Kolkata", "Asia/Krasnoyarsk",
		"Asia/Kuala_Lumpur", "Asia/Kuching", "Asia/Kuwait", "Asia/Macao", "Asia/Macau",
		"Asia/Magadan", "Asia/Makassar", "Asia/Manila", "Asia/Muscat", "Asia/Nicosia",
		"Asia/Novokuznetsk", "Asia/Novosibirsk", "Asia/Omsk", "Asia/Oral",
		"Asia/Phnom_Penh", "Asia/Pontianak", "Asia/Pyongyang", "Asia/Qatar",
		"Asia/Qyzylorda", "Asia/Rangoon", "Asia/Riyadh", "Asia/Saigon",
		"Asia/Sakhalin", "Asia/Samarkand", "Asia/Seoul", "Asia/Shanghai",
		"Asia/Singapore", "Asia/Srednekolymsk", "Asia/Taipei", "Asia/Tashkent",
		"Asia/Tbilisi", "Asia/Tehran", "Asia/Tel_Aviv", "Asia/Thimbu", "Asia/Thimphu",
		"Asia/Tokyo", "Asia/Ujung_Pandang", "Asia/Ulaanbaatar", "Asia/Ulan_Bator",
		"Asia/Urumqi", "Asia/Ust-Nera", "Asia/Vientiane", "Asia/Vladivostok",
		"Asia/Yakutsk", "Asia/Yekaterinburg", "Asia/Yerevan", "Atlantic/Azores",
		"Atlantic/Bermuda", "Atlantic/Canary", "Atlantic/Cape_Verde",
		"Atlantic/Faeroe", "Atlantic/Faroe", "Atlantic/Jan_Mayen", "Atlantic/Madeira",
		"Atlantic/Reykjavik", "Atlantic/South_Georgia", "Atlantic/Stanley",
		"Atlantic/St_Helena", "Australia/ACT", "Australia/Adelaide",
		"Australia/Brisbane", "Australia/Broken_Hill", "Australia/Canberra",
		"Australia/Currie", "Australia/Darwin", "Australia/Eucla", "Australia/Hobart",
		"Australia/LHI", "Australia/Lindeman", "Australia/Lord_Howe",
		"Australia/Melbourne", "Australia/North", "Australia/NSW", "Australia/Perth",
		"Australia/Queensland", "Australia/South", "Australia/Sydney",
		"Australia/Tasmania", "Australia/Victoria", "Australia/West",
		"Australia/Yancowinna", "Brazil/Acre", "Brazil/DeNoronha", "Brazil/East",
		"Brazil/West", "Canada/Atlantic", "Canada/Central", "Canada/Eastern",
		"Canada/East-Saskatchewan", "Canada/Mountain", "Canada/Newfoundland",
		"Canada/Pacific", "Canada/Saskatchewan", "Canada/Yukon", "CET",
		"Chile/Continental", "Chile/EasterIsland", "CST6CDT", "Cuba", "EET",
		"Egypt", "Eire", "EST", "EST5EDT", "Etc/GMT", "Etc/GMT-0", "Etc/GMT+0",
		"Etc/GMT0", "Etc/GMT-1", "Etc/GMT+1", "Etc/GMT-10", "Etc/GMT+10",
		"Etc/GMT-11", "Etc/GMT+11", "Etc/GMT-12", "Etc/GMT+12", "Etc/GMT-13",
		"Etc/GMT-14", "Etc/GMT-2", "Etc/GMT+2", "Etc/GMT-3", "Etc/GMT+3", "Etc/GMT-4",
		"Etc/GMT+4", "Etc/GMT-5", "Etc/GMT+5", "Etc/GMT-6", "Etc/GMT+6", "Etc/GMT-7",
		"Etc/GMT+7", "Etc/GMT-8", "Etc/GMT+8", "Etc/GMT-9", "Etc/GMT+9",
		"Etc/Greenwich", "Etc/UCT", "Etc/Universal", "Etc/UTC", "Etc/Zulu",
		"Europe/Amsterdam", "Europe/Andorra", "Europe/Athens", "Europe/Belfast",
		"Europe/Belgrade", "Europe/Berlin", "Europe/Bratislava", "Europe/Brussels",
		"Europe/Bucharest", "Europe/Budapest", "Europe/Busingen", "Europe/Chisinau",
		"Europe/Copenhagen", "Europe/Dublin", "Europe/Gibraltar", "Europe/Guernsey",
		"Europe/Helsinki", "Europe/Isle_of_Man", "Europe/Istanbul", "Europe/Jersey",
		"Europe/Kaliningrad", "Europe/Kiev", "Europe/Lisbon", "Europe/Ljubljana",
		"Europe/London", "Europe/Luxembourg", "Europe/Madrid", "Europe/Malta",
		"Europe/Mariehamn", "Europe/Minsk", "Europe/Monaco", "Europe/Moscow",
		"Europe/Nicosia", "Europe/Oslo", "Europe/Paris", "Europe/Podgorica",
		"Europe/Prague", "Europe/Riga", "Europe/Rome", "Europe/Samara",
		"Europe/San_Marino", "Europe/Sarajevo", "Europe/Simferopol", "Europe/Skopje",
		"Europe/Sofia", "Europe/Stockholm", "Europe/Tallinn", "Europe/Tirane",
		"Europe/Tiraspol", "Europe/Uzhgorod", "Europe/Vaduz", "Europe/Vatican",
		"Europe/Vienna", "Europe/Vilnius", "Europe/Volgograd", "Europe/Warsaw",
		"Europe/Zagreb", "Europe/Zaporozhye", "Europe/Zurich", "GB", "GB-Eire",
		"GMT", "GMT-0", "GMT+0", "GMT0", "Greenwich", "Hongkong", "HST", "Iceland",
		"Indian/Antananarivo", "Indian/Chagos", "Indian/Christmas", "Indian/Cocos",
		"Indian/Comoro", "Indian/Kerguelen", "Indian/Mahe", "Indian/Maldives",
		"Indian/Mauritius", "Indian/Mayotte", "Indian/Reunion", "Iran", "Israel",
		"Jamaica", "Japan", "Kwajalein", "Libya", "MET", "Mexico/BajaNorte",
		"Mexico/BajaSur", "Mexico/General", "MST", "MST7MDT", "Navajo", "NZ",
		"NZ-CHAT", "Pacific/Apia", "Pacific/Auckland", "Pacific/Bougainville",
		"Pacific/Chatham", "Pacific/Chuuk", "Pacific/Easter", "Pacific/Efate",
		"Pacific/Enderbury", "Pacific/Fakaofo", "Pacific/Fiji", "Pacific/Funafuti",
		"Pacific/Galapagos", "Pacific/Gambier", "Pacific/Guadalcanal", "Pacific/Guam",
		"Pacific/Honolulu", "Pacific/Johnston", "Pacific/Kiritimati", "Pacific/Kosrae",
		"Pacific/Kwajalein", "Pacific/Majuro", "Pacific/Marquesas", "Pacific/Midway",
		"Pacific/Nauru", "Pacific/Niue", "Pacific/Norfolk", "Pacific/Noumea",
		"Pacific/Pago_Pago", "Pacific/Palau", "Pacific/Pitcairn", "Pacific/Pohnpei",
		"Pacific/Ponape", "Pacific/Port_Moresby", "Pacific/Rarotonga",
		"Pacific/Saipan", "Pacific/Samoa", "Pacific/Tahiti", "Pacific/Tarawa",
		"Pacific/Tongatapu", "Pacific/Truk", "Pacific/Wake", "Pacific/Wallis",
		"Pacific/Yap", "Poland", "Portugal", "PRC", "PST8PDT", "ROC", "ROK",
		"Singapore", "Turkey", "UCT", "Universal", "US/Alaska", "US/Aleutian",
		"US/Arizona", "US/Central", "US/Eastern", "US/East-Indiana", "US/Hawaii",
		"US/Indiana-Starke", "US/Michigan", "US/Mountain", "US/Pacific",
		"US/Pacific-New", "US/Samoa", "UTC", "WET", "W-SU", "Zulu",
	}
}

//watch mem in top :)    I wish valgrind can run go progs...
//warn 5 min test !!!!
func zzTestMem(t *testing.T) {
	cn, _, _, _ := runtime.Caller(0)
	fmt.Println(runtime.FuncForPC(cn).Name())

	for now := time.Now().Add(time.Minute * 5); now.After(time.Now()); {
		TestTruncate(t)
		TestSelect1(t)
		TestInterval1(t)
		TestInterval2(t)
		TestInterval3(t)
		TestInterval4(t)
		TestIntervals5(t)
		TestTime1(t)
		TestTime2(t)
		TestTime3(t)
		TestBytes1(t)
		TestBytes2(t)
		TestString1(t)
		TestString2(t)
		TestString3(t)
		TestFooLargeBlob(t)
		TestSmallBlob(t)
		TestFooRowid(t)
		TestTransaction1(t)
		TestBigClob(t)
		TestSmallClob(t)
		TestNvarchar(t)
		TestNumber1(t)
		TestNumber2(t)
		TestFloat1(t)
		TestBinFloat1(t)
		TestBinFloat2(t)
		TestNchar(t)
		TestChar(t)
		TestDate(t)
		TestTimestamp(t)
		TestTimestampTz(t)
		TestTimestampLtz(t)
	}
}
