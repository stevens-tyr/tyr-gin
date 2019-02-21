package tyrgin

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	godotenv "github.com/joho/godotenv"
)

// MongoTyrRSStatusEndpoint is for healthcheck api to know about mongo replica sets.
var MongoTyrRSStatusEndpoint StatusEndpoint

func init() {
	env := os.Getenv("ENV")
	if env == "" {
		os.Setenv("ENV", "dev")
		env = "dev"
	}

	if env == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Could not load .env file.")
		}
	}

	checkSession, err := GetMongoDB(os.Getenv("DB_NAME"))
	if err != nil {
		log.Println("Could not get Mongo connection")
	}

	MongoTyrRSStatusEndpoint = StatusEndpoint{
		Name:          "Mongo Tyr Replica Set Check",
		Slug:          "mongo",
		Type:          "internal",
		IsTraversable: false,
		StatusCheck: MongoRPLStatusChecker{
			RPL: checkSession,
		},
		TraverseCheck: nil,
	}
}

// action takes the APIAction method and creates a gin route of that type.
// Also makes the route private if it labeled as private in the apiaction.
func (a *APIAction) action(route *gin.RouterGroup, jwt *jwt.GinJWTMiddleware) {
	if a.Private {
		
	}

	switch a.Method {
	case GET:
		route.GET(a.Route, a.Func)
		break
	case DELETE:
		route.DELETE(a.Route, a.Func)
		break
	case PATCH:
		route.PATCH(a.Route, a.Func)
		break
	case POST:
		route.POST(a.Route, a.Func)
		break
	case PUT:
		route.PUT(a.Route, a.Func)
		break
	}
}

// AddRoutes takes a gin server, gin jwt instance, version number as a string,
// api endpoint name and a list of APIActions to add to it.
func AddRoutes(router *gin.Engine, jwt *jwt.GinJWTMiddleware, version, api string, fns []APIAction) {

	ver := router.Group("/api/v" + version)
	{
		route := ver.Group(api)
		{

			for _, fn := range fns {
				fn.action(route, jwt)
			}

		}
	}

}

// NotFound a general 404 error message.
func NotFound(c *gin.Context) {
	fmt.Println(c.Request.URL.Path[1:])
	c.JSON(
		http.StatusNotFound,
		gin.H{
			"status_code": http.StatusNotFound,
			"message":     NotFoundError,
		},
	)
}

// SetupRouter returns an instance to a *gin.Enginer that is has
// some preconfigurations already set up.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(Logger())
	router.Use(gin.Recovery())
	route.Use(jwt.MiddlewareFunc())

	router.GET(
		"/status/:slug",
		HealthPointHandler(
			[]StatusEndpoint{
				MongoTyrRSStatusEndpoint,
			},
			"./about.json",
			"./version.txt",
			make(map[string]interface{}),
		),
	)

	return router
}

// ServeReact is a function to serve react from a(n) service.
func ServeReact(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./static", true)))
}

// ErrorHandler handles gin errors in a more clean way
func ErrorHandler(err error, c *gin.Context, sc int, json interface{}) {
	c.Writer.Header().Add("Content-Type", "application/json+error")
	c.AbortWithStatusJSON(sc, json)
	c.Error(err)
	return
}
