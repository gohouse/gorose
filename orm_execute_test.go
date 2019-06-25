package gorose

import (
	"fmt"
	"testing"
)

func TestOrm_Update(t *testing.T) {
	db := DB()

	var u = Users{
		Name: "gorose2",
		Age:  19,
	}

	aff, err := db.Force().Data(u).Update()
	fmt.Println(aff, err)
	fmt.Println(db.LastSql())
}
