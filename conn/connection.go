package conn

import (
	"net"
	"time"
)

type ConnTimeout struct {
	addr      string
	conn      net.Conn
	readTime  time.Duration
	writeTime time.Duration
}

func DialTimeout(addr string, connTime time.Duration) (*ConnTimeout, error) {
	conn, err := net.DialTimeout("tcp", addr, connTime)
	if err != nil {
		return nil, err
	}

	conn.(*net.TCPConn).SetNoDelay(true)
	conn.(*net.TCPConn).SetLinger(0)
	conn.(*net.TCPConn).SetKeepAlive(true)
	return &ConnTimeout{conn: conn, addr: addr}, nil
}

func NewConnTimeout(conn net.Conn) *ConnTimeout {
	if conn == nil {
		return nil
	}

	conn.(*net.TCPConn).SetNoDelay(true)
	conn.(*net.TCPConn).SetLinger(0)
	conn.(*net.TCPConn).SetKeepAlive(true)
	return &ConnTimeout{conn: conn, addr: conn.RemoteAddr().String()}
}

func (c *ConnTimeout) SetReadTimeout(timeout time.Duration) {
	c.readTime = timeout
}

func (c *ConnTimeout) SetWriteTimeout(timeout time.Duration) {
	c.writeTime = timeout
}

func (c *ConnTimeout) Read(p []byte) (int, error) {
	if c.readTime.Nanoseconds() > 0 {
		err := c.conn.SetReadDeadline(time.Now().Add(c.readTime))
		if err != nil {
			return 0, err
		}
	}

	return c.conn.Read(p)
}

func (c *ConnTimeout) Write(p []byte) (int, error) {
	if c.writeTime.Nanoseconds() > 0 {
		err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTime))
		if err != nil {
			return 0, err
		}
	}

	return c.conn.Write(p)
}

func (c *ConnTimeout) RemoteAddr() string {
	return c.addr
}

func (c *ConnTimeout) Close() error {
	return c.conn.Close()
}
