package conn

import (
	"net"
	"time"
)

// Pool ...
type Pool struct {
	objects chan *Object
	mincap  int
	maxcap  int
	target  string
	timeout int64
}

// NewPool ...
func NewPool(min, max int, timeout int64, target string) (p *Pool) {
	p = new(Pool)
	p.mincap = min
	p.maxcap = max
	p.target = target
	p.objects = make(chan *Object, max)
	p.timeout = timeout
	p.initAllConnect()
	return p
}

func (p *Pool) initAllConnect() {
	for i := 0; i < p.mincap; i++ {
		c, err := net.Dial("tcp", p.target)
		if err == nil {
			conn := c.(*net.TCPConn)
			conn.SetKeepAlive(true)
			conn.SetNoDelay(true)
			o := &Object{conn: conn, idle: time.Now().UnixNano()}
			p.PutConnectObjectToPool(o)
		}
	}
}

// PutConnectObjectToPool ...
func (p *Pool) PutConnectObjectToPool(o *Object) {
	select {
	case p.objects <- o:
		return
	default:
		if o.conn != nil {
			o.conn.Close()
		}
		return
	}
}

func (p *Pool) autoRelease() {
	connectLen := len(p.objects)
	for i := 0; i < connectLen; i++ {
		select {
		case o := <-p.objects:
			if time.Now().UnixNano()-int64(o.idle) > p.timeout {
				o.conn.Close()
			} else {
				p.PutConnectObjectToPool(o)
			}
		default:
			return
		}
	}
}

// NewConnect ...
func (p *Pool) NewConnect(target string) (c *net.TCPConn, err error) {
	var connect net.Conn
	connect, err = net.Dial("tcp", p.target)
	if err == nil {
		conn := connect.(*net.TCPConn)
		conn.SetKeepAlive(true)
		conn.SetNoDelay(true)
		c = conn
	}
	return
}

// GetConnectFromPool ...
func (p *Pool) GetConnectFromPool() (c *net.TCPConn, err error) {
	var o *Object
	for i := 0; i < len(p.objects); i++ {
		select {
		case o = <-p.objects:
			if time.Now().UnixNano()-int64(o.idle) > p.timeout {
				o.conn.Close()
				o = nil
				break
			}
			return o.conn, nil
		default:
			return p.NewConnect(p.target)
		}
	}

	return p.NewConnect(p.target)
}
