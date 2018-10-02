package tyrgin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// translateStatusLists takes a StatusList and converts it to a slice of JSONResponse.
func translateStatusList(s StatusList) []JSONResponse {

	if len(s.StatusList) <= 0 {
		return []JSONResponse{
			CRITICAL,
			Status{
				Description: "Invalid status response",
				Result:      CRITICAL,
				Details:     "StatusList empty",
			},
		}
	}

	r := s.StatusList[0]

	if r.Result == OK {
		return []JSONResponse{
			OK,
		}
	}

	return []JSONResponse{
		r.Result,
		r,
	}

}

// SerializeStatusList serializes a status list.
func SerializeStatusList(s StatusList) string {

	statusListJSONResponse := translateStatusList(s)
	statusListJSON, err := json.Marshal(statusListJSONResponse)
	ErrorLogger(err, fmt.Sprintf(`["CRIT", {"description":"Invalid StatusList","result":"CRIT","detials":"Error serializing StatusList: %v error:%s"}]`, s, err))

	return string(statusListJSON)
}

// ExecuteStatusCheck takes a StatusEndpoint and calls the StatusCheck CheckStatus function
// for the StatusCheck field.
func ExecuteStatusCheck(s *StatusEndpoint) string {
	result := s.StatusCheck.CheckStatus(s.Name)
	return SerializeStatusList(result)
}

// FindStatusEndpoint finds a StatusEndpoint from a slice of them and returns the one.
func FindStatusEndpoint(statusEndpoints []StatusEndpoint, slug string) *StatusEndpoint {

	if slug == "" {
		return nil
	}

	for _, se := range statusEndpoints {
		if slug == se.Slug {
			return &se
		}
	}

	return nil
}

// HealthPointHandler returns a gin.HandlerFunc that responds to status check requests. It should be registered at `/status/...`
func HealthPointHandler(statusEndpoints []StatusEndpoint, aboutFilePath, versionFilePath string, customData map[string]interface{}) gin.HandlerFunc {
	handler := HealthPointHandlerFunc(statusEndpoints, aboutFilePath, versionFilePath, customData)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// HealthPointHandlerFunc returns a http.HandlerFunc that responds to status check requests. It should be registered at `/status/...`
func HealthPointHandlerFunc(statusEndpoints []StatusEndpoint, aboutFilePath, versionFilePath string, customData map[string]interface{}) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		slug := strings.Split(r.URL.Path, "/")
		endpoint := slug[2]

		switch endpoint {
		case "about":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			io.WriteString(w, About(statusEndpoints, AboutProtocolHTTP, aboutFilePath, versionFilePath, customData))

		case "aggregate":
			typeFilter := r.URL.Query().Get("type")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			io.WriteString(w, Aggregate(statusEndpoints, typeFilter))

		case "am-i-up":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			io.WriteString(w, "OK")

		case "traverse":
			action := r.URL.Query().Get("action")
			if action == "" {
				action = "about"
			}
			dependencies := []string{}
			queryDependencies := r.URL.Query().Get("dependencies")
			if queryDependencies != "" {
				dependencies = strings.Split(queryDependencies, ",")
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			io.WriteString(w, Traverse(statusEndpoints, dependencies, action, AboutProtocolHTTP, aboutFilePath, versionFilePath, customData))
		default:
			endpoint := FindStatusEndpoint(statusEndpoints, endpoint)
			if endpoint == nil {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, SerializeStatusList(StatusList{
					StatusList: []Status{
						{
							Description: "Unknow Status endpoint",
							Result:      CRITICAL,
							Details:     fmt.Sprintf("Status endpoint does not exist: %s", r.URL.Path),
						},
					},
				}))
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			io.WriteString(w, ExecuteStatusCheck(endpoint))
		}

	})

}
