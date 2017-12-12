package main

import (
	"time"

	"github.com/ultrattp/ultrattp"
)

func main() {
	cnt := "Hi"
	var now []byte
	now = time.Now().In(time.UTC).AppendFormat(now, time.RFC1123)
	copy(now[len(now)-3:], []byte("GMT"))
	hi := []byte("HTTP/1.1 200 OK\r\n" +
		"Cache-Control: private\r\n" +
		"Content-Type: text/plain;charset=utf8\r\n" +
		"Content-Language: en-US\r\n" +
		"Server: ultrattp/1.0\r\n" +
		"Date: " + string(now) + "\r\n" +
		// "Content-Length: " + fmt.Sprintf("%d", len(cnt)) + "\r\n" +
		"Connection: closed\r\n" +
		"Last-Modified: " + time.Unix(0, 0).In(time.UTC).Format(time.RFC1123) + "\r\n" +
		"\r\n" +
		// cnt)
		cnt +
		"\r\n")
	ultrattp.ListenAndServe(":1025", func(ctx *ultrattp.RequestCtx) {
		ctx.SetConnectionClosed(true)
		ctx.Write(hi)
	})
}
