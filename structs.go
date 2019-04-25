package tyrgin

import (
	"bufio"
	"bytes"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/gridfs"
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

// Errors
var (
	// ErrorEmailNotValid an error to throw when an email format is not valid
	ErrorEmailNotValid = errors.New("EMAIL NOT VALID")
	// ErrorUnresolvedEmailHost an error to throw when the email host is unresolvable
	ErrorUnresolvableEmailHost = errors.New("EMAIL HOST UNRESOLVABLE")
	// UserNotFoundError an error to throw for when a User is not found.
	ErrorUserNotFound = errors.New("USER NOT FOUND")
	// IncorrectPasswordError an error to throw for when an inccorect passowrd is entered.
	ErrorIncorrectPassword = errors.New("INCORRECT PASSWORD")
	// MongoSessionFailure an error to throw for when a mongo Session fails.
	ErrorMongoSessionFailure = errors.New("FAILED TO GET MONGO SESSION")
	// MongoCollectionFailure an error to throw for when a mongo collection does not exist.
	ErrorMongoCollectionFailure = errors.New("MONGO COLLECTION DOES NOT EXIST")
)

// APIAction is the core of how you can easily add routes to the server.
type APIAction struct {
	Func   func(gin *gin.Context)
	Route  string
	Method httpMethod
}

// NewRoute takes a function that takes gin context, endpoint, whether the route should be login protected, and method type.
// This returns a pointer to a APIAction.
func NewRoute(action func(c *gin.Context), endpoint string, method httpMethod) APIAction {
	return APIAction{
		Func:   action,
		Route:  endpoint,
		Method: method,
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

	// MongoReplStatus a struct to unpack message from checking mongo replicatset.
	MongoReplStatus struct {
		OK       int    `bson:"ok" binding:"required"`
		ErrorMsg string `bson:"errmsg" binding:"required"`
	}

	// MongoRPLStatusChecker struct for when we eventually add mongo.
	MongoRPLStatusChecker struct {
		RPL *mongo.Database
	}
)

// Logger Types/Structs

// bufferedWriter a writer to add on top of
type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}

// GridFS Types/Structs

// Bucket struct contains the gridfs Bucket as well its its chunk size and name.
type Bucket struct {
	Bucket         *gridfs.Bucket
	ChunkSizeBytes *int32
	Name           *string
}
