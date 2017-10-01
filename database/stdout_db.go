package database

import (
	"fmt"
)

type StdoutDatabase struct{}

func (db *StdoutDatabase) Write(r *Record) {
	fmt.Println(r)
}
