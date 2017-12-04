package main

import "github.com/ultrattp/ultrattp"

func main() {
	ultrattp.ListenAndServe(":1025", func(*ultrattp.RequestCtx) {})
}
