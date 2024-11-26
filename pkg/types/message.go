package types

import (
	"fmt"
)

// Message represents a communication element in the protocol
type Message struct {
	Role    Role    `json:"role"`
	Content Content `json:"content"`
}

func (m *Message) Validate() error {
	switch m.Role {
	case RoleUser, RoleAssistant:
		// valid roles
	default:
		return fmt.Errorf("invalid role: %s", m.Role)
	}

	// The Content type already handles its own validation through JSON marshaling
	return nil
}

// SamplingMessage represents a message used in sampling operations
type SamplingMessage struct {
	Role    Role    `json:"role"`
	Content Content `json:"content"`
}

// ModelPreferences represents preferences for model selection
type ModelPreferences struct {
	// Optional hints for model selection
	Hints []ModelHint `json:"hints,omitempty"`

	// Priority values (0-1 range)
	CostPriority         *float64 `json:"costPriority,omitempty"`
	SpeedPriority        *float64 `json:"speedPriority,omitempty"`
	IntelligencePriority *float64 `json:"intelligencePriority,omitempty"`
}

func (mp *ModelPreferences) Validate() error {
	if mp == nil {
		return nil
	}

	priorities := map[string]*float64{
		"costPriority":         mp.CostPriority,
		"speedPriority":        mp.SpeedPriority,
		"intelligencePriority": mp.IntelligencePriority,
	}

	for name, priority := range priorities {
		if priority != nil {
			if *priority < 0 || *priority > 1 {
				return fmt.Errorf("%s must be between 0 and 1, got %f", name, *priority)
			}
		}
	}

	return nil
}

type ModelHint struct {
	// Name suggests a model name or family
	// Examples: "claude-3-5-sonnet", "claude", "gemini-1.5-flash"
	Name *string `json:"name,omitempty"`
}

// CreateMessageParams represents parameters for creating a message
type CreateMessageParams struct {
	Messages         []SamplingMessage `json:"messages"`
	ModelPreferences *ModelPreferences `json:"modelPreferences,omitempty"`
	SystemPrompt     *string           `json:"systemPrompt,omitempty"`
	IncludeContext   *IncludeContext   `json:"includeContext,omitempty"`
	Temperature      *float64          `json:"temperature,omitempty"`
	MaxTokens        int               `json:"maxTokens"`
	StopSequences    []string          `json:"stopSequences,omitempty"`
	Metadata         map[string]any    `json:"metadata,omitempty"`
}

type IncludeContext string

const (
	IncludeContextNone       IncludeContext = "none"
	IncludeContextThisServer IncludeContext = "thisServer"
	IncludeContextAllServers IncludeContext = "allServers"
)

func (p *CreateMessageParams) Validate() error {
	if len(p.Messages) == 0 {
		return fmt.Errorf("messages cannot be empty")
	}

	if p.MaxTokens <= 0 {
		return fmt.Errorf("maxTokens must be positive")
	}

	if p.Temperature != nil {
		if *p.Temperature < 0 || *p.Temperature > 1 {
			return fmt.Errorf("temperature must be between 0 and 1")
		}
	}

	if p.IncludeContext != nil {
		switch *p.IncludeContext {
		case IncludeContextNone, IncludeContextThisServer, IncludeContextAllServers:
			// valid values
		default:
			return fmt.Errorf("invalid includeContext value: %s", *p.IncludeContext)
		}
	}

	if p.ModelPreferences != nil {
		if err := p.ModelPreferences.Validate(); err != nil {
			return fmt.Errorf("invalid model preferences: %w", err)
		}
	}

	return nil
}

/* Usage Example:
func ExampleMessage() {
    // Create a simple text message
    msg := Message{
        Role: RoleAssistant,
        Content: *NewTextContent("Hello, world!", nil),
    }

    // Marshal to JSON
    data, err := json.Marshal(msg)
    if err != nil {
        log.Fatal(err)
    }

    // Will produce:
    // {
    //     "role": "assistant",
    //     "content": {
    //         "type": "text",
    //         "text": "Hello, world!"
    //     }
    // }

    // Create a message with model preferences
    createParams := CreateMessageParams{
        Messages: []SamplingMessage{{
            Role: RoleUser,
            Content: *NewTextContent("What is Go?", nil),
        }},
        ModelPreferences: &ModelPreferences{
            Hints: []ModelHint{{Name: ptr("claude-3")}},
            IntelligencePriority: ptr(0.8),
            SpeedPriority: ptr(0.4),
        },
        MaxTokens: 1000,
        Temperature: ptr(0.7),
        IncludeContext: ptr(IncludeContextThisServer),
    }

    if err := createParams.Validate(); err != nil {
        log.Fatal(err)
    }
}

// Helper function for string pointers
func ptr(s string) *string {
    return &s
}

// Helper function for float64 pointers
func ptrf(f float64) *float64 {
    return &f
}
*/
