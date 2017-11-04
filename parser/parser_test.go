package parser

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/ultrattp/ultrattp/internal/rawhttp"
)

var basicTestFile, btfReadErr = ioutil.ReadFile("../testdata/basic.rawhttp")

func TestBasicParser(t *testing.T) {
	if btfReadErr != nil {
		t.Fatalf("Can't read cases: %s", btfReadErr)
	}

	cases := rawhttp.Parse(basicTestFile)

	for id, _case := range cases {
		result := Parse(_case.Body)

		if _case.GuessIsResponse && result.Type() != HTTPTypeResponse {
			t.Fatalf(
				"invalid type: got %q but response guessed, test id #%v",
				result.Type(),
				id,
			)
		}
	}
}

func TestHeaders(t *testing.T) {
	headerTestFile, err := ioutil.ReadFile("../testdata/headers_1.rawhttp")
	if err != nil {
		t.Fatalf("Can't read cases: %s", err)
	}

	cases := rawhttp.Parse(headerTestFile)

	for id, _case := range cases {
		parsed := Parse(_case.Body)

		value, wasFound := parsed.ParseHeader("Name")
		if !wasFound {
			t.Fatalf(
				"Name header required but was not found in case #%d",
				id,
			)
		}
		if !bytes.Equal(value, []byte(_case.Name)) {
			t.Fatalf(
				"Name header was found but got %q instead of expected %q in case %d",
				string(value),
				_case.Name,
				id,
			)
		}
		value, wasFound = parsed.ParseHeader("Host")
		if !_case.GuessIsResponse {
			if !wasFound {
				t.Fatalf(
					"Host header required but was not found in case #%d",
					id,
				)
			}
			if !bytes.Equal(value, []byte("trolohost")) {
				t.Fatalf(
					"Host header was found but got %q instead of expected %q in case %d",
					string(value),
					_case.Name,
					id,
				)
			}
		}
	}

}

var benchCase = rawhttp.Parse(basicTestFile)[0]

func BenchmarkBasicParser(b *testing.B) {
	if btfReadErr != nil {
		b.Fatalf("Can't read cases: %s", btfReadErr)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := Parse(benchCase.Body)

		if benchCase.GuessIsResponse && result.Type() != HTTPTypeResponse {
			b.Fatalf(
				"invalid type: got %q but response guessed, test id #0",
				result.Type(),
			)
		}
	}
}
