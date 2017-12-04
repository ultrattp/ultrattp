package server

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/google/uuid"
	"github.com/ultrattp/ultrattp/util"
)

var LoggingLevel = log.DebugLevel

var logHandler = cli.New(os.Stdout)
var globalLog = &log.Logger{
	Level:   LoggingLevel,
	Handler: logHandler,
}

var respNoContent = "HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/plain;charset=utf8\r\n" +
	"Content-Language: en-US\r\n" +
	"Server: ultrattp/1.0\r\n" +
	"Date: %s\r\n" +
	"Content-Length: 3\r\n" +
	"Connection: closed\r\n" +
	"Last-Modified: " + time.Unix(0, 0).In(time.UTC).Format(time.RFC1123) + "\r\n" +
	"\r\n" +
	"Hi\n" +
	"\r\n"

var defaultServer = &Server{}

func getRespNoContent() []byte {
	now := time.Now().In(time.UTC).Format(time.RFC1123)
	return util.StringToBytes(fmt.Sprintf(respNoContent, now))
}

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
