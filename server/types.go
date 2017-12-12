package server

import (
	"github.com/apex/log"
	"github.com/valyala/bytebufferpool"
)

// Server instance
type Server struct {
}

// RequestCtx is the context of the request
type RequestCtx struct {
	buf *bytebufferpool.ByteBuffer

	connClosed bool
}

// Acceptor instance. This struct stands here for better API
// usage experience and to simplify the code.
type Acceptor struct {
	s       *Server
	id      string
	log     log.Interface
	handler func(*RequestCtx)
}

/*

type requestInterface interface {
	// func AcquireRequest() *Request
	AppendBody(p []byte)
	AppendBodyString(s string)
	Body() []byte
	BodyGunzip() ([]byte, error)
	BodyInflate() ([]byte, error)
	BodyWriteTo(w io.Writer) error
	BodyWriter() io.Writer
	ConnectionClose() bool
	ContinueReadBody(r *bufio.Reader, maxBodySize int) error
	CopyTo(dst *requestInterface)
	Host() []byte
	IsBodyStream() bool
	MayContinue() bool
	MultipartForm() (*multipart.Form, error)
	PostArgs() *Args
	Read(r *bufio.Reader) error
	ReadLimitBody(r *bufio.Reader, maxBodySize int) error
	ReleaseBody(size int)
	RemoveMultipartFormFiles()
	RequestURI() []byte
	Reset()
	ResetBody()
	SetBody(body []byte)
	SetBodyStream(bodyStream io.Reader, bodySize int)
	SetBodyStreamWriter(sw StreamWriter)
	SetBodyString(body string)
	SetConnectionClose()
	SetHost(host string)
	SetHostBytes(host []byte)
	SetRequestURI(requestURI string)
	SetRequestURIBytes(requestURI []byte)
	String() string
	SwapBody(body []byte) []byte
	URI() *URI
	Write(w *bufio.Writer) error
	WriteTo(w io.Writer) (int64, error)
}

type requestCtxInterface interface {
	ConnID() uint64
	ConnRequestNum() uint64
	ConnTime() time.Time
	Error(msg string, statusCode int)
	FormFile(key string) (*multipart.FileHeader, error)
	FormValue(key string) []byte
	Hijack(handler HijackHandler)
	Hijacked() bool
	Host() []byte
	ID() uint64
	IfModifiedSince(lastModified time.Time) bool
	Init(req *Request, remoteAddr net.Addr, logger Logger)
	Init2(conn net.Conn, logger Logger, reduceMemoryUsage bool)
	IsBodyStream() bool
	IsDelete() bool
	IsGet() bool
	IsHead() bool
	IsPost() bool
	IsPut() bool
	IsTLS() bool
	LastTimeoutErrorResponse() *Response
	LocalAddr() net.Addr
	LocalIP() net.IP
	Logger() Logger
	Method() []byte
	MultipartForm() (*multipart.Form, error)
	NotFound()
	NotModified()
	Path() []byte
	PostArgs() *Args
	PostBody() []byte
	QueryArgs() *Args
	Redirect(uri string, statusCode int)
	RedirectBytes(uri []byte, statusCode int)
	Referer() []byte
	RemoteAddr() net.Addr
	RemoteIP() net.IP
	RequestURI() []byte
	ResetBody()
	SendFile(path string)
	SendFileBytes(path []byte)
	SetBody(body []byte)
	SetBodyStream(bodyStream io.Reader, bodySize int)
	SetBodyStreamWriter(sw StreamWriter)
	SetBodyString(body string)
	SetConnectionClose()
	SetContentType(contentType string)
	SetContentTypeBytes(contentType []byte)
	SetStatusCode(statusCode int)
	SetUserValue(key string, value interface{})
	SetUserValueBytes(key []byte, value interface{})
	String() string
	Success(contentType string, body []byte)
	SuccessString(contentType, body string)
	TLSConnectionState() *tls.ConnectionState
	Time() time.Time
	TimeoutError(msg string)
	TimeoutErrorWithCode(msg string, statusCode int)
	TimeoutErrorWithResponse(resp *Response)
	URI() *uriInterface
	UserAgent() []byte
	UserValue(key string) interface{}
	UserValueBytes(key []byte) interface{}
	VisitUserValues(visitor func([]byte, interface{}))
	Write(p []byte) (int, error)
	WriteString(s string) (int, error)
}

type uriInterface interface {
	// func AcquireURI() *URI
	AppendBytes(dst []byte) []byte
	CopyTo(dst uriInterface)
	FullURI() []byte
	Hash() []byte
	Host() []byte
	LastPathSegment() []byte
	Parse(host, uri []byte)
	Path() []byte
	PathOriginal() []byte
	QueryArgs() *ArgsInterface
	QueryString() []byte
	RequestURI() []byte
	Reset()
	Scheme() []byte
	SetHash(hash string)
	SetHashBytes(hash []byte)
	SetHost(host string)
	SetHostBytes(host []byte)
	SetPath(path string)
	SetPathBytes(path []byte)
	SetQueryString(queryString string)
	SetQueryStringBytes(queryString []byte)
	SetScheme(scheme string)
	SetSchemeBytes(scheme []byte)
	String() string
	Update(newURI string)
	UpdateBytes(newURI []byte)
	WriteTo(w io.Writer) (int64, error)
}

type ArgsInterface interface {
	// func AcquireArgs() *Args
	Add(key, value string)
	AddBytesK(key []byte, value string)
	AddBytesKV(key, value []byte)
	AddBytesV(key string, value []byte)
	AppendBytes(dst []byte) []byte
	CopyTo(dst *ArgsInterface)
	Del(key string)
	DelBytes(key []byte)
	GetBool(key string) bool
	GetUfloat(key string) (float64, error)
	GetUfloatOrZero(key string) float64
	GetUint(key string) (int, error)
	GetUintOrZero(key string) int
	Has(key string) bool
	HasBytes(key []byte) bool
	Len() int
	Parse(s string)
	ParseBytes(b []byte)
	Peek(key string) []byte
	PeekBytes(key []byte) []byte
	PeekMulti(key string) [][]byte
	PeekMultiBytes(key []byte) [][]byte
	QueryString() []byte
	Reset()
	Set(key, value string)
	SetBytesK(key []byte, value string)
	SetBytesKV(key, value []byte)
	SetBytesV(key string, value []byte)
	SetUint(key string, value int)
	SetUintBytes(key []byte, value int)
	String() string
	VisitAll(f func(key, value []byte))
	WriteTo(w io.Writer) (int64, error)
}
*/
