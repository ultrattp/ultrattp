package server

import (
	"bufio"
	"io"
	"net"

	"github.com/valyala/bytebufferpool"
)

func (a *Acceptor) RunTCP(ln *net.TCPListener) {
	for {
		a.log.Debug("accepting connection")
		conn, err := ln.AcceptTCP()
		if err != nil {
			a.log.Errorf("can't accept connection to %q: %s", conn.RemoteAddr(), err)
			continue
		}
		go a.process(conn)
	}
}

func (a *Acceptor) process(conn *net.TCPConn) {
	conn.SetLinger(10)
	conn.SetNoDelay(true)
	conn.SetReadBuffer(128)
	a.log.Debug("reading data from connection")
	b := bytebufferpool.Get()
	connReader := bufio.NewReader(conn)

	startPos := int64(b.Len())
	max := int64(64)
	currentPos := startPos
	b.Set(make([]byte, max))
	for {
		if currentPos == max {
			max *= 2
			bNew := make([]byte, max)
			copy(bNew, b.B)
			b.Set(bNew)
		}
		readCount, err := connReader.Read(b.B[currentPos:])
		currentPos += int64(readCount)
		if err != nil {
			// currentPos -= startPos
			b.Set(b.B[:currentPos])
			if err != io.EOF {
				a.log.Errorf("error while reading from connection: %s", err)
			}
			break
		}

		if b.Len() != 0 && connReader.Buffered() == 0 {
			b.Set(b.B[:currentPos])
			break
		}
	}

	bytebufferpool.Put(b)

	err := conn.CloseRead()
	if err != nil {
		a.log.Errorf("error while closing read pipe")
	}
	// time.Sleep(10 * time.Millisecond)
	n, err := conn.Write(getRespNoContent())
	a.log.Debugf("written %d bytes", n)
	if err != nil {
		a.log.Errorf("error while writing to connection n=[%d]: %s", n, err)
	}
	err = conn.Close()
	if err != nil {
		a.log.Errorf("error while closing connection n=[%d]: %s", n, err)
	}
}
