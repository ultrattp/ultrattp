package server

import (
	"errors"
	"net"
	"os"
	"runtime"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp/reuseport"
)

var LoggingLevel = log.DebugLevel

// var LoggingLevel = log.ErrorLevel

var logHandler = cli.New(os.Stdout)
var globalLog = &log.Logger{
	Level:   LoggingLevel,
	Handler: logHandler,
}

var ErrInvalidHandler = errors.New("invalid handler")
var ErrInvalidListener = errors.New("invalid listener")

var defaultServer = &Server{}

func Serve(ln *net.TCPListener, h func(*RequestCtx)) error {
	if ln == nil {
		return ErrInvalidListener
	}
	if h == nil {
		return ErrInvalidHandler
	}
	for i := 0; i < runtime.NumCPU()*2; i++ {
		aID := uuid.New().String()
		a := &Acceptor{
			id:  aID,
			s:   defaultServer,
			log: nil,
			// log: (&log.Logger{
			// 	Level:   LoggingLevel,
			// 	Handler: logHandler,
			// }).WithField("acceptorID", aID),
			handler: h,
		}
		go a.RunTCP(ln)
	}
	<-make(chan bool)
	return nil
}

func ListenAndServe(addr string, h func(*RequestCtx)) error {
	// ln, err := net.Listen("tcp4", addr)
	ln, err := reuseport.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	return Serve(ln.(*net.TCPListener), h)
}
