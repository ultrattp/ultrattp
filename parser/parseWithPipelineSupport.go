package parser

import (
	"bytes"
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

var LoggingLevel = log.DebugLevel

// var LoggingLevel = log.ErrorLevel

var logHandler = cli.New(os.Stdout)
var globalLog = &log.Logger{
	Level:   LoggingLevel,
	Handler: logHandler,
}

// ErrBrokenPipeline notifies about broken HTTP pipeline request
var ErrBrokenPipeline = errors.New("can't parse broken pipeline")

var partSplitter = []byte("\r\n\r\n")

// ParsePipeline is a shortcut for parsing an HTTP pipeline.
// When pipeline is broken, it returns successfully parsed
// results, if any
func ParsePipeline(body []byte) ([]*Result, error) {
	var r *Result
	var res []*Result
	headerEnd := bytes.Index(body, partSplitter)
	if headerEnd == -1 {
		return res, ErrBrokenPipeline
	}
	bodyEnd := bytes.Index(body[headerEnd:], partSplitter)
	if bodyEnd == -1 { // only one request
		return []*Result{
			Parse(body),
		}, nil
	}
	position := 0
	currPosition := 0
	bodyLen := len(body)
	for position < bodyLen {
		currPosition = position
		headerEnd = bytes.Index(body[position:], partSplitter)
		if headerEnd == -1 {
			return res, nil
		}
		currPosition += headerEnd + 1
		if currPosition < bodyLen {
			bodyEnd = bytes.Index(body[currPosition:], partSplitter)
			if bodyEnd == -1 { // only one request
				return []*Result{
					Parse(body),
				}, nil
			}
			currPosition += bodyEnd
		}
		globalLog.Debugf("len=%d position=%d headerEnd=%d bodyEnd=%d", bodyLen, position, headerEnd, bodyEnd)
		r = Parse(body[position:currPosition])
		if r.isBroken {
			globalLog.Errorf("broken pipeline?")
			return res, ErrBrokenPipeline
		}
		res = append(res, r)
		position = currPosition + 1
	}

	return res, nil
}
