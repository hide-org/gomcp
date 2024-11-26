package types

import (
	"fmt"
)

// InitializeRequestOption configures InitializeRequest
type InitializeRequestOption func(*InitializeRequest) error

// InitializeRequest represents the initial request from client to server
type InitializeRequest struct {
    Method          string             `json:"method"`
    Params          InitializeParams   `json:"params"`
}

type InitializeParams struct {
    ProtocolVersion string             `json:"protocolVersion"`
    Capabilities    ClientCapabilities `json:"capabilities"`
    ClientInfo      Implementation     `json:"clientInfo"`
}

func NewInitializeRequest(clientInfo Implementation, opts ...InitializeRequestOption) (*InitializeRequest, error) {
    req := &InitializeRequest{
        Method: "initialize",
        Params: InitializeParams{
            ProtocolVersion: LatestProtocolVersion,
            ClientInfo:      clientInfo,
        },
    }

    for _, opt := range opts {
        if err := opt(req); err != nil {
            return nil, fmt.Errorf("applying initialize request option: %w", err)
        }
    }

    return req, nil
}

// InitializeRequest options

func WithProtocolVersion(version string) InitializeRequestOption {
    return func(r *InitializeRequest) error {
        r.Params.ProtocolVersion = version
        return nil
    }
}

func WithClientCapabilities(opts ...ClientCapabilityOption) InitializeRequestOption {
    return func(r *InitializeRequest) error {
        caps, err := NewClientCapabilities(opts...)
        if err != nil {
            return fmt.Errorf("creating client capabilities: %w", err)
        }
        r.Params.Capabilities = *caps
        return nil
    }
}

// InitializeResultOption configures InitializeResult
type InitializeResultOption func(*InitializeResult) error

// InitializeResult represents the server's response to initialization
type InitializeResult struct {
    ProtocolVersion string             `json:"protocolVersion"`
    ServerInfo      Implementation     `json:"serverInfo"`
    Capabilities    ServerCapabilities `json:"capabilities"`
    Instructions    *string            `json:"instructions,omitempty"`
}

func NewInitializeResult(serverInfo Implementation, opts ...InitializeResultOption) (*InitializeResult, error) {
    result := &InitializeResult{
        ProtocolVersion: LatestProtocolVersion,
        ServerInfo:      serverInfo,
    }

    for _, opt := range opts {
        if err := opt(result); err != nil {
            return nil, fmt.Errorf("applying initialize result option: %w", err)
        }
    }

    return result, nil
}

// InitializeResult options

func WithServerCapabilities(opts ...ServerCapabilityOption) InitializeResultOption {
    return func(r *InitializeResult) error {
        caps, err := NewServerCapabilities(opts...)
        if err != nil {
            return fmt.Errorf("creating server capabilities: %w", err)
        }
        r.Capabilities = *caps
        return nil
    }
}

func WithInstructions(instructions string) InitializeResultOption {
    return func(r *InitializeResult) error {
        r.Instructions = &instructions
        return nil
    }
}

// Implementation represents an MCP implementation
type Implementation struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

func NewImplementation(name, version string) (*Implementation, error) {
    if name == "" {
        return nil, fmt.Errorf("implementation name cannot be empty")
    }
    if version == "" {
        return nil, fmt.Errorf("implementation version cannot be empty")
    }

    return &Implementation{
        Name:    name,
        Version: version,
    }, nil
}

// InitializedNotification represents the notification sent after initialization
type InitializedNotification struct {
    Method string                  `json:"method"`
    Params *InitializedParams     `json:"params,omitempty"`
}

type InitializedParams struct {
    Meta map[string]interface{} `json:"_meta,omitempty"`
}

func NewInitializedNotification(meta map[string]interface{}) *InitializedNotification {
    params := &InitializedParams{}
    if len(meta) > 0 {
        params.Meta = meta
    }

    return &InitializedNotification{
        Method: "notifications/initialized",
        Params: params,
    }
}

/* Usage Example:
func ExampleInitialize() {
    // Create client implementation info
    clientInfo, err := NewImplementation(
        "example-client",
        "1.0.0",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create initialize request
    request, err := NewInitializeRequest(
        *clientInfo,
        WithClientCapabilities(
            WithClientRoots(true),
            WithClientSampling(),
            WithClientExperimental("customFeature", map[string]interface{}{
                "enabled": true,
            }),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create server implementation info
    serverInfo, err := NewImplementation(
        "example-server",
        "2.0.0",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create initialize result
    result, err := NewInitializeResult(
        *serverInfo,
        WithServerCapabilities(
            WithServerLogging(),
            WithServerPrompts(true),
            WithServerResources(true, true),
            WithServerTools(true),
        ),
        WithInstructions(`
            This server provides code analysis and generation capabilities.
            Available tools:
            - searchCode: Search for code in the repository
            - formatCode: Format source code
            - analyzeCode: Analyze code for issues
        `),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create initialized notification
    notification := NewInitializedNotification(map[string]interface{}{
        "clientId": "client-123",
        "sessionStarted": time.Now().Unix(),
    })

    // Example JSON output for request:
    // {
    //     "method": "initialize",
    //     "params": {
    //         "protocolVersion": "2024-11-05",
    //         "clientInfo": {
    //             "name": "example-client",
    //             "version": "1.0.0"
    //         },
    //         "capabilities": {
    //             "roots": {
    //                 "listChanged": true
    //             },
    //             "sampling": {},
    //             "experimental": {
    //                 "customFeature": {
    //                     "enabled": true
    //                 }
    //             }
    //         }
    //     }
    // }
}
*/
