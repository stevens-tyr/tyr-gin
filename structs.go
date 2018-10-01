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

// Default Fields
const (
	AboutFieldNa      string = "N/A"
	AboutProtocolHTTP string = "http"
	VersionNa         string = "N/A"
)

type (
	// ConfigAbout this struct to configure an About API call.
	ConfigAbout struct {
		ID          string                 `json:"id"`
		Summary     string                 `json:"sumamry"`
		Description string                 `json:"description"`
		Maintainers []string               `json:"maintainers"`
		ProjectRepo string                 `json:"projectRepo"`
		ProjectHome string                 `json:"projectHome"`
		LogsLinks   []string               `json:"logsLinks"`
		StatsLinks  []string               `json:"statsLinks"`
		CustomData  map[string]interface{} `json:"customData"`
	}

	// AboutResponse is the response the Aboout API response.
	AboutResponse struct {
		ID           string                 `json:"id"`
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

	// Dependency is the dependency struct to go inside the AboutResponse struct
	// to detail the dependencies of the service.
	Dependency struct {
		Name           string         `json:"name"`
		Status         []JSONResponse `json:"status"`
		StatusDuration float64        `json:"statusDuration"`
		StatusPath     string         `json:"statusPath"`
		Type           string         `json:"type"`
		IsTraversable  bool           `json:"isTraversable"`
	}

	// dependencyPosition is a simple struct to keep track of where dependencies are
	// in a slice of dependencyPositions.
	dependencyPosition struct {
		item     Dependency
		position int
	}
)

// Health Check Types/Structs

// AlertLevel is a string wrapper to tell the Health check the alert level of items.
type AlertLevel string

// Some constanst to keep track of logging level
const (
	OK       AlertLevel = "OK"
	WARNING  AlertLevel = "WARN"
	CRITICAL AlertLevel = "CRIT"
)

type (
	// StatusResponse is just a struct to represent the status of responses and
	// goes inside the StatusEndpoint struct.
	StatusResponse struct {
		Status string `jon:"status"`
	}

	// StatusEndpoint is the struct to track information about an endpoints status.
	StatusEndpoint struct {
		Name          string
		Slug          string
		Type          string
		IsTraversable bool
		StatusCheck   StatusCheck
		TraverseCheck TraverseCheck
	}

	// Status struct keeps track of the details, description and results of a status.
	Status struct {
		Description string     `json:"description"`
		Result      AlertLevel `json:"result"`
		Details     string     `json:"details"`
	}

	// StatusList to convienently pass around a slice of the Status struct.
	StatusList struct {
		StatusList []Status
	}

	// StatusCheck itnerface this way when we add new services we can have it implement this
	// interace so that we can just call this function on the StatusCheck field from
	// the StatusEndpoint struct.
	StatusCheck interface {
		CheckStatus(name string) StatusList
	}

	// JSONResponse is a wrapper of the interface{} type.
	JSONResponse interface{}

	// TraverseCheck interface this way when we add new services we can have it implement this
	// interace so that we can just call this function on the TraverseCheck field from
	// the StatusEndpoint struct.
	TraverseCheck interface {
		Traverse(traversalPath []string, action string) (string, error)
	}

	// MongoStatusChecker struct for when we eventually add mongo.
	MongoStatusChecker struct{}
)

// Logger Types/Structs

// bufferedWriter a writer to add on top of
type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}
