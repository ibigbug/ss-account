package database

import (
	"context"
	"sync"
)

func NewAsyncStorage(db Database) *AsyncStorage {
	s := AsyncStorage{
		pending: make(chan *Record, 100),
		db:      db,
	}
	s.StartFlush()
	return &s
}

type AsyncStorage struct {
	pending    chan *Record
	db         Database
	cond       *sync.Cond
	ctx        context.Context
	cancelFunc func()
}

func (s *AsyncStorage) Write(r *Record) error {
	s.pending <- r
	return nil
}

func (s *AsyncStorage) BindPort(b *Binding) error {
	return s.db.BindPort(b)
}

func (s *AsyncStorage) GetAllActiveBinding() ([]*Binding, error) {
	return s.db.GetAllActiveBinding()
}

func (s *AsyncStorage) GetAllUserUsage() ([]*UserUsage, error) {
	return s.db.GetAllUserUsage()
}

// StartFlush not goroutine safe
func (s *AsyncStorage) StartFlush() {
	for i := 0; i < 10; i++ {
		go func() {
			for r := range s.pending {
				s.db.Write(r)
			}
		}()
	}
}

func (s *AsyncStorage) StopFlush() {
	close(s.pending)
}
