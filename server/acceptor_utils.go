package server

import (
	"runtime"
	"sync"

	"github.com/valyala/bytebufferpool"
)

var initialMax int64 = 64

var acceptorBufPool = make(chan []byte, runtime.NumCPU()*2)

func acceptorBufAllocator() {
	for {
		acceptorBufPool <- make([]byte, 64)
	}
}

func init() {
	go acceptorBufAllocator()
}

func (a *Acceptor) readBody(b *bytebufferpool.ByteBuffer) {

}

var reqCtxPool = sync.Pool{
	New: func() interface{} {
		return &RequestCtx{}
	},
}
