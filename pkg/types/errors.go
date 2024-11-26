package types

import (
	"encoding/json"
	"fmt"
)

const (
	ErrParse          = -32700
	ErrInvalidRequest = -32600
	ErrMethodNotFound = -32601
	ErrInvalidParams  = -32602
	ErrInternal       = -32603
)

// ErrorData represents different types of error details
type ErrorData interface {
	isErrorData()
	ErrorType() string // helps with unmarshaling
}

type ValidationFailure struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationError struct {
	Validation []ValidationFailure `json:"validation"`
}

func (ValidationError) isErrorData()      {}
func (ValidationError) ErrorType() string { return "validation" }

type ToolExecutionError struct {
	ToolName string `json:"toolName"`
	ErrType  string `json:"errorType"`
	Details  string `json:"details"`
}

func (ToolExecutionError) isErrorData()      {}
func (ToolExecutionError) ErrorType() string { return "toolExecution" }

// ErrorInfo represents a JSON-RPC error
type ErrorInfo struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    ErrorData `json:"data,omitempty"`
}

// MarshalJSON implements custom marshaling for ErrorInfo
func (e ErrorInfo) MarshalJSON() ([]byte, error) {
	type Alias ErrorInfo
	aux := struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data,omitempty"`
	}{
		Code:    e.Code,
		Message: e.Message,
	}

	if e.Data != nil {
		data, err := json.Marshal(e.Data)
		if err != nil {
			return nil, fmt.Errorf("marshaling error data: %w", err)
		}
		aux.Data = data
	}

	return json.Marshal(aux)
}

// UnmarshalJSON implements custom unmarshaling for ErrorInfo
func (e *ErrorInfo) UnmarshalJSON(data []byte) error {
	type Alias ErrorInfo
	aux := struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data,omitempty"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Code = aux.Code
	e.Message = aux.Message

	if aux.Data != nil {
		// First unmarshal into a temporary structure to get the error type
		var temp struct {
			ErrorType string `json:"errorType"`
		}
		if err := json.Unmarshal(aux.Data, &temp); err != nil {
			return err
		}

		// Based on error code and/or type, unmarshal into appropriate structure
		switch e.Code {
		case ErrInvalidParams:
			var validationErr ValidationError
			if err := json.Unmarshal(aux.Data, &validationErr); err != nil {
				return err
			}
			e.Data = validationErr
		case ErrInternal:
			switch temp.ErrorType {
			case "toolExecution":
				var toolErr ToolExecutionError
				if err := json.Unmarshal(aux.Data, &toolErr); err != nil {
					return err
				}
				e.Data = toolErr
			default:
				return fmt.Errorf("unknown error type: %s", temp.ErrorType)
			}
		}
	}

	return nil
}

// Helper functions to create common errors
func NewValidationError(failures []ValidationFailure) *ErrorInfo {
	return &ErrorInfo{
		Code:    ErrInvalidParams,
		Message: "Invalid parameters",
		Data:    ValidationError{Validation: failures},
	}
}

func NewToolExecutionError(toolName, errorType, details string) *ErrorInfo {
	return &ErrorInfo{
		Code:    ErrInternal,
		Message: "Tool execution failed",
		Data: ToolExecutionError{
			ToolName: toolName,
			ErrType:  errorType,
			Details:  details,
		},
	}
}

// Usage examples:
/*
// Example 1: Validation error during parameter parsing
validationErr := NewValidationError([]ValidationFailure{
    {Field: "maxTokens", Error: "must be positive"},
    {Field: "temperature", Error: "must be between 0 and 1"},
})

// Example 2: Tool execution error
toolErr := NewToolExecutionError(
    "searchCode",
    "timeout",
    "Operation timed out after 30s",
)

// Example 3: Deserializing error from JSON
jsonData := `{
    "code": -32602,
    "message": "Invalid parameters",
    "data": {
        "validation": [
            {"field": "maxTokens", "error": "must be positive"}
        ]
    }
}`
var err ErrorInfo
if unmarshalErr := json.Unmarshal([]byte(jsonData), &err); unmarshalErr != nil {
    // handle error
}
*/
