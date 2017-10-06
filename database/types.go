package database

import "fmt"

type Record struct {
	Username  string
	BytesUsed int
	Dir       int // up -> 0, down -> 1
	Time      int64
}

func (r *Record) String() string {
	return fmt.Sprintf("%s:%v:%d:%d", r.Username, r.Dir, r.BytesUsed, r.Time)
}

type Binding struct {
	Username string
	Port     string
	Active   bool
}

func (b *Binding) String() string {
	return fmt.Sprintf("%s <-> %s, %v", b.Username, b.Port, b.Active)
}

// Storage provides a middle layer that
// support/manage async write in case we
// use db drivers which doens't support
// async write
type Storage interface {
	Write(r *Record) error
	BindPort(b *Binding) error
	GetAllActiveBinding() ([]*Binding, error)
}

// Database provides the write options
// beyond the underlying db-driver(mysql,redis,etc.)
type Database interface {
	Write(r *Record) error
	BindPort(b *Binding) error
	GetAllActiveBinding() ([]*Binding, error)
}
