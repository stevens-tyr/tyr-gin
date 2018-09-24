package tyrgin

import (
	"os"
	"time"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	formatter := runtime.Formatter{ChildFormatter: &log.JSONFormatter{}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// Logger a logging middleware to be used with gin.
// Logs standard information based of the information given.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()

		c.Next()
		//after request
		latency := time.Since(t)
		contextLog := log.WithFields(log.Fields{
			"RequestMethod":   c.Request.Method,
			"RequestUrl":      c.Request.URL,
			"RequestHeaders":  c.Request.Header,
			"RequestBody":     c.Request.Body,
			"ResponseStatus":  c.Request.Response.Status,
			"Latency":         latency,
			"ResponseHeaders": c.Request.Response.Header,
			"ResponseBody":    c.Request.Response.Body,
		})

		switch c.Request.Response.Status {
		default:
			contextLog.Info("temp")
			break
		}
	}
}
