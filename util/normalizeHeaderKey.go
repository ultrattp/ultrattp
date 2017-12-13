// The MIT License (MIT)
//
// Copyright (c) 2015-2016 Aliaksandr Valialkin, VertaMedia
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package util

const ToLower = 'a' - 'A'

var ToLowerTable = func() [256]byte {
	var a [256]byte
	for i := 0; i < 256; i++ {
		c := byte(i)
		if c >= 'A' && c <= 'Z' {
			c += ToLower
		}
		a[i] = c
	}
	return a
}()

var ToUpperTable = func() []byte {
	a := make([]byte, 256)
	for i := 0; i < 256; i++ {
		c := byte(i)
		if c >= 'a' && c <= 'z' {
			c -= ToLower
		}
		a[i] = c
	}
	return a
}()

func LowercaseBytes(b []byte) {
	for i := 0; i < len(b); i++ {
		p := &b[i]
		*p = ToLowerTable[*p]
	}
}

func NormalizeHeaderKeyStr(s string) string {
	n := len(s)
	if n == 0 {
		return s
	}

	res := []byte(s)
	res[0] = ToUpperTable[res[0]]
	for i := 1; i < n; i++ {
		p := &res[i]
		if *p == '-' {
			i++
			if i < n {
				res[i] = ToUpperTable[res[i]]
			}
			continue
		}
		*p = ToLowerTable[*p]
	}

	return BytesToString(res)
}

func NormalizeHeaderKey(b []byte) {
	n := len(b)
	if n == 0 {
		return
	}

	b[0] = ToUpperTable[b[0]]
	for i := 1; i < n; i++ {
		p := &b[i]
		if *p == '-' {
			i++
			if i < n {
				b[i] = ToUpperTable[b[i]]
			}
			continue
		}
		*p = ToLowerTable[*p]
	}
}
