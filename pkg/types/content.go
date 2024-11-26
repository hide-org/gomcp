package types

import (
	"encoding/json"
	"fmt"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Annotations provides optional metadata for content
type Annotations struct {
	// TODO: check how this relates to Message.Role
	Audience []Role   `json:"audience,omitempty"`
	Priority *float64 `json:"priority,omitempty"`
}

func (a *Annotations) Validate() error {
	if a == nil {
		return nil
	}

	if a.Priority != nil {
		if *a.Priority < 0 || *a.Priority > 1 {
			return fmt.Errorf("priority must be between 0 and 1, got %f", *a.Priority)
		}
	}

	if len(a.Audience) > 0 {
		for _, role := range a.Audience {
			switch role {
			case RoleUser, RoleAssistant:
				// valid roles
			default:
				return fmt.Errorf("invalid role in audience: %s", role)
			}
		}
	}

	return nil
}

// ContentType for type discrimination
type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeImage    ContentType = "image"
	ContentTypeResource ContentType = "resource"
)

// Content represents the interface that all content types must implement
type Content struct {
	// Type field for JSON discrimination
	Type ContentType `json:"type"`
	// Only one of these will be non-nil
	TextContent     *TextContent     `json:"text,omitempty"`
	ImageContent    *ImageContent    `json:"image,omitempty"`
	ResourceContent *ResourceContent `json:"resource,omitempty"`
}

type TextContent struct {
	Text        string       `json:"text"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

type ImageContent struct {
	Data        string       `json:"data"` // base64 encoded
	MimeType    string       `json:"mimeType"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// Custom JSON marshaling/unmarshaling
func (c *Content) UnmarshalJSON(data []byte) error {
	// First unmarshal the discriminator
	type typeOnly struct {
		Type ContentType `json:"type"`
	}
	var t typeOnly
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	// Then unmarshal the full content based on type
	switch t.Type {
	case ContentTypeText:
		var text TextContent
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}
		c.Type = ContentTypeText
		c.TextContent = &text
	case ContentTypeImage:
		var img ImageContent
		if err := json.Unmarshal(data, &img); err != nil {
			return err
		}
		c.Type = ContentTypeImage
		c.ImageContent = &img
	case ContentTypeResource:
		var res ResourceContent
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		c.Type = ContentTypeResource
		c.ResourceContent = &res
	default:
		return fmt.Errorf("unknown content type: %s", t.Type)
	}

	return nil
}

func (c Content) MarshalJSON() ([]byte, error) {
	switch c.Type {
	case ContentTypeText:
		if c.TextContent == nil {
			return nil, fmt.Errorf("text content is nil")
		}
		return json.Marshal(struct {
			Type ContentType `json:"type"`
			*TextContent
		}{
			Type:        ContentTypeText,
			TextContent: c.TextContent,
		})
	case ContentTypeImage:
		if c.ImageContent == nil {
			return nil, fmt.Errorf("image content is nil")
		}
		return json.Marshal(struct {
			Type ContentType `json:"type"`
			*ImageContent
		}{
			Type:         ContentTypeImage,
			ImageContent: c.ImageContent,
		})
	case ContentTypeResource:
		if c.ResourceContent == nil {
			return nil, fmt.Errorf("resource content is nil")
		}
		return json.Marshal(struct {
			Type ContentType `json:"type"`
			*ResourceContent
		}{
			Type:            ContentTypeResource,
			ResourceContent: c.ResourceContent,
		})
	default:
		return nil, fmt.Errorf("unknown content type: %s", c.Type)
	}
}

// Helper constructors
func NewTextContent(text string, annotations *Annotations) *Content {
	return &Content{
		Type: ContentTypeText,
		TextContent: &TextContent{
			Text:        text,
			Annotations: annotations,
		},
	}
}

func NewImageContent(data, mimeType string, annotations *Annotations) *Content {
	return &Content{
		Type: ContentTypeImage,
		ImageContent: &ImageContent{
			Data:        data,
			MimeType:    mimeType,
			Annotations: annotations,
		},
	}
}

/* Usage Example:
message := Content{
    Type: ContentTypeText,
    TextContent: &TextContent{
        Text: "Hello, world!",
        Annotations: &Annotations{
            Audience: []Role{RoleAssistant},
            Priority: ptr(0.8),
        },
    },
}

// Will produce JSON like:
{
    "type": "text",
    "text": "Hello, world!",
    "annotations": {
        "audience": ["assistant"],
        "priority": 0.8
    }
}

// Unmarshaling example:
var content Content
err := json.Unmarshal(data, &content)
if err != nil {
    // handle error
}

switch content.Type {
case ContentTypeText:
    fmt.Println(content.TextContent.Text)
case ContentTypeImage:
    fmt.Println(content.ImageContent.MimeType)
case ContentTypeResource:
    fmt.Println(content.ResourceContent.URI)
}
*/
