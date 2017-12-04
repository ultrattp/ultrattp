package ultrattp

import (
	"net"

	"github.com/ultrattp/ultrattp/server"
)

// Server instance
type Server = server.Server

type RequestCtx = server.RequestCtx

func Serve(ln *net.TCPListener) error {
	return server.Serve(ln)
}

func ListenAndServe(addr string, h func(ctx *RequestCtx)) error {
	return server.ListenAndServe(addr, h)
}
