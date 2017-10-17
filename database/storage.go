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
	return &s
}

type AsyncStorage struct {
	pending    chan *Record
	db         Database
	ctx        context.Context
	cancelFunc func()
}

func (s *AsyncStorage) Write(r *Record) error {
	go s.db.Write(r)
	return nil
}

func (s *AsyncStorage) BindPort(b *Binding) error {
	return s.db.BindPort(b)
}

func (s *AsyncStorage) GetAllActiveBinding() ([]*Binding, error) {
	return s.db.GetAllActiveBinding()
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
				s.StopFlush()
			}
		}
	}()
}

func (s *AsyncStorage) StopFlush() {
	s.cancelFunc()
}
