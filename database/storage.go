package database

import (
	"context"
	"fmt"
	"sync"
)

func NewAsyncStorage(db Database) *AsyncStorage {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s := AsyncStorage{
		pending:    make(chan *Record, 100),
		ctx:        ctx,
		cancelFunc: cancelFunc,
		db:         db,
		cond:       sync.NewCond(&sync.Mutex{}),
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
	s.cond.L.Lock()
	fmt.Println("enqueue, len q", len(s.pending))
	s.pending <- r
	s.cond.L.Unlock()
	fmt.Println("notifying, len q", len(s.pending))
	s.cond.Broadcast()
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
	s.ctx, s.cancelFunc = context.WithCancel(s.ctx)
	go func() {
		for {
			s.cond.L.Lock()
			for len(s.pending) == 0 {
				s.cond.Wait()
			}
			select {
			case <-s.ctx.Done():
				return
			case r := <-s.pending:
				s.db.Write(r)
			}
			s.cond.L.Unlock()
		}
	}()
}

func (s *AsyncStorage) StopFlush() {
	s.cancelFunc()
}
