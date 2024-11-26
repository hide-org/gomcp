package types

import (
	"fmt"
)

// CompleteRequestOption configures CompleteRequest
type CompleteRequestOption func(*CompleteRequest) error

// CompleteRequest represents a request for completion options
type CompleteRequest struct {
	Method string         `json:"method"`
	Params CompleteParams `json:"params"`
}

type CompleteParams struct {
	Ref      Reference     `json:"ref"`
	Argument CompletionArg `json:"argument"`
}

// Reference represents either a prompt or resource reference
type Reference struct {
	Type string `json:"type"`
	// Only one of these should be set based on Type
	Name *string `json:"name,omitempty"` // for prompt references
	URI  *string `json:"uri,omitempty"`  // for resource references
}

func NewPromptReference(name string) Reference {
	return Reference{
		Type: "ref/prompt",
		Name: &name,
	}
}

func NewResourceReference(uri string) Reference {
	return Reference{
		Type: "ref/resource",
		URI:  &uri,
	}
}

type CompletionArg struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewCompleteRequest(ref Reference, argName, argValue string, opts ...CompleteRequestOption) (*CompleteRequest, error) {
	if err := validateReference(ref); err != nil {
		return nil, fmt.Errorf("invalid reference: %w", err)
	}

	if argName == "" {
		return nil, fmt.Errorf("argument name cannot be empty")
	}

	req := &CompleteRequest{
		Method: "completion/complete",
		Params: CompleteParams{
			Ref: ref,
			Argument: CompletionArg{
				Name:  argName,
				Value: argValue,
			},
		},
	}

	for _, opt := range opts {
		if err := opt(req); err != nil {
			return nil, fmt.Errorf("applying complete request option: %w", err)
		}
	}

	return req, nil
}

func validateReference(ref Reference) error {
	switch ref.Type {
	case "ref/prompt":
		if ref.Name == nil || *ref.Name == "" {
			return fmt.Errorf("prompt reference requires name")
		}
		if ref.URI != nil {
			return fmt.Errorf("prompt reference should not have URI")
		}
	case "ref/resource":
		if ref.URI == nil || *ref.URI == "" {
			return fmt.Errorf("resource reference requires URI")
		}
		if ref.Name != nil {
			return fmt.Errorf("resource reference should not have name")
		}
	default:
		return fmt.Errorf("invalid reference type: %s", ref.Type)
	}
	return nil
}

// CompleteResultOption configures CompleteResult
type CompleteResultOption func(*CompleteResult) error

// CompleteResult represents the response to a completion request
type CompleteResult struct {
	Completion CompletionInfo `json:"completion"`
}

type CompletionInfo struct {
	Values  []string `json:"values"`
	Total   *int     `json:"total,omitempty"`
	HasMore *bool    `json:"hasMore,omitempty"`
}

func NewCompleteResult(values []string, opts ...CompleteResultOption) (*CompleteResult, error) {
	if len(values) > 100 {
		return nil, fmt.Errorf("completion values cannot exceed 100 items")
	}

	result := &CompleteResult{
		Completion: CompletionInfo{
			Values: values,
		},
	}

	for _, opt := range opts {
		if err := opt(result); err != nil {
			return nil, fmt.Errorf("applying complete result option: %w", err)
		}
	}

	return result, nil
}

// CompleteResult options

func WithResultTotal(total int) CompleteResultOption {
	return func(r *CompleteResult) error {
		if total < len(r.Completion.Values) {
			return fmt.Errorf("total (%d) cannot be less than number of values (%d)", total, len(r.Completion.Values))
		}
		r.Completion.Total = &total
		return nil
	}
}

func WithHasMore(hasMore bool) CompleteResultOption {
	return func(r *CompleteResult) error {
		r.Completion.HasMore = &hasMore
		return nil
	}
}

/* Usage Example:
func ExampleCompletion() {
    // Create a completion request for a prompt argument
    promptRef := NewPromptReference("generateCode")
    promptRequest, err := NewCompleteRequest(
        promptRef,
        "language",  // argument name
        "py",       // partial value for completion
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example completion result for programming languages
    languageResult, err := NewCompleteResult(
        []string{
            "python",
            "python3",
            "pypy",
        },
        WithTotal(3),
        WithHasMore(false),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create a completion request for a resource
    resourceRef := NewResourceReference("file:///configs/{env}")
    resourceRequest, err := NewCompleteRequest(
        resourceRef,
        "env",      // argument name
        "pr",       // partial value for completion
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example completion result for environments
    envResult, err := NewCompleteResult(
        []string{
            "prod",
            "pre-prod",
            "preview",
        },
        WithTotal(5),    // indicates there are more matches
        WithHasMore(true), // indicates more results available
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example of paginated completions
    func GetCompletions(ref Reference, argName, value string, cursor *string) (*CompleteResult, error) {
        // Simulated completion values
        allValues := []string{
            "value1", "value2", "value3", "value4", "value5",
            "value6", "value7", "value8", "value9", "value10",
        }

        // Simulate pagination
        pageSize := 3
        start := 0
        if cursor != nil {
            // Parse cursor to get start index
            // In real implementation, use proper cursor handling
            start = 3
        }

        end := start + pageSize
        if end > len(allValues) {
            end = len(allValues)
        }

        values := allValues[start:end]
        hasMore := end < len(allValues)

        return NewCompleteResult(
            values,
            WithTotal(len(allValues)),
            WithHasMore(hasMore),
        )
    }
}

// Example of structured completions
func ExampleStructuredCompletions() {
    // Configuration completion
    configRef := NewResourceReference("file:///app/config.yaml")
    configRequest, _ := NewCompleteRequest(
        configRef,
        "setting",
        "log.",
    )

    configResult, _ := NewCompleteResult(
        []string{
            "log.level",
            "log.format",
            "log.output",
        },
    )

    // API endpoint completion
    apiRef := NewPromptReference("apiEndpoint")
    apiRequest, _ := NewCompleteRequest(
        apiRef,
        "path",
        "/api/v1/us",
    )

    apiResult, _ := NewCompleteResult(
        []string{
            "/api/v1/users",
            "/api/v1/users/{id}",
            "/api/v1/user-groups",
        },
    )
}
*/
