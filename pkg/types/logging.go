package types

import (
	"fmt"
)

// LoggingLevel represents the severity of log messages
type LoggingLevel string

const (
    LogLevelDebug     LoggingLevel = "debug"
    LogLevelInfo      LoggingLevel = "info"
    LogLevelNotice    LoggingLevel = "notice"
    LogLevelWarning   LoggingLevel = "warning"
    LogLevelError     LoggingLevel = "error"
    LogLevelCritical  LoggingLevel = "critical"
    LogLevelAlert     LoggingLevel = "alert"
    LogLevelEmergency LoggingLevel = "emergency"
)

// SetLevelRequestOption configures SetLevelRequest
type SetLevelRequestOption func(*SetLevelRequest) error

// SetLevelRequest represents a request to set logging level
type SetLevelRequest struct {
    Method string           `json:"method"`
    Params SetLevelParams   `json:"params"`
}

type SetLevelParams struct {
    Level LoggingLevel `json:"level"`
}

func NewSetLevelRequest(level LoggingLevel, opts ...SetLevelRequestOption) (*SetLevelRequest, error) {
    if err := validateLoggingLevel(level); err != nil {
        return nil, err
    }

    req := &SetLevelRequest{
        Method: "logging/setLevel",
        Params: SetLevelParams{
            Level: level,
        },
    }

    for _, opt := range opts {
        if err := opt(req); err != nil {
            return nil, fmt.Errorf("applying set level request option: %w", err)
        }
    }

    return req, nil
}

func validateLoggingLevel(level LoggingLevel) error {
    switch level {
    case LogLevelDebug, LogLevelInfo, LogLevelNotice, LogLevelWarning,
         LogLevelError, LogLevelCritical, LogLevelAlert, LogLevelEmergency:
        return nil
    default:
        return fmt.Errorf("invalid logging level: %s", level)
    }
}

// LoggingMessageOption configures LoggingMessage
type LoggingMessageOption func(*LoggingMessageNotification) error

// LoggingMessageNotification represents a log message notification
type LoggingMessageNotification struct {
    Method string                  `json:"method"`
    Params LoggingMessageParams    `json:"params"`
}

type LoggingMessageParams struct {
    Level  LoggingLevel  `json:"level"`
    Data   interface{}   `json:"data"`
    Logger *string       `json:"logger,omitempty"`
}

func NewLoggingMessage(level LoggingLevel, data interface{}, opts ...LoggingMessageOption) (*LoggingMessageNotification, error) {
    if err := validateLoggingLevel(level); err != nil {
        return nil, err
    }

    if data == nil {
        return nil, fmt.Errorf("log data cannot be nil")
    }

    msg := &LoggingMessageNotification{
        Method: "notifications/message",
        Params: LoggingMessageParams{
            Level: level,
            Data:  data,
        },
    }

    for _, opt := range opts {
        if err := opt(msg); err != nil {
            return nil, fmt.Errorf("applying logging message option: %w", err)
        }
    }

    return msg, nil
}

// LoggingMessage options

func WithLogger(logger string) LoggingMessageOption {
    return func(msg *LoggingMessageNotification) error {
        if logger == "" {
            return fmt.Errorf("logger name cannot be empty")
        }
        msg.Params.Logger = &logger
        return nil
    }
}

// Helper functions for creating log messages with specific levels

func NewDebugMessage(data interface{}, opts ...LoggingMessageOption) (*LoggingMessageNotification, error) {
    return NewLoggingMessage(LogLevelDebug, data, opts...)
}

func NewInfoMessage(data interface{}, opts ...LoggingMessageOption) (*LoggingMessageNotification, error) {
    return NewLoggingMessage(LogLevelInfo, data, opts...)
}

func NewWarningMessage(data interface{}, opts ...LoggingMessageOption) (*LoggingMessageNotification, error) {
    return NewLoggingMessage(LogLevelWarning, data, opts...)
}

func NewErrorMessage(data interface{}, opts ...LoggingMessageOption) (*LoggingMessageNotification, error) {
    return NewLoggingMessage(LogLevelError, data, opts...)
}

func NewCriticalMessage(data interface{}, opts ...LoggingMessageOption) (*LoggingMessageNotification, error) {
    return NewLoggingMessage(LogLevelCritical, data, opts...)
}

/* Usage Example:
func ExampleLogging() {
    // Set logging level
    setLevelReq, err := NewSetLevelRequest(LogLevelInfo)
    if err != nil {
        log.Fatal(err)
    }

    // Create simple info message
    infoMsg, err := NewInfoMessage(
        "Server started successfully",
        WithLogger("system"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create structured debug message
    debugData := struct {
        Operation string `json:"operation"`
        Duration  int    `json:"duration"`
        Success   bool   `json:"success"`
    }{
        Operation: "database_query",
        Duration:  150,
        Success:   true,
    }

    debugMsg, err := NewDebugMessage(
        debugData,
        WithLogger("database"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create error message
    errorData := struct {
        Code    int    `json:"code"`
        Message string `json:"message"`
        Details string `json:"details"`
    }{
        Code:    500,
        Message: "Database connection failed",
        Details: "Connection timeout after 30s",
    }

    errorMsg, err := NewErrorMessage(
        errorData,
        WithLogger("database"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example JSON output for error message:
    // {
    //     "method": "notifications/message",
    //     "params": {
    //         "level": "error",
    //         "logger": "database",
    //         "data": {
    //             "code": 500,
    //             "message": "Database connection failed",
    //             "details": "Connection timeout after 30s"
    //         }
    //     }
    // }
}

// Example of using different log levels with structured data
func ExampleStructuredLogging() {
    // Define some common structured data types
    type MetricData struct {
        Name      string            `json:"name"`
        Value     float64           `json:"value"`
        Tags      map[string]string `json:"tags"`
        Timestamp int64            `json:"timestamp"`
    }

    type ErrorData struct {
        Error     string    `json:"error"`
        Component string    `json:"component"`
        Stack     []string  `json:"stack,omitempty"`
    }

    // Create metric log
    metricMsg, _ := NewInfoMessage(
        MetricData{
            Name:      "request_duration",
            Value:     123.45,
            Tags:      map[string]string{"endpoint": "/api/v1/users"},
            Timestamp: time.Now().Unix(),
        },
        WithLogger("metrics"),
    )

    // Create error log
    errorMsg, _ := NewErrorMessage(
        ErrorData{
            Error:     "invalid configuration",
            Component: "config-loader",
            Stack:     []string{"main.go:123", "config.go:45"},
        },
        WithLogger("configuration"),
    )
}
*/
