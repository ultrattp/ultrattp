package parser

const (
	_ HTTPType = iota
	// HTTPTypeUnknown is unknown http state. May occurs on errors
	HTTPTypeUnknown
	// HTTPTypeRequest is a request type
	HTTPTypeRequest
	// HTTPTypeResponse is a response type
	HTTPTypeResponse
)

const emptyString = ""
