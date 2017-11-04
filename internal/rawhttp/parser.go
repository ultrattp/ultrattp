package rawhttp

import (
	"bytes"
	"strings"
)

// TestCase is a rawhttp case data
type TestCase struct {
	Name            string
	Body            []byte
	GuessIsResponse bool
}

// Parse rawhttp testdata format
func Parse(raw []byte) []TestCase {
	var res []TestCase

	caseOpened := false
	currentCase := TestCase{}
	captureName := false

	for i := 0; i < len(raw); i++ {
		switch {
		case bytes.Equal(raw[i:i+2], []byte(`\n`)):
			currentCase.Body = append(currentCase.Body, '\n')
			i++
		case bytes.Equal(raw[i:i+2], []byte(`\r`)):
			currentCase.Body = append(currentCase.Body, '\r')
			i++
		case bytes.Equal(raw[i:i+2], []byte(`\t`)):
			currentCase.Body = append(currentCase.Body, '\t')
			i++
		case raw[i] == '\n', raw[i] == '\t', raw[i] == '\r':
			// skip the real ones
		case captureName:
			captureName = false
			for ; raw[i] != '\n'; i++ {
				currentCase.Name += string(raw[i])
			}
			currentCase.Name = strings.TrimSpace(currentCase.Name)
		case bytes.Equal(raw[i:i+6], []byte("======")):
			if caseOpened {
				currentCase.GuessIsResponse = bytes.Equal(currentCase.Body[0:4], []byte("HTTP"))
				res = append(res, currentCase)
				caseOpened = false
				currentCase = TestCase{}
			} else {
				caseOpened = true
				captureName = true
			}
			i += 6
			continue
		default:
			currentCase.Body = append(currentCase.Body, raw[i])
		}
	}

	return res
}
