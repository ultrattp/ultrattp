package rawhttp

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name string
		args args
		want []TestCase
	}{
		{
			name: "basic",
			args: args{
				raw: []byte(`
======# super simple request
POST /test.ext?a=b+c HTTP/1.1\r\n
User-Agent: FooAgent\r\n
Host: trolohost\r\n\r\n
======

======# super simple request with weird http version
POST /test.ext?a=b+c HTTP/1.2\r\n
User-Agent: FooAgent\r\n
Host: trolohost\r\n\r\n
======

======# slightly weird response
HTTP/1.1 200 OK\r\n
Date: Mon, 23 May 2005 22:38:34 GMT\r\n
User-Agent: FooAgent\r\n
Content-Type: text/html; \n\t\t  charset=UTF-8\r\n
Host: trolohost\r\n\r\n
======
`),
			},
			want: []TestCase{
				{
					Name:            "super simple request",
					Body:            []byte("POST /test.ext?a=b+c HTTP/1.1\r\nUser-Agent: FooAgent\r\nHost: trolohost\r\n\r\n"),
					GuessIsResponse: false,
				},
				{
					Name:            "super simple request with weird http version",
					Body:            []byte("POST /test.ext?a=b+c HTTP/1.2\r\nUser-Agent: FooAgent\r\nHost: trolohost\r\n\r\n"),
					GuessIsResponse: false,
				},
				{
					Name:            "slightly weird response",
					Body:            []byte("HTTP/1.1 200 OK\r\nDate: Mon, 23 May 2005 22:38:34 GMT\r\nUser-Agent: FooAgent\r\nContent-Type: text/html; \n\t\t  charset=UTF-8\r\nHost: trolohost\r\n\r\n"),
					GuessIsResponse: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				for _, tcase := range got {
					t.Logf("tcase.Name=[%q] tcase.Body=[%s]", tcase.Name, string(tcase.Body))
				}
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
