package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogFormatterParams struct {

	// StatusCode is HTTP response code.
	StatusCode int `json:"status_code"`
	// Latency is how much time the server cost to process a certain request.
	Latency int64 `json:"latency"`
	// ClientIP equals Context's ClientIP method.
	ClientIP string `json:"client_ip"`
	// Method is the HTTP method given to the request.
	Method string `json:"method"`
	// Path is a path the client requests.
	Path string `json:"path"`
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string `json:"error_message"`
	// isTerm shows whether gin's output descriptor refers to a terminal.
	// BodySize is the size of the Response Body
	BodySize int    `json:"body_size"`
	User     string `json:"user_id"`
	// Keys are the keys set on the request's context.
	Keys map[string]any `json:"keys,omitempty"`
}

var formatJSON = func(log LogFormatterParams) string {
	logline, _ := json.Marshal(log)
	return fmt.Sprintf("%s\n", logline)
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(logger *logrus.Logger) func(c *gin.Context) {

	return func(c *gin.Context) {
		start := time.Now().UnixNano() / int64(time.Millisecond)
		path := c.Request.URL.Path
		// Process request
		c.Next()
		end := time.Now().UnixNano() / int64(time.Millisecond)
		param := LogFormatterParams{}
		// Stop timer
		param.Latency = end - start
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		param.User = c.GetString("user_id")
		param.Path = path
		//fmt.Fprint(os.Stdout, formatJSON(param))
	}
}
