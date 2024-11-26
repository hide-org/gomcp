package types

import (
	"encoding/json"
	"fmt"
)

// ServerCapabilityOption is used to configure ServerCapabilities
type ServerCapabilityOption func(*ServerCapabilities) error

// ClientCapabilityOption is used to configure ClientCapabilities
type ClientCapabilityOption func(*ClientCapabilities) error

// ServerCapabilities represents capabilities that a server supports
type ServerCapabilities struct {
	Experimental map[string]json.RawMessage `json:"experimental,omitempty"`
	Logging      *LoggingCapability         `json:"logging,omitempty"`
	Prompts      *PromptsCapability         `json:"prompts,omitempty"`
	Resources    *ResourcesCapability       `json:"resources,omitempty"`
	Tools        *ToolsCapability           `json:"tools,omitempty"`
}

// ClientCapabilities represents capabilities that a client supports
type ClientCapabilities struct {
	Experimental map[string]json.RawMessage `json:"experimental,omitempty"`
	Roots        *RootsCapability           `json:"roots,omitempty"`
	Sampling     *SamplingCapability        `json:"sampling,omitempty"`
}

// Specific capability types

type LoggingCapability struct{}

type PromptsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
	Subscribe   *bool `json:"subscribe,omitempty"`
	ListChanged *bool `json:"listChanged,omitempty"`
}

type ToolsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type RootsCapability struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type SamplingCapability struct{}

// Server capabilities constructor and options

func NewServerCapabilities(opts ...ServerCapabilityOption) (*ServerCapabilities, error) {
	sc := &ServerCapabilities{
		Experimental: make(map[string]json.RawMessage),
	}

	for _, opt := range opts {
		if err := opt(sc); err != nil {
			return nil, fmt.Errorf("applying server capability option: %w", err)
		}
	}

	return sc, nil
}

// Server capability options

func WithServerLogging() ServerCapabilityOption {
	return func(sc *ServerCapabilities) error {
		sc.Logging = &LoggingCapability{}
		return nil
	}
}

func WithServerPrompts(listChanged bool) ServerCapabilityOption {
	return func(sc *ServerCapabilities) error {
		sc.Prompts = &PromptsCapability{
			ListChanged: &listChanged,
		}
		return nil
	}
}

func WithServerResources(subscribe, listChanged bool) ServerCapabilityOption {
	return func(sc *ServerCapabilities) error {
		sc.Resources = &ResourcesCapability{
			Subscribe:   &subscribe,
			ListChanged: &listChanged,
		}
		return nil
	}
}

func WithServerTools(listChanged bool) ServerCapabilityOption {
	return func(sc *ServerCapabilities) error {
		sc.Tools = &ToolsCapability{
			ListChanged: &listChanged,
		}
		return nil
	}
}

func WithServerExperimental(name string, data interface{}) ServerCapabilityOption {
	return func(sc *ServerCapabilities) error {
		rawData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshaling experimental data: %w", err)
		}
		sc.Experimental[name] = rawData
		return nil
	}
}

// Client capabilities constructor and options

func NewClientCapabilities(opts ...ClientCapabilityOption) (*ClientCapabilities, error) {
	cc := &ClientCapabilities{
		Experimental: make(map[string]json.RawMessage),
	}

	for _, opt := range opts {
		if err := opt(cc); err != nil {
			return nil, fmt.Errorf("applying client capability option: %w", err)
		}
	}

	return cc, nil
}

// Client capability options

func WithClientRoots(listChanged bool) ClientCapabilityOption {
	return func(cc *ClientCapabilities) error {
		cc.Roots = &RootsCapability{
			ListChanged: &listChanged,
		}
		return nil
	}
}

func WithClientSampling() ClientCapabilityOption {
	return func(cc *ClientCapabilities) error {
		cc.Sampling = &SamplingCapability{}
		return nil
	}
}

func WithClientExperimental(name string, data interface{}) ClientCapabilityOption {
	return func(cc *ClientCapabilities) error {
		rawData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshaling experimental data: %w", err)
		}
		cc.Experimental[name] = rawData
		return nil
	}
}

/* Usage Example:
func ExampleCapabilities() {
    // Create server capabilities
    serverCaps, err := NewServerCapabilities(
        WithServerLogging(),
        WithServerPrompts(true),  // with list changes enabled
        WithServerResources(true, true),  // with subscribe and list changes
        WithServerTools(true),  // with list changes
        WithServerExperimental("streaming", map[string]interface{}{
            "feature": "customStreaming",
            "config": map[string]interface{}{
                "maxChunkSize": 1024,
            },
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create client capabilities
    clientCaps, err := NewClientCapabilities(
        WithClientRoots(true),  // with list changes
        WithClientSampling(),
        WithClientExperimental("customFeature", map[string]interface{}{
            "enabled": true,
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
}
*/
