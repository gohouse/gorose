package oci8

import (
	"database/sql"
	"fmt"
)

type RowMap struct {
	columns []string
	byName  map[string]int
	values  []interface{}
	rows    *sql.Rows
}

func NewS(r *sql.Rows, e error) (*RowMap, error) {
	if e != nil {
		return nil, e
	}
	rm := &RowMap{rows: r}
	e = rm.init()
	if e != nil {
		return nil, e
	}
	return rm, nil
}

func (r *RowMap) Key(key string) interface{} {
	if v, ok := r.Get(key); ok {
		return v
	}
	return nil

}

func (r *RowMap) Get(key string) (interface{}, bool) {
	if i, ok := r.byName[key]; ok {
		return r.values[i], true
	}
	return nil, false
}

func (r *RowMap) Row() []interface{} {
	return r.values
}

func (r *RowMap) Map() map[string]interface{} {
	ret := make(map[string]interface{}, len(r.columns))
	for i, v := range r.columns {
		ret[v] = r.values[i]
	}
	return ret
}

func (r *RowMap) Next() bool {
	return r.rows.Next()
}

func (r *RowMap) Err() error {
	return r.rows.Err()
}

func (r *RowMap) Close() error {
	r.values = r.values[:0]
	r.columns = nil
	r.byName = nil
	return r.rows.Close()
}

func (r *RowMap) Init(rows *sql.Rows) error {
	if r.columns != nil {
		panic("must be closed first")
	}
	r.rows = rows
	return r.init()
}

func (r *RowMap) init() error {
	var err error

	r.columns, err = r.rows.Columns()
	if err != nil {
		return err
	}

	delta := len(r.columns) - len(r.values)
	if delta > 0 {
		r.values = append(r.values, make([]interface{}, delta)...)
	} else if delta < 0 {
		r.values = r.values[:len(r.columns)]
	}

	r.byName = make(map[string]int, len(r.columns))
	for i, key := range r.columns {
		r.byName[key] = i
	}

	return nil
}

func (r *RowMap) Scan() error {

	for i := range r.values {
		r.values[i] = &r.values[i]
	}

	return r.rows.Scan(r.values...)
}

func (r *RowMap) Print() {
	for i, v := range r.columns {
		switch z := r.values[i].(type) {
		case string:
			if len(z) > 40 {
				fmt.Printf("%s=[%d]%v...%v.\n", v, len(z), z[:10], z[len(z)-10:])
			} else {
				fmt.Printf("%s=[%d]%v.\n", v, len(z), z)
			}
		case []byte:
			if len(z) > 40 {
				fmt.Printf("%s=[%d]%v...%v.\n", v, len(z), z[:10], z[len(z)-10:])
			} else {
				fmt.Printf("%s=[%d]%v.\n", v, len(z), z)
			}
		default:
			fmt.Printf("%s=%v.\n", v, r.values[i])
		}
	}
	fmt.Println()

}
