package server

func (ctx *RequestCtx) Write(b []byte) {
	ctx.buf.Write(b)
}

func (ctx *RequestCtx) SetConnectionClosed(v bool) {
	ctx.connClosed = v
}
