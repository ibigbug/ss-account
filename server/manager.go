package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Direction int

const (
	DirectionUpload Direction = iota
	DirectionDownload
)

// DefaultManaged is default user manager
var DefaultManaged = &Managed{}
var DefaultStorage = NewAsyncStorage()

// GetManagerByUsername ...
func GetManagerByUsername(username string) *Manager {
	for _, m := range DefaultManaged.mrs {
		if m.Username == username {
			return m
		}
	}
	return nil
}

// Managed list of managers
type Managed struct {
	sync.Mutex
	mrs []*Manager
}

func (m *Managed) String() string {
	var s []string
	for _, mr := range m.mrs {
		s = append(s, mr.String())
	}
	return strings.Join(s, "\n")
}

// Add ...
func (m *Managed) Add(manager *Manager) {
	m.Lock()
	defer m.Unlock()
	m.mrs = append(m.mrs, manager)
}

// Remove ...
func (m *Managed) Remove(manager *Manager) {
	m.Lock()
	defer m.Unlock()
	for i, mgr := range m.mrs {
		if mgr == manager {
			m.mrs[i] = m.mrs[len(m.mrs)-1]
			m.mrs[len(m.mrs)-1] = nil
			m.mrs = m.mrs[:len(m.mrs)-1]
		}
	}
}

// Manager manage binding and piping
type Manager struct {
	Username       string
	Port           int
	Backend        string
	NumConnCreated int64
	NumConnClosed  int64
	BytesUpload    int64
	BytesDownload  int64

	l net.Listener
}

func (m *Manager) String() string {
	return fmt.Sprintf("%s,%d,%s,%d,%d,%d,%d", m.Username, m.Port, m.Backend, m.NumConnCreated, m.NumConnClosed, m.BytesUpload, m.BytesDownload)
}

// Bind ...
func (m *Manager) Bind() error {
	return nil
}

// Start ...
func (m *Manager) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", m.Port))
	if err != nil {
		return err
	}
	fmt.Println("listening", l.Addr())
	m.l = l

	go func() {
		for {
			c, err := l.Accept()
			atomic.AddInt64(&m.NumConnCreated, 1)
			connConnectCounter.WithLabelValues(m.getMetricsTags()...).Inc()
			if err != nil {
				log.Println(err)
				if nerr, ok := err.(net.Error); ok {
					if nerr.Temporary() {
						continue
					}
					break
				}
			}

			bc, err := net.Dial("tcp", m.Backend)
			if err != nil {
				log.Println(err)
				m.closeConn(c)
				continue
			}

			log.Println("pipe", m.Port, m.Backend)
			m.pipeWithMetrics(bc, c)
		}
	}()
	return nil
}

// Stop ...
func (m *Manager) Stop() {
	m.l.Close()
}

// Use metrics and store the user usage
// upload/download
func (m *Manager) Use(n int, dir Direction) {
	if dir == DirectionUpload {
		atomic.AddInt64(&m.BytesUpload, int64(n))
		bytesUploadVec.WithLabelValues(m.getMetricsTags()...).Observe(float64(n))
	} else {
		atomic.AddInt64(&m.BytesDownload, int64(n))
		bytesDownloadVec.WithLabelValues(m.getMetricsTags()...).Observe(float64(n))
	}

	DefaultStorage.Write(&Record{
		Username:  m.Username,
		BytesUsed: n,
		Dir:       dir,
		Time:      time.Now().Unix(),
	})
}

// bc for backend conn, c for client conn
func (m *Manager) pipeWithMetrics(bc, c net.Conn) {
	// c -> bc
	b1 := make([]byte, 512)
	go func() {
		defer c.Close()
		defer bc.Close()
		for {
			n, err := c.Read(b1)
			if n > 0 {
				m.Use(n, DirectionUpload)
				bc.Write(b1[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	b2 := make([]byte, 512)
	go func() {
		defer m.closeConn(c)
		defer bc.Close()
		for {
			n, err := bc.Read(b2)
			if n > 0 {
				m.Use(n, DirectionDownload)
				c.Write(b2[:n])
			}
			if err != nil {
				break
			}
		}
	}()
}

func (m *Manager) closeConn(c net.Conn) {
	atomic.AddInt64(&m.NumConnClosed, int64(1))
	connClosedCounter.WithLabelValues(m.getMetricsTags()...).Inc()
	c.Close()
}

func (m *Manager) getMetricsTags() []string {
	return []string{
		m.Username,
		fmt.Sprintf("%d", m.Port),
		m.Backend,
	}
}
