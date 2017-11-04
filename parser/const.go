package parser

const (
	// HTTPTypeUnknown is unknown http state. May occurs on errors
	HTTPTypeUnknown HTTPType = "unknown"
	// HTTPTypeRequest is a request type
	HTTPTypeRequest HTTPType = "request"
	// HTTPTypeResponse is a response type
	HTTPTypeResponse HTTPType = "response"
)
