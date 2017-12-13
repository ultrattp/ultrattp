package server

import (
	"io"
	"net"
	"time"

	"github.com/gramework/runtimer"

	"github.com/kirillDanshin/bufioReaderPool"
	"github.com/ultrattp/ultrattp/parser"
	"github.com/ultrattp/ultrattp/server/internals"
	"github.com/valyala/bytebufferpool"
)

const (
	http11 = "HTTP/1.1"
)

var http11Hash = runtimer.BytesHash([]byte(http11), 8)

func (a *Acceptor) RunTCP(ln *net.TCPListener) {
	for {
		if a.log != nil {
			a.log.Debug("accepting connection")
		}
		conn, err := ln.AcceptTCP()
		if err != nil {
			if a.log != nil {
				a.log.Errorf("can't accept connection to %q: %s", conn.RemoteAddr(), err)
			}
			continue
		}
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(10 * time.Minute)
		go a.process(conn)
	}
}

func (a *Acceptor) process(rawConn *net.TCPConn) {
	conn := internals.TCPConnPool.Get(rawConn)
	connReader := bufioReaderPool.Get(conn)
	b := bytebufferpool.Get()

	defer internals.TCPConnPool.Put(conn)
	defer bufioReaderPool.Put(connReader)
	conn.Setup()
	var (
		i               int
		shouldCloseConn bool
	)
	for {
		max := initialMax
		var currentPos int64
		b.Set(make([]byte, 64))
		for {
			if currentPos == max {
				max *= 2
				bNew := make([]byte, max)
				copy(bNew, b.B)
				b.Set(bNew)
			}
			readCount, err := connReader.Read(b.B[currentPos:])
			currentPos += int64(readCount)
			if err != nil || (b.Len() != 0 && connReader.Buffered() == 0) {
				b.Set(b.B[:currentPos])
				if err != nil && err != io.EOF {
					if a.log != nil {
						a.log.Errorf("error while reading from connection: %s", err)
					}
				}
				break
			}
		}

		if b.Len() == 0 {
			if a.log != nil {
				a.log.Errorf("read 0 bytes")
			}
			conn.Close()
			return
		}
		// conn.CloseRead()

		// pipelineReq, _ := parser.ParsePipeline(b.B)
		rawReq := parser.Parse(b.B)
		protoHash := runtimer.BytesHash(rawReq.Proto(), 8)

		// bytebufferpool.Put(b)

		// respBuf := bytebufferpool.Get()
		// for i, rawReq = range pipelineReq {
		if rawReq.Type() != parser.HTTPTypeRequest || rawReq.IsBroken() {
			a.log.Debugf("not a request or broken", i)
			// continue
			return
		}
		// ctx := &RequestCtx{
		// 	buf: respBuf,
		// }
		b.Reset()
		ctx := &RequestCtx{
			buf: b,
		}

		a.handler(ctx)
		if i != 0 {
			shouldCloseConn = ctx.connClosed
		}
		i++
		// }

		conn.SetNoDelay(true)
		n, err := conn.Write(b.B)
		b.Reset()
		// conn.CloseWrite()

		// time.Sleep(10 * time.Millisecond)
		// n, err := conn.Write(getRespNoContent())
		if a.log != nil {
			a.log.Debugf("written %d bytes i=%d", n, i)
		}
		if err != nil {
			if a.log != nil {
				a.log.Errorf("error while writing to connection n=[%d]: %s", n, err)
			}
		}
		if (i == 0 && shouldCloseConn) || protoHash != http11Hash {
			if a.log != nil {
				a.log.Errorf("closing conn! i %v shouldClose %v proto %q", i, shouldCloseConn, protoHash)
			}
			err = conn.Close()
			if err != nil {
				if a.log != nil {
					a.log.Errorf("error while closing connection n=[%d]: %s", n, err)
				}
			}
			return
		}
	}
}
