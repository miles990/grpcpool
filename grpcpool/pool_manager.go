package grpcpool

import (
	"sync"

	"google.golang.org/grpc"
)

const (
	DefaultConnPoolSize = 10
)

var (
	once     sync.Once
	instance *PoolManager
)

type PoolManager struct {
	mu            sync.Mutex
	pools         map[string]*ConnPool
	defaultTarget string // set default target
}

func GetManager() *PoolManager {
	once.Do(func() {
		instance = &PoolManager{
			pools: make(map[string]*ConnPool),
		}
	})
	return instance
}

// New grpc conn pool
func (m *PoolManager) NewConnPool(maxSize int, target string) {
	p := &ConnPool{
		pool:   make(chan *grpc.ClientConn, maxSize),
		conns:  make(map[*grpc.ClientConn]bool),
		target: target, // set grpc host
	}
	m.defaultTarget = target
	p.Resize(maxSize)
	m.pools[target] = p
}

// Get default connection pool, and create a new one if it does not exist
func (m *PoolManager) DefaultConnPool() *ConnPool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.defaultTarget == "" {
		panic("default grpc target is empty")
	}

	pool, exists := m.pools[m.defaultTarget]
	if !exists {
		m.NewConnPool(DefaultConnPoolSize, m.defaultTarget)
	}
	return pool
}

// Get the connection pool based on the specified target address, and create a new one if it does not exist
func (m *PoolManager) GetConnPool(target string) *ConnPool {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, exists := m.pools[target]
	if !exists {
		m.NewConnPool(DefaultConnPoolSize, target)
	}
	return pool
}
