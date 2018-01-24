// +build go1.8

package oci8

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"
)

func TestNamedParam(t *testing.T) {
	r := sqlstest(DB(), t, "select :foo||:bar as message from dual", sql.Named("foo", "hello"), sql.Named("bar", "world"))
	if "helloworld" != r["MESSAGE"].(string) {
		t.Fatal("message should be: helloworld", r)
	}
}

/* FIXME
func TestOutputBind(t *testing.T) {
	db := DB()

	s1 := "-----------------------------"
	s2 := 11
	s3 := false
	_, err := db.Exec(`begin  :a := 42; :b := 'ddddd' ; :c := 2; end;`,
		sql.Named("a", sql.Out{Dest: &s2}),
		sql.Named("b", sql.Out{Dest: &s1}),
		sql.Named("c", sql.Out{Dest: &s3}))
	if err != nil {
		t.Fatal(err)
	}
	s1want := "ddddd                        "
	if s1 != s1want {
		t.Fatalf("want %q but %q", s1want, s1)
	}
	if s2 != 42 {
		t.Fatalf("want %v but %v", 42, s2)
	}
	if !s3 {
		t.Fatalf("want %v but %v", true, s3)
	}
}
*/

func TestTimeout(t *testing.T) {
	db := DB()
	for i := 0; i < 2000; i++ {
		db.Exec("insert into foo(c3) values(:1)", i)
	}

	var wg sync.WaitGroup
	f := func(wg *sync.WaitGroup) {
		defer wg.Done()

		stmt, err := db.Prepare(`select * from foo order by c3`)
		if err != nil {
			t.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 200*time.Millisecond)
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		err = ctx.Err()
		if err != nil && err != context.DeadlineExceeded {
			t.Fatal(err)
		}
	}
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go f(&wg)
	}
	wg.Wait()
}
