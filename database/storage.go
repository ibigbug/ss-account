package database

import "context"

func NewAsyncStorage(db Database) *AsyncStorage {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s := AsyncStorage{
		pending:    make(chan *Record, 100),
		ctx:        ctx,
		cancelFunc: cancelFunc,
		db:         db,
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

func (s *AsyncStorage) BindPort(b *Binding) {
	s.db.BindPort(b)
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
