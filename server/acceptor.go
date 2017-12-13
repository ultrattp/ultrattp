package server

import (
	"io"
	"net"
	"strings"
	"time"

	"github.com/gramework/runtimer"
	"github.com/kirillDanshin/bufioReaderPool"
	"github.com/ultrattp/ultrattp/parser"
	"github.com/ultrattp/ultrattp/server/internals"
	"github.com/ultrattp/ultrattp/util"
	"github.com/valyala/bytebufferpool"
)

const (
	http11 = "HTTP/1.1"
)

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

	max := initialMax
	var currentPos int64
	b.Set(<-acceptorBufPool)
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
		conn.Close()
		return
	}
	conn.CloseRead()

	// a.log.Fatalf("req=\n%s\n", b.String())

	pipelineReq, _ := parser.ParsePipeline(b.B)

	bytebufferpool.Put(b)

	respBuf := bytebufferpool.Get()

	var (
		rawReq          *parser.Result
		i               int
		shouldCloseConn bool
	)
	for i, rawReq = range pipelineReq {
		if a.log != nil {
			a.log.Debugf("processing %d", i)
		}
		if rawReq.Type() != parser.HTTPTypeRequest || rawReq.IsBroken() {
			continue
		}
		// ctx := &RequestCtx{
		// 	buf: respBuf,
		// }
		ctx := (*RequestCtx)(runtimer.GetEfaceDataPtr(reqCtxPool.Get()))
		ctx.buf = respBuf

		a.handler(ctx)
		if ctx.connClosed {
			shouldCloseConn = ctx.connClosed
		}
	}

	n, err := conn.Write(b.B)
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
	if (i == 0 && shouldCloseConn) ||
		strings.Compare(strings.ToUpper(util.BytesToString(rawReq.Proto())), http11) != 0 {
		// if util.BytesToString(rawReq.Proto()) != http11 {
		// respBuf.WriteByte(byte(0x0))
		err = conn.Close()
		if err != nil {
			if a.log != nil {
				a.log.Errorf("error while closing connection n=[%d]: %s", n, err)
			}
		}
	}
}
