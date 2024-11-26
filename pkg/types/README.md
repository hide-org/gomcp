# Model Context Protocol Types (gomcp/pkg/types)

This package provides Go types and utilities for working with the Model Context Protocol (MCP). It offers type-safe, idiomatic Go implementations of MCP messages, with proper JSON marshaling/unmarshaling and validation.

## Usage Examples

### Error Handling

MCP provides structured error handling for different scenarios. Here are some common use cases:

```go
// Handling validation errors
validationErr := NewValidationError([]ValidationFailure{
    {Field: "maxTokens", Error: "must be positive"},
    {Field: "temperature", Error: "must be between 0 and 1"},
})

response := Response{
    JSONRPC: "2.0",
    ID:      1,
    Error:   validationErr,
}

// Will produce JSON:
// {
//     "jsonrpc": "2.0",
//     "id": 1,
//     "error": {
//         "code": -32602,
//         "message": "Invalid parameters",
//         "data": {
//             "validation": [
//                 {"field": "maxTokens", "error": "must be positive"},
//                 {"field": "temperature", "error": "must be between 0 and 1"}
//             ]
//         }
//     }
// }

// Handling tool execution errors
toolErr := NewToolExecutionError(
    "searchCode",
    "timeout",
    "Operation timed out after 30s",
)

// Will produce JSON:
// {
//     "code": -32603,
//     "message": "Tool execution failed",
//     "data": {
//         "toolName": "searchCode",
//         "errorType": "timeout",
//         "details": "Operation timed out after 30s"
//     }
// }

// Parsing errors from JSON
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
    log.Fatal(unmarshalErr)
}

// Access structured error data
if validationErr, ok := err.Data.(ValidationError); ok {
    for _, failure := range validationErr.Validation {
        fmt.Printf("Field %s: %s\n", failure.Field, failure.Error)
    }
}
```

### Content Handling

MCP supports different types of content (text, images, resources). Here's how to work with them:

```go
// Creating text content
textContent := NewTextContent("Hello, world!")

// Creating image content
imageContent := NewImageContent(base64Data, "image/png")

// Adding annotations
content := NewTextContent("Important message")
content.Annotations = &Annotations{
    Audience: []Role{RoleAssistant},
    Priority: ptr(0.8), // helper function for float64 pointer
}
```

(More examples will be added as we implement other features)

## Package Structure

```
types/
├── README.md      - This documentation
├── consts.go      - Protocol constants
├── errors.go      - Error types and handling
├── content.go     - Content type definitions
├── message.go     - Message type definitions
├── tool.go        - Tool-related types
├── resource.go    - Resource management types
├── prompt.go      - Prompt-related types
├── capabilities.go - Capability definitions
└── initialize.go  - Initialization types
```

## Error Codes

The package defines standard JSON-RPC error codes:

- `-32700`: Parse error
- `-32600`: Invalid request
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error

## Contributing

When adding new types or modifying existing ones:

1. Ensure proper JSON marshaling/unmarshaling
2. Include validation where appropriate
3. Add tests for new functionality
4. Update documentation with examples
5. Follow Go best practices and MCP specification

## References

- [MCP Specification](https://github.com/modelcontextprotocol/specification)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
