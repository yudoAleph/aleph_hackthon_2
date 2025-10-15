package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type JSONLogEntry struct {
	Timestamp     string      `json:"@timestamp"`
	Level         string      `json:"level"`
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Status        int         `json:"status"`
	Latency       float64     `json:"latency_ms"` // in milliseconds
	ClientIP      string      `json:"client_ip"`
	UserAgent     string      `json:"user_agent"`
	ErrorMessage  string      `json:"error_message,omitempty"`
	RequestBody   interface{} `json:"request_body,omitempty"`
	ResponseBody  interface{} `json:"response_body,omitempty"`
	CorrelationID string      `json:"correlation_id,omitempty"`
	UserID        uint        `json:"user_id,omitempty"`
}

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
		},
	})

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Set up daily rotating file
	currentTime := time.Now()
	logFileName := filepath.Join(logsDir, fmt.Sprintf("app-%s.log", currentTime.Format("2006-01-02")))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Write to both file and stdout
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

// rotateLogFile creates a new log file for the current day
func rotateLogFile() {
	logsDir := "logs"
	currentTime := time.Now()
	logFileName := filepath.Join(logsDir, fmt.Sprintf("app-%s.log", currentTime.Format("2006-01-02")))

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
		return
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

// JSONLogMiddleware is a Gin middleware that logs requests in JSON format
func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if we need to rotate the log file
		currentTime := time.Now()
		if currentTime.Hour() == 0 && currentTime.Minute() == 0 {
			rotateLogFile()
		}

		// Read the request body
		var requestBody interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			_ = json.Unmarshal(bodyBytes, &requestBody)
		}

		// Create a custom response writer to capture the response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Parse response body
		var responseBody interface{}
		if len(blw.body.String()) > 0 {
			_ = json.Unmarshal(blw.body.Bytes(), &responseBody)
		}

		// Get user ID from context if available
		var userID uint
		if id, exists := c.Get("user_id"); exists {
			userID = id.(uint)
		}

		// Create log entry
		entry := &JSONLogEntry{
			Timestamp:    time.Now().Format(time.RFC3339),
			Level:        getLogLevel(c.Writer.Status()),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Status:       c.Writer.Status(),
			Latency:      float64(time.Since(start)) / float64(time.Millisecond),
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			RequestBody:  requestBody,
			ResponseBody: responseBody,
			UserID:       userID,
		}

		// Add correlation ID if present
		if corrID := c.GetHeader("X-Correlation-ID"); corrID != "" {
			entry.CorrelationID = corrID
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			entry.ErrorMessage = c.Errors.String()
		}

		// Log the entry
		logJSON, _ := json.Marshal(entry)
		log.Info(string(logJSON))
	}
}

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// getLogLevel returns the appropriate log level based on status code
func getLogLevel(status int) string {
	switch {
	case status >= 500:
		return "error"
	case status >= 400:
		return "warning"
	default:
		return "info"
	}
}

// Error logs an error message with context
func Error(err error, context map[string]interface{}) {
	log.WithFields(logrus.Fields(context)).Error(err)
}

// Info logs an info message with context
func Info(msg string, context map[string]interface{}) {
	log.WithFields(logrus.Fields(context)).Info(msg)
}

// Warn logs a warning message with context
func Warn(msg string, context map[string]interface{}) {
	log.WithFields(logrus.Fields(context)).Warn(msg)
}

// Debug logs a debug message with context
func Debug(msg string, context map[string]interface{}) {
	log.WithFields(logrus.Fields(context)).Debug(msg)
}

// LogEndpointError logs structured endpoint errors for Kibana
func LogEndpointError(c *gin.Context, handler string, err error, statusCode int, additionalContext map[string]interface{}) {
	context := map[string]interface{}{
		"handler":       handler,
		"method":        c.Request.Method,
		"path":          c.Request.URL.Path,
		"status_code":   statusCode,
		"client_ip":     c.ClientIP(),
		"user_agent":    c.Request.UserAgent(),
		"error_type":    "endpoint_error",
		"error_message": err.Error(),
		"@timestamp":    time.Now().Format(time.RFC3339),
	}

	// Add user ID if available
	if userID, exists := c.Get("user_id"); exists {
		context["user_id"] = userID
	}

	// Add correlation ID if present
	if corrID := c.GetHeader("X-Correlation-ID"); corrID != "" {
		context["correlation_id"] = corrID
	}

	// Add additional context
	for k, v := range additionalContext {
		context[k] = v
	}

	Error(err, context)
}

// LogEndpointTimeout logs timeout errors for Kibana
func LogEndpointTimeout(c *gin.Context, handler string, timeout time.Duration, additionalContext map[string]interface{}) {
	context := map[string]interface{}{
		"handler":         handler,
		"method":          c.Request.Method,
		"path":            c.Request.URL.Path,
		"status_code":     408,
		"client_ip":       c.ClientIP(),
		"user_agent":      c.Request.UserAgent(),
		"error_type":      "timeout_error",
		"error_message":   "Request timeout",
		"timeout_seconds": timeout.Seconds(),
		"@timestamp":      time.Now().Format(time.RFC3339),
	}

	// Add user ID if available
	if userID, exists := c.Get("user_id"); exists {
		context["user_id"] = userID
	}

	// Add correlation ID if present
	if corrID := c.GetHeader("X-Correlation-ID"); corrID != "" {
		context["correlation_id"] = corrID
	}

	// Add additional context
	for k, v := range additionalContext {
		context[k] = v
	}

	Warn("Request timeout", context)
}

// LogValidationError logs validation errors for Kibana
func LogValidationError(c *gin.Context, handler string, validationErrors map[string]string, additionalContext map[string]interface{}) {
	context := map[string]interface{}{
		"handler":           handler,
		"method":            c.Request.Method,
		"path":              c.Request.URL.Path,
		"status_code":       400,
		"client_ip":         c.ClientIP(),
		"user_agent":        c.Request.UserAgent(),
		"error_type":        "validation_error",
		"error_message":     "Validation failed",
		"validation_errors": validationErrors,
		"@timestamp":        time.Now().Format(time.RFC3339),
	}

	// Add user ID if available
	if userID, exists := c.Get("user_id"); exists {
		context["user_id"] = userID
	}

	// Add correlation ID if present
	if corrID := c.GetHeader("X-Correlation-ID"); corrID != "" {
		context["correlation_id"] = corrID
	}

	// Add additional context
	for k, v := range additionalContext {
		context[k] = v
	}

	Warn("Validation error", context)
}

// LogAuthError logs authentication/authorization errors for Kibana
func LogAuthError(c *gin.Context, handler string, err error, additionalContext map[string]interface{}) {
	context := map[string]interface{}{
		"handler":       handler,
		"method":        c.Request.Method,
		"path":          c.Request.URL.Path,
		"status_code":   401,
		"client_ip":     c.ClientIP(),
		"user_agent":    c.Request.UserAgent(),
		"error_type":    "auth_error",
		"error_message": err.Error(),
		"@timestamp":    time.Now().Format(time.RFC3339),
	}

	// Add correlation ID if present
	if corrID := c.GetHeader("X-Correlation-ID"); corrID != "" {
		context["correlation_id"] = corrID
	}

	// Add additional context
	for k, v := range additionalContext {
		context[k] = v
	}

	Warn("Authentication error", context)
}
