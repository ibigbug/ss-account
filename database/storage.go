package database

import "context"
import "fmt"

type Record struct {
	Username  string
	BytesUsed int
	Dir       int
	Time      int64
}

func (r *Record) String() string {
	return fmt.Sprintf("%s:%v:%d:%d", r.Username, r.Dir, r.BytesUsed, r.Time)
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
		pending:    make(chan *Record, 100),
		ctx:        ctx,
		cancelFunc: cancelFunc,
		db:         &StdoutDatabase{},
	}
	//	s.StartFlush()
	return &s
}

type AsyncStorage struct {
	pending    chan *Record
	db         Database
	ctx        context.Context
	cancelFunc func()
}

func (s *AsyncStorage) Write(r *Record) {
	//	s.pending <- r
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
