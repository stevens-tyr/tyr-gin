package tyrgin

import (
	"github.com/gin-gonic/gin"
)

// Creator Types

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
