package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ibigbug/ss-account/database"
	"github.com/ibigbug/ss-account/metrics"
	"github.com/ibigbug/ss-account/utils"
)

type Direction int

const (
	DirectionUpload Direction = iota
	DirectionDownload
)

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

func (m *Managed) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.mrs)
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
	Username       string `json:"username,omitempty"`
	Port           string `json:"port,omitempty"`
	Backend        string `json:"backend,omitempty"`
	NumConnCreated int64  `json:"num_conn_created,omitempty"`
	NumConnClosed  int64  `json:"num_conn_closed,omitempty"`
	BytesUpload    int64  `json:"bytes_upload"`
	BytesDownload  int64  `json:"bytes_download"`
	Active         bool   `json:"active"`

	l net.Listener
}

func (m *Manager) String() string {
	return fmt.Sprintf("%s,%s,%s,%d,%d,%d,%d", m.Username, m.Port, m.Backend, m.NumConnCreated, m.NumConnClosed, m.BytesUpload, m.BytesDownload)
}

// Bind creates a port binding and username in database
func (m *Manager) Bind() error {
	DefaultManaged.Lock()
	defer DefaultManaged.Unlock()
	for _, u := range DefaultManaged.mrs {
		if u.Port == m.Port || u.Username == m.Username {
			return errors.New("Username or Port already existed")
		}
	}
	logger.Printf("bind user %s -> %s <- %s\n", m.Username, m.Port, m.Backend)
	return DefaultStorage.BindPort(&database.Binding{
		Username: m.Username,
		Port:     m.Port,
		Backend:  m.Backend,
		Active:   true,
	})
}

// Unbind remove the backend and port from database
func (m *Manager) Unbind() error {
	return DefaultStorage.BindPort(&database.Binding{
		Username: m.Username,
		Backend:  "",
		Port:     "",
		Active:   false,
	})
}

// Start ...
func (m *Manager) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", m.Port))
	if err != nil {
		return err
	}

	if err := m.Bind(); err != nil {
		return err
	}
	m.l = l

	go func() {
		for {
			c, err := l.Accept()
			atomic.AddInt64(&m.NumConnCreated, 1)
			metrics.ConnConnectCounter.WithLabelValues(m.getMetricsTags()...).Inc()
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

			m.pipeWithMetrics(bc, c)
		}
	}()
	return nil
}

// Stop ...
func (m *Manager) Stop() {
	m.l.Close()
	m.Unbind()
}

// Use metrics and store the user usage
// upload/download
func (m *Manager) Use(n int, dir Direction) {
	if dir == DirectionUpload {
		atomic.AddInt64(&m.BytesUpload, int64(n))
		metrics.BytesUploadVec.WithLabelValues(m.getMetricsTags()...).Add(float64(n))
	} else {
		atomic.AddInt64(&m.BytesDownload, int64(n))
		metrics.BytesDownloadVec.WithLabelValues(m.getMetricsTags()...).Add(float64(n))
	}

	DefaultStorage.Write(&database.Record{
		Username:  m.Username,
		BytesUsed: n,
		Dir:       int(dir),
		Time:      time.Now().Unix(),
	})
}

// bc for backend conn, c for client conn
func (m *Manager) pipeWithMetrics(bc, c net.Conn) {
	// c -> bc
	var b1 [1024]byte
	select {
	case b1 = <-utils.FreeList:
	default:
		b1 = [1024]byte{}
	}
	go func() {
		defer c.Close()
		defer bc.Close()
		defer func() {
			utils.FreeList <- b1
		}()
		for {
			n, err := c.Read(b1[:])
			if n > 0 {
				m.Use(n, DirectionUpload)
				bc.Write(b1[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	var b2 [1024]byte
	select {
	case b2 = <-utils.FreeList:
	default:
		b2 = [1024]byte{}
	}
	go func() {
		defer m.closeConn(c)
		defer bc.Close()
		defer func() {
			utils.FreeList <- b2
		}()
		for {
			n, err := bc.Read(b2[:])
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
	metrics.ConnClosedCounter.WithLabelValues(m.getMetricsTags()...).Inc()
	c.Close()
}

func (m *Manager) getMetricsTags() []string {
	return []string{
		m.Username,
		m.Port,
		m.Backend,
	}
}
