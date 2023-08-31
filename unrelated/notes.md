# 1. Intro

```go
package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}
```

## From `server.go`

### `type Handler interface`

- responds to an HTTP request

```go
type Handler {
  ServeHTTP(ResponseWriter, *Request)
}
```

- `ServeHTTP` should write reply headers and data to the `ResponseWriter` and the return.

### `type ResponseWriter interface`

```go
type ResponseWriter interface {
  Header() Header
  Write([]byte) (int, error)
  WriteHeader(statusCode int)
}
```

- `Header` returns the header map that will be sent by `WriteHeader`
- `Write`
  - writes the data to the connection as part of an HTTP reply. If `WriteHeader` ahs not yet been called, `Write` calls `WriteHead(http.StatusOK)` before writing the data.
  - if the `Header` doesn't contain a `Content-Type` line, `Write` adds a `Content-Type` set to the result passing the initial 512 bytes of written data to `DetectContentType`
- `WriteHeader`
  - sends an HTTP response with the provided status code.
  - provided code must be a valid HTTP 1xx-5xx status code

### `type Flusher interface`

- is implemented by `ResponseWriters`, allows an HTTP handler to flush buffered data to the client.

```go
type Flusher interface {
  Flush()
}
```

- `Flush` sends any buffered data to the client

### `type Hijacker interface`

- is implemented by `ResponseWriters` that allows an HTTP handler to take over the connection.

```go
type Hijacker interface {
  Hijack() (net.Conn, *buffio.ReadWriter, error)
}
```

### `type CloseNotifier interface`

- returns a channel that receives at most a single value (true) when the client connection has gone away.

```go
type CloseNotifier interface {
  CloseNotify() <-chan bool
}
```

### `type conn struct`

```go
// A conn represents the server side of an HTTP connection.
type conn struct {
	// server is the server on which the connection arrived.
	// Immutable; never nil.
	server *Server

	// cancelCtx cancels the connection-level context.
	cancelCtx context.CancelFunc

	// rwc is the underlying network connection.
	// This is never wrapped by other types and is the value given out
	// to CloseNotifier callers. It is usually of type *net.TCPConn or
	// *tls.Conn.
	rwc net.Conn

	// remoteAddr is rwc.RemoteAddr().String(). It is not populated synchronously
	// inside the Listener's Accept goroutine, as some implementations block.
	// It is populated immediately inside the (*conn).serve goroutine.
	// This is the value of a Handler's (*Request).RemoteAddr.
	remoteAddr string

	// tlsState is the TLS connection state when using TLS.
	// nil means not TLS.
	tlsState *tls.ConnectionState

	// werr is set to the first write error to rwc.
	// It is set via checkConnErrorWriter{w}, where bufw writes.
	werr error

	// r is bufr's read source. It's a wrapper around rwc that provides
	// io.LimitedReader-style limiting (while reading request headers)
	// and functionality to support CloseNotifier. See *connReader docs.
	r *connReader

	// bufr reads from r.
	bufr *bufio.Reader

	// bufw writes to checkConnErrorWriter{c}, which populates werr on error.
	bufw *bufio.Writer

	// lastMethod is the method of the most recent request
	// on this connection, if any.
	lastMethod string

	curReq atomic.Pointer[response] // (which has a Request in it)

	curState atomic.Uint64 // packed (unixtime<<8|uint8(ConnState))

	// mu guards hijackedv
	mu sync.Mutex

	// hijackedv is whether this connection has been hijacked
	// by a Handler with the Hijacker interface.
	// It is guarded by mu.
	hijackedv bool
}
```

### `type response struct`

### `type ServeMux struct`

```go
// ServeMux is an HTTP request multiplexer.
// It matches the URL of each incoming request against a list of registered
// patterns and calls the handler for the pattern that
// most closely matches the URL.
//
// Patterns name fixed, rooted paths, like "/favicon.ico",
// or rooted subtrees, like "/images/" (note the trailing slash).
// Longer patterns take precedence over shorter ones, so that
// if there are handlers registered for both "/images/"
// and "/images/thumbnails/", the latter handler will be
// called for paths beginning with "/images/thumbnails/" and the
// former will receive requests for any other paths in the
// "/images/" subtree.
//
// Note that since a pattern ending in a slash names a rooted subtree,
// the pattern "/" matches all paths not matched by other registered
// patterns, not just the URL with Path == "/".
//
// If a subtree has been registered and a request is received naming the
// subtree root without its trailing slash, ServeMux redirects that
// request to the subtree root (adding the trailing slash). This behavior can
// be overridden with a separate registration for the path without
// the trailing slash. For example, registering "/images/" causes ServeMux
// to redirect a request for "/images" to "/images/", unless "/images" has
// been registered separately.
//
// Patterns may optionally begin with a host name, restricting matches to
// URLs on that host only. Host-specific patterns take precedence over
// general patterns, so that a handler might register for the two patterns
// "/codesearch" and "codesearch.google.com/" without also taking over
// requests for "http://www.google.com/".
//
// ServeMux also takes care of sanitizing the URL request path and the Host
// header, stripping the port number and redirecting any request containing . or
// .. elements or repeated slashes to an equivalent, cleaner URL.
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}
```

### We can use `ServeMux` like this:

```go
package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", mux)
}
```

### `func ListenAndServe(addr string, handler handler) error`

# 2. Dynamic Reloading

`$ go install github.com/cortesi/modd/cmd/modd@latest`

- we install it locally

then, create `modd.conf`:

```.conf
**/*.go {
  prep: go test @dirmods
}
```

- any time any go file changes, run `go test`

```.conf
**/*.go {
  prep: go test @dirmods
}

**/*.go !**/*_test.go {
  prep: go build -o webapp .
  daemon +sigterm: ./webapp
}
```

# 3. Content Type Headers

The `Content-Type` representation header is used to indicate the original media type of the resource (prior to any content encoding applied for sending).

```go
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	// ...
}

```

# 4. Getting url path:

```go
func pathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.URL.Path)
}
```

# 5. URL.Path vs URL.RawPath

```go
type URL struct {
	Scheme      string
	Opaque      string    // encoded opaque data
	User        *Userinfo // username and password information
	Host        string    // host or host:port
	Path        string    // path (relative paths may omit leading slash)
	RawPath     string    // encoded path hint (see EscapedPath method)
	OmitHost    bool      // do not emit empty host (authority)
	ForceQuery  bool      // append a query ('?') even if RawQuery is empty
	RawQuery    string    // encoded query values, without '?'
	Fragment    string    // fragment for references, without '#'
	RawFragment string    // encoded fragment hint (see EscapedFragment method)
}
```

# 6. Not Found Handler

```go
http.Error(w, "Page not found", http.StatusNotFound)

// or
http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
```

# 7. Router

```go
package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // redundant
	w.WriteHeader(http.StatusOK) // redundant
	fmt.Fprint(w, "<h1>Welcome to my great site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:john@example.com\">john</a></p>")
}

type Router struct {}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request ) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}
```

# 8. http.HandlerFunc conversion

```go
func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func main() {
	var router http.HandlerFunc = pathHandler
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		panic(err)
	}
}
```

or even like this:

```go
func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func main() {
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", http.HandlerFunc(pathHandler))
	if err != nil {
		panic(err)
	}
}
```

- ` http.HandlerFunc(pathHandler)` is not a function call! This is almost like type conversion!

```go
types:

-	http.Handler - interface with the ServeHTTP method
- http.HandlerFunc - a function type that accepts same args as ServeHTTP method,
	- also implements http.Handler

http.HandlerFunc(pathHandler)

// vs

http.Handle("/", http.Handler)
http.HandleFunc("/", pathHandler)
```

# 9. http.HandleFunc

- Registers the handler function for the given pattern in the DefaultServeMux

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

compare this two:

```go
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/contact", contactHandler)
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
```

with this one

```go
func main() {
	http.Handle("/", http.HandlerFunc(homeHandler))
	http.Handle("/contact", http.HandlerFunc(contactHandler))
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
```

- `http.HandlerFunc` is provided to us to simplify things (very slightly)

we can even do this stuff:

```go
func main() {
	http.HandleFunc("/", http.HandlerFunc(homeHandler).ServeHTTP)
	http.HandleFunc("/contact", http.HandlerFunc(contactHandler).ServeHTTP)
	fmt.Println("Starting the server on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
```

- to simplify things: handle is for handlers
- and handle func is for functions that look like handle funcs
