package parser

import (
	"bytes"
	"strconv"
	"strings"
)

var lineSplitter = []byte("\r\n")

// ParseHeader reads all headers until header with given name
// in case-insensetive way
func (r *Result) ParseHeader(name string) ([]byte, bool) {
	if r.isBroken || strings.Compare(string(r.Type()), string(HTTPTypeUnknown)) == 0 {
		return nil, false
	}
	r.readFirstLine()
	if r.isBroken {
		return nil, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.parsedHeaders == nil {
		r.parsedHeaders = make(map[string][]byte, 0)
	}

	var (
		key string
	)
	i := r.preambulaEnd + 2
	if r.lastParseHeadersState != 0 {
		i = r.lastParseHeadersState
	}
	for ; i < r.rawHeadersEnd; i++ {
		idx := bytes.IndexByte(r.body[i:], ':')
		key = strings.TrimSpace(string(r.body[i : i+idx]))
		i += idx + 1
		for ; i < len(r.body) && (r.body[i] == ' ' || r.body[i] == '\t'); i++ {
		}
		valueIdx := bytes.Index(r.body[i:], lineSplitter)
		r.parsedHeaders[key] = r.body[i : i+valueIdx]
		i += valueIdx

		r.lastParseHeadersState = i
		if key == name {
			return r.parsedHeaders[key], true
		}
	}

	return nil, false
}

func (r *Result) readFirstLine() {
	var i int
	if strings.Compare(string(r.Type()), string(HTTPTypeRequest)) == 0 {
		for ; i < len(r.body) && (r.body[i] != ' ' && r.body[i] != '\t'); i++ {
		}
		r.reqMethod = r.body[:i]
		for ; i < len(r.body) && (r.body[i] == ' ' || r.body[i] == '\t'); i++ {
			// skip spaces
		}
		var reqPathStartIdx = i
		for ; i < len(r.body) && (r.body[i] != ' ' && r.body[i] != '\t'); i++ {
		}
		r.reqPath = r.body[reqPathStartIdx:i]
		for ; i < len(r.body) && (r.body[i] == ' ' || r.body[i] == '\t'); i++ {
			// skip spaces
		}

		r.proto = r.body[i : i+bytes.Index(r.body[i:], lineSplitter)]
		r.preambulaEnd += len(r.proto)
		// parsed line
	} else {
		for ; i < len(r.body) && (r.body[i] != ' ' && r.body[i] != '\t'); i++ {
		}
		r.proto = r.body[:i]
		for ; i < len(r.body) && (r.body[i] == ' ' || r.body[i] == '\t'); i++ {
			// skip spaces
		}

		var statusCodeStartIdx = i
		for ; i < len(r.body) && (r.body[i] != ' ' && r.body[i] != '\t'); i++ {
		}
		statusCode, err := strconv.ParseUint(string(r.body[statusCodeStartIdx:i]), 10, 64)
		if err != nil {
			r.isBroken = true
			return
		}
		r.statusCode = int(statusCode)

		r.statusString = r.body[i : i+bytes.Index(r.body[i:], lineSplitter)]
		r.preambulaEnd += len(r.statusString)
	}
	r.preambulaEnd += i
}
