package tyrgin

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"time"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// determineLogLevel given a string of the level it converts to a
// logrus logger level and returns that type to be used. Defaults
// to debug level.
func determineLogLevel(level string) log.Level {
	switch level {
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.DebugLevel
	}
}

func init() {
	// Format to json because jq tool is amazing.
	formatter := runtime.Formatter{ChildFormatter: &log.JSONFormatter{}}
	formatter.Line = true
	log.SetFormatter(&formatter)

	// Grab log file name from enviornment. Default is log.json.
	logFileName := os.Getenv("LOG_FILE")
	if logFileName == "" {
		logFileName = "log.json"
	}
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Could not create log file. Server Quiting...")
	}
	log.SetOutput(logFile)

	// Grab log level from env.
	logLevelEnv := os.Getenv("LOG_LEVEL")
	log.SetLevel(determineLogLevel(logLevelEnv))
	log.Info("Logging starting...")
}

// Write the function to make buferredWriter type part of go's
// Writer interface.
func (b *bufferedWriter) Write(data []byte) (int, error) {
	b.Buffer.Write(data)
	return b.out.Write(data)
}

// ErrorLogger takes an error and a message, if the error is not
// null log with warning message.
func ErrorLogger(err error, msg string) {
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"msg":   msg,
		}).Warn("Code Error")
	}
}

// Logger a logging middleware to be used with gin.
// Logs standard information based of the information given.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()

		// Give our context additional context to see response body.
		w := bufio.NewWriter(c.Writer)
		buff := bytes.Buffer{}
		newWriter := &bufferedWriter{c.Writer, w, buff}

		c.Writer = newWriter

		// You have to manually flush the buffer at the end
		defer func() {
			w.Flush()
		}()

		bytesBody, err := ioutil.ReadAll(c.Request.Body)
		ErrorLogger(err, "Failed to read Request Body.")

		c.Next()
		//after request
		latency := int64(time.Since(t) / time.Millisecond)

		contextLog := log.WithFields(log.Fields{
			"RequestMethod":   c.Request.Method,
			"RequestUrl":      c.Request.URL,
			"RequestHeaders":  c.Request.Header,
			"RequestBody":     string(bytesBody),
			"ResponseStatus":  c.Writer.Status(),
			"Latency(ms)":     latency,
			"ResponseHeaders": c.Writer.Header(),
			"ResponseBody":    string(newWriter.Buffer.Bytes()),
		})

		// Should never panic or fatal because that will exit server.
		switch c.Writer.Status() {
		case 200, 201, 202:
			contextLog.Info("OK")
			break
		case 400, 401, 403, 404, 405, 415:
			contextLog.Error("Client Error")
		case 500:
			contextLog.Error("Server Error")
		default:
			contextLog.Info("Unexpected")
		}
	}
}
