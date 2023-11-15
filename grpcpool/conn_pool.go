package grpcpool

import (
	"log/slog"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnPool struct {
	pool   chan *grpc.ClientConn
	mu     sync.Mutex
	conns  map[*grpc.ClientConn]bool
	target string // 用於動態更改目標地址
}

// Resize 方法用於動態調整連接池的大小
func (p *ConnPool) Resize(newSize int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	currentSize := len(p.pool)
	// 縮小連接池
	for currentSize > newSize {
		conn := <-p.pool
		delete(p.conns, conn)
		conn.Close()
		currentSize--
	}
	// 擴大連接池
	for currentSize < newSize {
		conn, err := grpc.Dial(p.target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			slog.Error("NewConnPool", "Failed to connect: %v", err)
			return
		}
		p.pool <- conn
		p.conns[conn] = true
		currentSize++
	}
}

func (p *ConnPool) Get() *grpc.ClientConn {
	return <-p.pool
}

func (p *ConnPool) Put(conn *grpc.ClientConn) {
	// p.pool <- conn

	// 如果池已經滿了，關閉並釋放連接，否則放回池中
	if len(p.pool) == cap(p.pool) {
		delete(p.conns, conn)
		conn.Close()
	} else {
		p.pool <- conn
	}
}

func (p *ConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.pool)
	for conn := range p.pool {
		delete(p.conns, conn)
		conn.Close()
	}
}
