package parser

import "bytes"

// Type returns http type: request, response or unknown if state is invalid
func (r *Result) Type() HTTPType {
	if r.httpType == "" {
		if r.rawHeadersEnd < 4 {
			r.isBroken = true
			r.httpType = HTTPTypeUnknown
			return r.httpType
		}
		if bytes.HasPrefix(r.body, strHTTP) {
			r.httpType = HTTPTypeResponse
		} else {
			r.httpType = HTTPTypeRequest
		}
	}

	return r.httpType
}

func (r *Result) IsBroken() bool {
	return r.isBroken
}
