package types

import (
	"fmt"
)

// ProgressToken is used to associate progress notifications with requests
type ProgressToken int64

// ProgressNotificationOption configures ProgressNotification
type ProgressNotificationOption func(*ProgressNotification) error

// ProgressNotification represents a progress update for a long-running operation
type ProgressNotification struct {
	Method string         `json:"method"`
	Params ProgressParams `json:"params"`
}

type ProgressParams struct {
	ProgressToken ProgressToken `json:"progressToken"`
	Progress      float64       `json:"progress"`
	Total         *float64      `json:"total,omitempty"`
}

func NewProgressNotification(token ProgressToken, progress float64, opts ...ProgressNotificationOption) (*ProgressNotification, error) {
	if progress < 0 {
		return nil, fmt.Errorf("progress cannot be negative")
	}

	notification := &ProgressNotification{
		Method: "notifications/progress",
		Params: ProgressParams{
			ProgressToken: token,
			Progress:      progress,
		},
	}

	for _, opt := range opts {
		if err := opt(notification); err != nil {
			return nil, fmt.Errorf("applying progress notification option: %w", err)
		}
	}

	return notification, nil
}

// Progress notification options

func WithProgressTotal(total float64) ProgressNotificationOption {
	return func(n *ProgressNotification) error {
		if total <= 0 {
			return fmt.Errorf("total must be positive")
		}
		if n.Params.Progress > total {
			return fmt.Errorf("progress (%f) cannot exceed total (%f)", n.Params.Progress, total)
		}
		n.Params.Total = &total
		return nil
	}
}

// RequestProgressOption configures progress tracking for requests
type RequestProgressOption func(*RequestProgressMeta) error

// RequestProgressMeta represents progress-related metadata for requests
type RequestProgressMeta struct {
	ProgressToken ProgressToken `json:"progressToken,omitempty"`
}

// Helper functions for creating progress notifications at different completion stages

func NewProgressStart(token ProgressToken) (*ProgressNotification, error) {
	return NewProgressNotification(token, 0)
}

func NewProgressComplete(token ProgressToken) (*ProgressNotification, error) {
	return NewProgressNotification(token, 100)
}

func NewProgressPercentage(token ProgressToken, percentage float64) (*ProgressNotification, error) {
	if percentage < 0 || percentage > 100 {
		return nil, fmt.Errorf("percentage must be between 0 and 100")
	}
	return NewProgressNotification(token, percentage)
}

// Progress tracking for specific operations

func NewProgressWithItems(token ProgressToken, completed, total int) (*ProgressNotification, error) {
	if completed < 0 {
		return nil, fmt.Errorf("completed items cannot be negative")
	}
	if total <= 0 {
		return nil, fmt.Errorf("total items must be positive")
	}
	if completed > total {
		return nil, fmt.Errorf("completed items (%d) cannot exceed total items (%d)", completed, total)
	}

	percentage := (float64(completed) / float64(total)) * 100
	return NewProgressNotification(
		token,
		percentage,
		WithProgressTotal(float64(total)),
	)
}

/* Usage Example:
func ExampleProgress() {
    // Simple progress tracking
    token := ProgressToken(1)

    // Start operation
    startNotification, err := NewProgressStart(token)
    if err != nil {
        log.Fatal(err)
    }

    // Update progress at 50%
    halfwayNotification, err := NewProgressPercentage(token, 50)
    if err != nil {
        log.Fatal(err)
    }

    // Complete operation
    completeNotification, err := NewProgressComplete(token)
    if err != nil {
        log.Fatal(err)
    }

    // Progress with known total items
    itemProgress, err := NewProgressWithItems(
        token,
        75,  // completed items
        100, // total items
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example of processing a collection with progress updates
    func ProcessItemsWithProgress(items []string, token ProgressToken) error {
        total := len(items)
        for i, item := range items {
            // Process item...

            // Send progress notification
            notification, err := NewProgressWithItems(token, i+1, total)
            if err != nil {
                return fmt.Errorf("creating progress notification: %w", err)
            }

            // Send notification...
        }
        return nil
    }

    // Example of custom progress tracking
    type FileUploadProgress struct {
        BytesUploaded int64  `json:"bytesUploaded"`
        TotalBytes    int64  `json:"totalBytes"`
        Filename      string `json:"filename"`
    }

    func TrackFileUpload(token ProgressToken, uploaded, total int64, filename string) (*ProgressNotification, error) {
        if total <= 0 {
            return nil, fmt.Errorf("total bytes must be positive")
        }
        if uploaded < 0 || uploaded > total {
            return nil, fmt.Errorf("invalid bytes uploaded")
        }

        percentage := (float64(uploaded) / float64(total)) * 100

        return NewProgressNotification(
            token,
            percentage,
            WithTotal(float64(total)),
        )
    }
}

// Example of using progress with a request
func ExampleRequestWithProgress() {
    type LongRunningRequest struct {
        Method string      `json:"method"`
        Params interface{} `json:"params"`
        Meta   *RequestProgressMeta `json:"_meta,omitempty"`
    }

    // Create request with progress tracking
    req := LongRunningRequest{
        Method: "someOperation",
        Params: map[string]interface{}{
            "data": "example",
        },
        Meta: &RequestProgressMeta{
            ProgressToken: ProgressToken(123),
        },
    }

    // Track progress at various stages
    progress := []struct {
        Completed int
        Total     int
    }{
        {0, 4},   // 0%
        {1, 4},   // 25%
        {2, 4},   // 50%
        {3, 4},   // 75%
        {4, 4},   // 100%
    }

    for _, p := range progress {
        notification, err := NewProgressWithItems(
            req.Meta.ProgressToken,
            p.Completed,
            p.Total,
        )
        if err != nil {
            log.Fatal(err)
        }
        // Send notification...
    }
}
*/
