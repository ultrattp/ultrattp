package parser

import "sync"

// Parse the given body
func Parse(body []byte) *Result {
	res := &Result{
		body: body,
		mu:   &sync.RWMutex{},
	}

	res.loadHeaders()
	res.loadBody()

	return res
}
