package internals

import (
	"net"
	"sync"

	"github.com/gramework/runtimer"
)

type tcpConnPool struct {
	pool sync.Pool
}

func (p *tcpConnPool) Get(rawConn *net.TCPConn) (conn *TCPConn) {
	conn = (*TCPConn)(runtimer.GetEfaceDataPtr(p.pool.Get()))
	conn.TCPConn = rawConn
	return
}

func (p *tcpConnPool) Put(conn *TCPConn) {
	if conn != nil {
		p.pool.Put(conn)
	}
}

var TCPConnPool = &tcpConnPool{
	sync.Pool{
		New: func() interface{} {
			return &TCPConn{}
		},
	},
}
