package server

import "github.com/apex/log"

// Server instance
type Server struct {
}

// RequestCtx is the context of the request
type RequestCtx struct {
}

// Acceptor instance. This struct stands here for better API
// usage experience and to simplify the code.
type Acceptor struct {
	s   *Server
	id  string
	log log.Interface
}
