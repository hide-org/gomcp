package types

import (
	"fmt"
)

// PromptOption configures a Prompt
type PromptOption func(*Prompt) error

// Prompt represents a prompt or prompt template
type Prompt struct {
    Name        string           `json:"name"`
    Description *string          `json:"description,omitempty"`
    Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument represents an argument that a prompt can accept
type PromptArgument struct {
    Name        string  `json:"name"`
    Description *string `json:"description,omitempty"`
    Required    *bool   `json:"required,omitempty"`
}

// PromptMessage represents a message returned as part of a prompt
type PromptMessage struct {
    Role    Role    `json:"role"`
    Content Content `json:"content"`
}

// NewPrompt creates a new Prompt with the given name and options
func NewPrompt(name string, opts ...PromptOption) (*Prompt, error) {
    if name == "" {
        return nil, fmt.Errorf("prompt name cannot be empty")
    }

    p := &Prompt{
        Name:      name,
        Arguments: make([]PromptArgument, 0),
    }

    for _, opt := range opts {
        if err := opt(p); err != nil {
            return nil, fmt.Errorf("applying prompt option: %w", err)
        }
    }

    return p, nil
}

// Prompt options

func WithPromptDescription(description string) PromptOption {
    return func(p *Prompt) error {
        p.Description = &description
        return nil
    }
}

func WithPromptArgument(name string, opts ...PromptArgumentOption) PromptOption {
    return func(p *Prompt) error {
        arg := PromptArgument{
            Name: name,
        }

        for _, opt := range opts {
            if err := opt(&arg); err != nil {
                return fmt.Errorf("applying argument option: %w", err)
            }
        }

        p.Arguments = append(p.Arguments, arg)
        return nil
    }
}

// PromptArgumentOption configures a PromptArgument
type PromptArgumentOption func(*PromptArgument) error

func WithArgumentDescription(description string) PromptArgumentOption {
    return func(a *PromptArgument) error {
        a.Description = &description
        return nil
    }
}

func WithArgumentRequired(required bool) PromptArgumentOption {
    return func(a *PromptArgument) error {
        a.Required = &required
        return nil
    }
}

// GetPromptRequest represents a request to get a prompt
type GetPromptRequest struct {
    Name      string            `json:"name"`
    Arguments map[string]string `json:"arguments,omitempty"`
}

// GetPromptResult represents the response to a get prompt request
type GetPromptResult struct {
    Description *string         `json:"description,omitempty"`
    Messages    []PromptMessage `json:"messages"`
}

// ListPromptsResult represents the response to a list prompts request
type ListPromptsResult struct {
    NextCursor *string  `json:"nextCursor,omitempty"`
    Prompts    []Prompt `json:"prompts"`
}

/* Usage Example:
func ExamplePrompt() {
    // Create a new prompt with arguments
    prompt, err := NewPrompt("generateCode",
        WithPromptDescription("Generates code based on description"),
        WithPromptArgument("language",
            WithArgumentDescription("Programming language to use"),
            WithArgumentRequired(true),
        ),
        WithPromptArgument("description",
            WithArgumentDescription("What the code should do"),
            WithArgumentRequired(true),
        ),
        WithPromptArgument("style",
            WithArgumentDescription("Coding style preferences"),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create a request to get the prompt with arguments
    request := GetPromptRequest{
        Name: "generateCode",
        Arguments: map[string]string{
            "language": "go",
            "description": "HTTP server with two endpoints",
            "style": "clean code",
        },
    }

    // Example of prompt messages in response
    result := GetPromptResult{
        Description: prompt.Description,
        Messages: []PromptMessage{
            {
                Role: RoleUser,
                Content: *NewTextContent(
                    "Please generate a Go HTTP server with two endpoints following clean code principles",
                    nil,
                ),
            },
        },
    }

    // Example of listing prompts
    listResult := ListPromptsResult{
        Prompts: []Prompt{*prompt},
    }
}
*/
