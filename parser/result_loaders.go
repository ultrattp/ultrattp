package parser

import (
	"bytes"
)

var bodySplitter = []byte("\r\n\r\n")

func (r *Result) loadHeaders() {
	if idx := bytes.Index(r.body, bodySplitter); idx > 4 {
		r.rawHeadersEnd = idx
		return
	}

	r.httpType = HTTPTypeUnknown
}

func (r *Result) loadBody() {
	if len(r.body) > r.rawHeadersEnd+4 {
		r.hasBody = true
	}
}
