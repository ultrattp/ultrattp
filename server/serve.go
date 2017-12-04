package server

import (
	"net"
	"os"
	"runtime"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/google/uuid"
)

var LoggingLevel = log.DebugLevel

var logHandler = cli.New(os.Stdout)
var globalLog = &log.Logger{
	Level:   LoggingLevel,
	Handler: logHandler,
}

var respNoContent = []byte(
	"HTTP/1.1 200 OK\r\n" +
		"Server: ultrattp/1.0\r\n" +
		"\r\n" +
		"Hi\r\n" +
		"\r\n",
)

var defaultServer = &Server{}

func Serve(ln *net.TCPListener) error {
	for i := 0; i < runtime.NumCPU(); i++ {
		aID := uuid.New().String()
		a := &Acceptor{
			id: aID,
			s:  defaultServer,
			log: (&log.Logger{
				Level:   LoggingLevel,
				Handler: logHandler,
			}).WithField("acceptorID", aID),
		}
		go a.RunTCP(ln)
	}
	<-make(chan bool)
	return nil
}

func ListenAndServe(addr string, h func(*RequestCtx)) error {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	return Serve(ln.(*net.TCPListener))
}
