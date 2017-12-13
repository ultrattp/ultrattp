package internals

import "net"

// TCPConn is a helper type used in acceptors
type TCPConn struct {
	*net.TCPConn
}

func (conn *TCPConn) Setup() {
	// conn.SetLinger(10)
	// conn.SetNoDelay(true)
	// conn.SetReadBuffer(128)
}
