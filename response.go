package lessgo

import (
	"bufio"
	"net"
	"net/http"
)

// Response wraps an http.ResponseWriter and implements its interface to be used
// by an HTTP handler to construct an HTTP response.
// See [http.ResponseWriter](https://golang.org/pkg/net/http/#ResponseWriter)
type Response struct {
	writer    http.ResponseWriter
	status    int
	size      int64
	committed bool
}

// NewResponse creates a new instance of Response.
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{writer: w}
}

// SetWriter sets the http.ResponseWriter instance for this Response.
func (resp *Response) SetWriter(w http.ResponseWriter) {
	resp.writer = w
}

// Writer returns the http.ResponseWriter instance for this Response.
func (resp *Response) Writer() http.ResponseWriter {
	return resp.writer
}

// Header returns the header map for the writer that will be sent by
// WriteHeaderesp. Changing the header after a call to WriteHeader (or Write) has
// no effect unless the modified headers were declared as trailers by setting
// the "Trailer" header before the call to WriteHeader (see example)
// To suppress implicit response headers, set their value to nil.
// Example [ResponseWriteresp.Trailers](https://golang.org/pkg/net/http/#example_ResponseWriter_trailers)
func (resp *Response) Header() http.Header {
	return resp.writer.Header()
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (resp *Response) WriteHeader(code int) {
	if resp.committed {
		Log.Warn("response already committed")
		return
	}
	resp.status = code
	resp.writer.WriteHeader(code)
	resp.committed = true
}

// Write wraps and implements the http.Response.Write specification.
// Additionally, Write will increment the size of the current response.
// See [http.Response.Write](https://golang.org/pkg/net/http/#Response.Write)
func (resp *Response) Write(b []byte) (int, error) {
	n, err := resp.writer.Write(b)
	resp.size += int64(n)
	return n, err
}

// SetCookie adds a Set-Cookie header.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (r *Response) SetCookie(cookie *http.Cookie) {
	r.Header().Add("Set-Cookie", cookie.String())
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (resp *Response) Flush() {
	resp.writer.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (resp *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return resp.writer.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotifier interface to allow detecting
// when the underlying connection has gone away.
// This mechanism can be used to cancel long operations on the server if the
// client has disconnected before the response is ready.
// See [http.CloseNotifier](https://golang.org/pkg/net/http/#CloseNotifier)
func (resp *Response) CloseNotify() <-chan bool {
	return resp.writer.(http.CloseNotifier).CloseNotify()
}

// Status returns the HTTP status code of the response.
func (resp *Response) Status() int {
	return resp.status
}

// Size returns the current size, in bytes, of the response.
func (resp *Response) Size() int64 {
	return resp.size
}

// Committed asserts whether or not the response has been committed to.
func (resp *Response) Committed() bool {
	return resp.committed
}

func (resp *Response) free() {
	resp.writer = nil
	resp.size = 0
	resp.status = http.StatusOK
	resp.committed = false
}
