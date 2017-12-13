package parser

import "sync"

// HTTPType is a type that used to store http body type
// which can be request or response. can be empty if unknown state occurred
type HTTPType uint

// Result of the parsing
type Result struct {
	body          []byte
	bodyLen       int
	rawHeadersEnd int
	preambulaEnd  int

	// result caches

	httpType         HTTPType
	hasBody          bool
	allHeadersParsed bool
	isBroken         bool
	firstLineParsed  bool
	isLoaded         bool

	// remove padding if needed
	padding [3]bool

	lastParseHeadersState int
	parsedHeaders         map[string][]byte
	reqMethod             []byte
	reqPath               []byte
	proto                 []byte

	statusCode   int
	statusString []byte

	mu *sync.RWMutex
}
