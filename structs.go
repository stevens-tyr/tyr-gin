package tyrgin

import (
	"bufio"
	"bytes"

	"github.com/gin-gonic/gin"
)

// Creator Types/Structs

// The http type is a type for HTTP request types.
type http string

// The http request types.
const (
	GET    = http("GET")
	DELETE = http("DELETE")
	PATCH  = http("PATCh")
	POST   = http("POST")
	PUT    = http("PUT")
)

// APIAction is the core of how you can easily add routes to the server.
type APIAction struct {
	Func    func(gin *gin.Context)
	Route   string
	Private bool
	Method  http
}

// NewRoute takes a function that takes gin context, endpoint, whether the route should be login protected, and method type.
// This returns a pointer to a APIAction.
func NewRoute(action func(gin *gin.Context), endpoint string, private bool, method http) APIAction {
	return APIAction{
		Func:    action,
		Route:   endpoint,
		Private: private,
		Method:  method,
	}
}

// Logger Types/Structs

// bufferedWriter a writer to add on top of
type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}

// Write the function to make buferredWriter type part of go's
// Writer interface.
func (b *bufferedWriter) Write(data []byte) (int, error) {
	b.Buffer.Write(data)
	return b.out.Write(data)
}
