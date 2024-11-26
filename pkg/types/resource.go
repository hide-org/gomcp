package types

import (
	"fmt"
)

// ResourceOption configures a Resource
type ResourceOption func(*Resource) error

// Resource represents a known resource that the server can read
type Resource struct {
	URI         string       `json:"uri"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	MimeType    *string      `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

func NewResource(uri, name string, opts ...ResourceOption) (*Resource, error) {
	if uri == "" {
		return nil, fmt.Errorf("resource URI cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("resource name cannot be empty")
	}

	r := &Resource{
		URI:  uri,
		Name: name,
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, fmt.Errorf("applying resource option: %w", err)
		}
	}

	return r, nil
}

// Resource options

func WithResourceDescription(description string) ResourceOption {
	return func(r *Resource) error {
		r.Description = &description
		return nil
	}
}

func WithResourceMimeType(mimeType string) ResourceOption {
	return func(r *Resource) error {
		r.MimeType = &mimeType
		return nil
	}
}

func WithResourceAnnotations(annotations *Annotations) ResourceOption {
	return func(r *Resource) error {
		r.Annotations = annotations
		return nil
	}
}

// ResourceTemplate represents a template for resources
type ResourceTemplateOption func(*ResourceTemplate) error

type ResourceTemplate struct {
	Name        string       `json:"name"`
	URITemplate string       `json:"uriTemplate"`
	Description *string      `json:"description,omitempty"`
	MimeType    *string      `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

func NewResourceTemplate(name, uriTemplate string, opts ...ResourceTemplateOption) (*ResourceTemplate, error) {
	if name == "" {
		return nil, fmt.Errorf("template name cannot be empty")
	}
	if uriTemplate == "" {
		return nil, fmt.Errorf("URI template cannot be empty")
	}

	rt := &ResourceTemplate{
		Name:        name,
		URITemplate: uriTemplate,
	}

	for _, opt := range opts {
		if err := opt(rt); err != nil {
			return nil, fmt.Errorf("applying template option: %w", err)
		}
	}

	return rt, nil
}

// Resource template options

func WithTemplateDescription(description string) ResourceTemplateOption {
	return func(rt *ResourceTemplate) error {
		rt.Description = &description
		return nil
	}
}

func WithTemplateMimeType(mimeType string) ResourceTemplateOption {
	return func(rt *ResourceTemplate) error {
		rt.MimeType = &mimeType
		return nil
	}
}

func WithTemplateAnnotations(annotations *Annotations) ResourceTemplateOption {
	return func(rt *ResourceTemplate) error {
		rt.Annotations = annotations
		return nil
	}
}

// ResourceContentOption configures ResourceContent
type ResourceContentOption func(*ResourceContent) error

// ResourceContent represents the contents of a specific resource
type ResourceContent struct {
	URI         string       `json:"uri"`
	Text        *string      `json:"text,omitempty"`
	Blob        *string      `json:"blob,omitempty"` // base64 encoded
	MimeType    *string      `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

func NewResourceContent(uri string, opts ...ResourceContentOption) (*ResourceContent, error) {
	if uri == "" {
		return nil, fmt.Errorf("resource URI cannot be empty")
	}

	rc := &ResourceContent{
		URI: uri,
	}

	for _, opt := range opts {
		if err := opt(rc); err != nil {
			return nil, fmt.Errorf("applying content option: %w", err)
		}
	}

	// Validate that either Text or Blob is set, but not both
	if (rc.Text != nil && rc.Blob != nil) || (rc.Text == nil && rc.Blob == nil) {
		return nil, fmt.Errorf("exactly one of text or blob must be set")
	}

	return rc, nil
}

// Resource content options

func WithContentText(text string) ResourceContentOption {
	return func(rc *ResourceContent) error {
		if rc.Blob != nil {
			return fmt.Errorf("cannot set text when blob is already set")
		}
		rc.Text = &text
		return nil
	}
}

func WithContentBlob(blob string) ResourceContentOption {
	return func(rc *ResourceContent) error {
		if rc.Text != nil {
			return fmt.Errorf("cannot set blob when text is already set")
		}
		rc.Blob = &blob
		return nil
	}
}

func WithContentMimeType(mimeType string) ResourceContentOption {
	return func(rc *ResourceContent) error {
		rc.MimeType = &mimeType
		return nil
	}
}

func WithContentAnnotations(annotations *Annotations) ResourceContentOption {
	return func(rc *ResourceContent) error {
		rc.Annotations = annotations
		return nil
	}
}

// Request/Response types

type ReadResourceRequest struct {
	URI string `json:"uri"`
}

type ReadResourceResult struct {
	Contents []ResourceContent `json:"contents"`
}

type ListResourcesResult struct {
	NextCursor *string    `json:"nextCursor,omitempty"`
	Resources  []Resource `json:"resources"`
}

type ListResourceTemplatesResult struct {
	NextCursor        *string            `json:"nextCursor,omitempty"`
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates"`
}

/* Usage Example:
func ExampleResource() {
    // Create a new resource
    resource, err := NewResource(
        "file:///path/to/config.yaml",
        "Configuration",
        WithResourceDescription("Application configuration file"),
        WithResourceMimeType("application/yaml"),
        WithResourceAnnotations(&Annotations{
            Audience: []Role{RoleAssistant},
            Priority: ptr(0.8),
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create a resource template
    template, err := NewResourceTemplate(
        "ConfigFile",
        "file:///configs/{env}/config.yaml",
        WithTemplateDescription("Environment-specific configuration"),
        WithTemplateMimeType("application/yaml"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create resource content
    content, err := NewResourceContent(
        "file:///path/to/config.yaml",
        WithContentText("key: value\nother: data"),
        WithContentMimeType("application/yaml"),
        WithContentAnnotations(&Annotations{
            Audience: []Role{RoleAssistant},
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example of resource listing
    listResult := ListResourcesResult{
        Resources: []Resource{*resource},
    }

    // Example of reading resource
    readResult := ReadResourceResult{
        Contents: []ResourceContent{*content},
    }
}

// Helper function for float64 pointers
func ptr(f float64) *float64 {
    return &f
}
*/
