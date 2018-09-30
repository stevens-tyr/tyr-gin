package tyrgin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func translateStatusList(s StatusList) []JsonResponse {

	if len(s.StatusList) <= 0 {
		return []JsonResponse{
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
		return []JsonResponse{
			OK,
		}
	} else {
		return []JsonResponse{
			r.Result,
			r,
		}
	}

}

func SerializeStatusList(s StatusList) string {

	statusListJsonResponse := translateStatusList(s)
	statusListJson, err := json.Marshal(statusListJsonResponse)
	ErrorLogger(err, fmt.Sprintf(`["CRIT", {"description":"Invalid StatusList","result":"CRIT","detials":"Error serializing StatusList: %v error:%s"}]`, s, err))

	return string(statusListJson)
}

func ExecuteStatusCheck(s *StatusEndpoint) string {
	result := s.StatusCheck.CheckStatus(s.Name)
	return SerializeStatusList(result)
}

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

// HandlerPointHandlerFunc returns a http.HandlerFunc that responds to status check requests. It should be registered at `/status/...`
func HealthPointHandlerFunc(statusEndpoints []StatusEndpoint, aboutFilePath, versionFilePath string, customData map[string]interface{}) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		slug := strings.Split(r.URL.Path, "/")
		endpoint := slug[2]

		switch endpoint {
		case "about":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			io.WriteString(w, About(statusEndpoints, AboutProtocolHttp, aboutFilePath, versionFilePath, customData))

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
			io.WriteString(w, Traverse(statusEndpoints, dependencies, action, AboutProtocolHttp, aboutFilePath, versionFilePath, customData))
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
