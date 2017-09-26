package server

import "context"

type Record struct {
	Username  string
	BytesUsed int
	Dir       Direction
	Time      int64
}
type Storage interface {
	Write(r *Record)
	Flush()
}

type Database interface {
	Write(r *Record)
}

func NewAsyncStorage() *AsyncStorage {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s := AsyncStorage{
		pending:    make(chan *Record),
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
	return &s
}

type AsyncStorage struct {
	pending    chan *Record
	db         Database
	ctx        context.Context
	cancelFunc func()
}

func (s *AsyncStorage) Write(r *Record) {
	s.pending <- r
}

// StartFlush not goroutine safe
func (s *AsyncStorage) StartFlush() {
	s.ctx, s.cancelFunc = context.WithCancel(s.ctx)
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case r := <-s.pending:
				s.db.Write(r)
			default:
			}

		}
	}()
}
