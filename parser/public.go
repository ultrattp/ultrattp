package parser

import (
	"sync"
)

func (r *Result) reset() {
	*r = Result{
		mu: &sync.RWMutex{},
	}
}

// Parse the given body
func Parse(body []byte) *Result {
	r := &Result{
		body:    body,
		bodyLen: len(body),
		mu:      &sync.RWMutex{},
	}

	r.loadHeaders()
	r.loadBody()

	r.readFirstLine()
	r.isLoaded = true

	return r
}
