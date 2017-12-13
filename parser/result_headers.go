package parser

import (
	"bytes"
	"errors"

	"github.com/gramework/runtimer"

	"github.com/ultrattp/ultrattp/util"
)

var (
	spaceByte = byte(' ')
	tabByte   = byte('\t')
	colonByte = byte(':')
)

// ParseHeader reads all headers until header with given name
// in case-insensetive way
func (r *Result) ParseHeader(name string) ([]byte, bool) {
	if len(name) == 0 || r.isBroken || r.Type() == HTTPTypeUnknown {
		return nil, false
	}

	name = util.NormalizeHeaderKeyStr(name)
	var (
		nameBytes = util.StringToBytes(name)
		key       []byte
		keyStr    string
		i         = r.preambulaEnd + 2
	)

	if r.parsedHeaders == nil {
		r.parsedHeaders = make(map[string][]byte, 8)
	} else {
		if v, ok := r.parsedHeaders[name]; ok {
			return v, ok
		}
	}

	if r.lastParseHeadersState != 0 {
		i = r.lastParseHeadersState
	}
	for idx := 0; i < r.rawHeadersEnd; i++ {
		idx = bytes.IndexByte(r.body[i:], colonByte)
		key = bytes.TrimSpace(r.body[i : i+idx])
		util.NormalizeHeaderKey(key)
		keyStr = runtimer.Gostringnocopy(&key[0])
		i += idx + 1
		// for ; i < len(r.body) && (r.body[i] == spaceByte || r.body[i] == tabByte); i++ {
		// }
		valueIdx := bytes.Index(r.body[i:], lineSplitter)
		valueEnd := i + valueIdx
		r.parsedHeaders[keyStr] = bytes.TrimSpace(r.body[i:valueEnd])

		r.lastParseHeadersState = valueEnd

		if runtimer.BytesHash(key, 8) == runtimer.BytesHash(nameBytes, 8) {
			return r.parsedHeaders[keyStr], true
		}
		i = valueEnd
	}

	return nil, false
}

func (r *Result) Proto() []byte {
	if r.isBroken == true {
		return nil
	}
	return r.proto
}

func (r *Result) readFirstLine() {
	if r.firstLineParsed {
		return
	}
	if r.bodyLen == 0 {
		r.isBroken = true
		return
	}
	var i int
	if r.Type() == HTTPTypeRequest {
		for ; i < r.bodyLen && (r.body[i] != spaceByte && r.body[i] != tabByte); i++ {
		}
		r.reqMethod = r.body[:i]
		for ; i < r.bodyLen && (r.body[i] == spaceByte || r.body[i] == tabByte); i++ {
			// skip spaces
		}
		var reqPathStartIdx = i
		for ; i < r.bodyLen && (r.body[i] != spaceByte && r.body[i] != tabByte); i++ {
		}
		r.reqPath = r.body[reqPathStartIdx:i]
		for ; i < r.bodyLen && (r.body[i] == spaceByte || r.body[i] == tabByte); i++ {
			// skip spaces
		}
		protoEnding := bytes.Index(r.body[i:], lineSplitter)
		if protoEnding != -1 {
			r.proto = r.body[i : i+protoEnding]
			r.preambulaEnd += len(r.proto)
		} else {
			r.isBroken = true
		}
		// parsed line
	} else {
		for ; i < r.bodyLen && (r.body[i] != spaceByte && r.body[i] != tabByte); i++ {
		}
		r.proto = r.body[:i]
		for ; i < r.bodyLen && (r.body[i] == spaceByte || r.body[i] == tabByte); i++ {
			// skip spaces
		}

		var statusCodeStartIdx = i
		for ; i < r.bodyLen && (r.body[i] != spaceByte && r.body[i] != tabByte); i++ {
		}
		// statusCode, err := strconv.ParseUint(string(r.body[statusCodeStartIdx:i]), 10, 64)
		statusCode, _, err := parseUintBuf(r.body[statusCodeStartIdx:i])
		if err != nil {
			r.isBroken = true
			return
		}
		r.statusCode = statusCode

		r.statusString = r.body[i : i+bytes.Index(r.body[i:], lineSplitter)]
		r.preambulaEnd += len(r.statusString)
	}
	r.preambulaEnd += i
	r.firstLineParsed = true
}

var (
	errEmptyInt               = errors.New("empty integer")
	errUnexpectedFirstChar    = errors.New("unexpected first char found. Expecting 0-9")
	errUnexpectedTrailingChar = errors.New("unexpected traling char found. Expecting 0-9")
	errTooLongInt             = errors.New("too long int")
)

func parseUintBuf(b []byte) (int, int, error) {
	n := len(b)
	if n == 0 {
		return -1, 0, errEmptyInt
	}
	v := 0
	for i := 0; i < n; i++ {
		c := b[i]
		k := c - '0'
		if k > 9 {
			if i == 0 {
				return -1, i, errUnexpectedFirstChar
			}
			return v, i, nil
		}
		if i >= maxIntChars {
			return -1, i, errTooLongInt
		}
		v = 10*v + int(k)
	}
	return v, n, nil
}
