package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/ultrattp/ultrattp"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	cnt := "hi"

	// hi := []byte("HTTP/1.1 200 OK\r\nDate: Wed, 13 Dec 2017 11:41:04 GMT\r\nContent-Length: 2\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nhi")
	var lastMod []byte
	lastMod = time.Unix(0, 0).In(time.UTC).AppendFormat(lastMod, time.RFC1123)
	copy(lastMod[len(lastMod)-3:], []byte("GMT"))
	h := func(ctx *ultrattp.RequestCtx) {
		var now []byte
		now = time.Now().In(time.UTC).AppendFormat(now, time.RFC1123)
		copy(now[len(now)-3:], []byte("GMT"))
		hi := []byte("HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain;charset=utf8\r\n" +
			"Server: ultrattp/1.0\r\n" +
			"Date: " + string(now) + "\r\n" +
			"Content-Length: " + fmt.Sprintf("%d", len(cnt)) + "\r\n" +
			// "Connection: close\r\n" +
			"Last-Modified: " + string(lastMod) + "\r\n" +
			"\r\n" +
			// cnt)
			cnt +
			"\r\n")
		// ctx.SetConnectionClosed(true)
		ctx.Write(hi)
	}
	if *memprofile != "" || *cpuprofile != "" {
		go ultrattp.ListenAndServe(":1025", h)
		<-time.After(10 * time.Second)
	} else {
		ultrattp.ListenAndServe(":1025", h)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
