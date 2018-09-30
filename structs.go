package tyrgin

import (
	"bufio"
	"bytes"

	"github.com/gin-gonic/gin"
)

// Creator Types/Structs

// The http type is a type for HTTP request types.
type httpMethod string

// The http request types.
const (
	GET    httpMethod = "GET"
	DELETE httpMethod = "DELETE"
	PATCH  httpMethod = "PATCH"
	POST   httpMethod = "POST"
	PUT    httpMethod = "PUT"
)

// Default Error Messages
const (
	NotFoundError = "404 PAGE NOT FOUND"
)

// APIAction is the core of how you can easily add routes to the server.
type APIAction struct {
	Func    func(gin *gin.Context)
	Route   string
	Private bool
	Method  httpMethod
}

// NewRoute takes a function that takes gin context, endpoint, whether the route should be login protected, and method type.
// This returns a pointer to a APIAction.
func NewRoute(action func(gin *gin.Context), endpoint string, private bool, method httpMethod) APIAction {
	return APIAction{
		Func:    action,
		Route:   endpoint,
		Private: private,
		Method:  method,
	}
}

// About Check Types/Structs
const (
	AboutFieldNa      string = "N/A"
	AboutProtocolHttp string = "http"
	VersionNa         string = "N/A"
)

type (
	ConfigAbout struct {
		Id          string                 `json:"id"`
		Summary     string                 `json:"sumamry"`
		Description string                 `json:"description"`
		Maintainers []string               `json:"maintainers"`
		ProjectRepo string                 `json:"projectRepo"`
		ProjectHome string                 `json:"projectHome"`
		LogsLinks   []string               `json:"logsLinks"`
		StatsLinks  []string               `json:"statsLinks"`
		CustomData  map[string]interface{} `json:"customData"`
	}

	AboutResponse struct {
		Id           string                 `json:"id"`
		Name         string                 `json:"name"`
		Description  string                 `json:"description"`
		Protocol     string                 `json:"protocol"`
		Owners       []string               `json:"owners"`
		Version      string                 `json:"version"`
		Host         string                 `json:"host"`
		ProjectRepo  string                 `json:"projectRepo"`
		ProjectHome  string                 `json:"projectHome"`
		LogsLinks    []string               `json:"logsLinks"`
		StatsLinks   []string               `json:"statsLinks"`
		Dependencies []Dependency           `json:"dependencies"`
		CustomData   map[string]interface{} `json:"customData"`
	}

	Dependency struct {
		Name           string         `json:"name"`
		Status         []JsonResponse `json:"status"`
		StatusDuration float64        `json:"statusDuration"`
		StatusPath     string         `json:"statusPath"`
		Type           string         `json:"type"`
		IsTraversable  bool           `json:"isTraversable"`
	}

	dependencyPosition struct {
		item     Dependency
		position int
	}
)

// Health Check Types/Structs
type AlertLevel string

const (
	OK       AlertLevel = "OK"
	WARNING  AlertLevel = "WARN"
	CRITICAL AlertLevel = "CRIT"
)

type (
	StatusResponse struct {
		Status string `jon:"status"`
	}

	StatusEndpoint struct {
		Name          string
		Slug          string
		Type          string
		IsTraversable bool
		StatusCheck   StatusCheck
		TraverseCheck TraverseCheck
	}

	Status struct {
		Description string     `json:"description"`
		Result      AlertLevel `json:"result"`
		Details     string     `json:"details"`
	}

	StatusList struct {
		StatusList []Status
	}

	StatusCheck interface {
		CheckStatus(name string) StatusList
	}

	JsonResponse interface{}

	TraverseCheck interface {
		Traverse(traversalPath []string, action string) (string, error)
	}

	MongoStatusChecker struct{}
)

// Logger Types/Structs

// bufferedWriter a writer to add on top of
type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}
