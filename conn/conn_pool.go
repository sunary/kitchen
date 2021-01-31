package conn

import (
	"net"
	"sync"
	"time"
)

const (
	connectIdleTime = 30
)

// Object ...
type Object struct {
	conn *net.TCPConn
	idle int64
}

// ConnectPool ...
type ConnectPool struct {
	sync.RWMutex
	pools   map[string]*Pool
	mincap  int
	maxcap  int
	timeout int64
}

// NewConnectPool ...
func NewConnectPool() *ConnectPool {
	cp := &ConnectPool{
		pools:   make(map[string]*Pool),
		mincap:  5,
		maxcap:  80,
		timeout: int64(time.Second * connectIdleTime),
	}
	go cp.autoRelease()

	return cp
}

// DailTimeOut ...
func DailTimeOut(target string, timeout time.Duration) (c *net.TCPConn, err error) {
	var connect net.Conn
	connect, err = net.DialTimeout("tcp", target, timeout)
	if err == nil {
		conn := connect.(*net.TCPConn)
		conn.SetKeepAlive(true)
		conn.SetNoDelay(true)
		c = conn
	}
	return
}

// GetConnect ...
func (cp *ConnectPool) GetConnect(targetAddr string) (c *net.TCPConn, err error) {
	cp.RLock()
	pool, ok := cp.pools[targetAddr]
	cp.RUnlock()
	if !ok {
		cp.Lock()
		pool = NewPool(cp.mincap, cp.maxcap, cp.timeout, targetAddr)
		cp.pools[targetAddr] = pool
		cp.Unlock()
	}

	return pool.GetConnectFromPool()
}

// PutConnect ...
func (cp *ConnectPool) PutConnect(c *net.TCPConn, forceClose bool) {
	if c == nil {
		return
	}
	if forceClose {
		c.Close()
		return
	}
	addr := c.RemoteAddr().String()
	cp.RLock()
	pool, ok := cp.pools[addr]
	cp.RUnlock()
	if !ok {
		c.Close()
		return
	}
	object := &Object{conn: c, idle: time.Now().UnixNano()}
	pool.PutConnectObjectToPool(object)
}

func (cp *ConnectPool) autoRelease() {
	for {
		pools := make([]*Pool, 0)
		cp.RLock()
		for _, pool := range cp.pools {
			pools = append(pools, pool)
		}
		cp.RUnlock()
		for _, pool := range pools {
			pool.autoRelease()
		}
		time.Sleep(time.Second)
	}

}
