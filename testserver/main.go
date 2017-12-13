package main

import (
	"fmt"
	"time"

	"github.com/ultrattp/ultrattp"
)

func main() {
	cnt := "hi"
	// hi := []byte("HTTP/1.1 200 OK\r\nDate: Wed, 13 Dec 2017 11:41:04 GMT\r\nContent-Length: 2\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nhi")
	ultrattp.ListenAndServe(":1025", func(ctx *ultrattp.RequestCtx) {
		var now, lastMod []byte
		now = time.Now().In(time.UTC).AppendFormat(now, time.RFC1123)
		lastMod = time.Unix(0, 0).In(time.UTC).AppendFormat(lastMod, time.RFC1123)
		copy(now[len(now)-3:], []byte("GMT"))
		copy(lastMod[len(now)-3:], []byte("GMT"))
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
	})
}
