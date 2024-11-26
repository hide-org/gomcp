package types

// JSONSchemaType represents valid JSON Schema types
type JSONSchemaType string

const (
    TypeObject  JSONSchemaType = "object"
    TypeArray   JSONSchemaType = "array"
    TypeString  JSONSchemaType = "string"
    TypeNumber  JSONSchemaType = "number"
    TypeInteger JSONSchemaType = "integer"
    TypeBoolean JSONSchemaType = "boolean"
    TypeNull    JSONSchemaType = "null"
)

// SchemaEnum represents possible enum values in JSON Schema
type SchemaEnum []interface{}

// JSONSchema represents a JSON Schema object for tool input validation
type JSONSchema struct {
    Type       JSONSchemaType          `json:"type"`
    Properties map[string]JSONSchema   `json:"properties,omitempty"`
    Required   []string               `json:"required,omitempty"`
    Items      *JSONSchema            `json:"items,omitempty"`
    Enum       SchemaEnum             `json:"enum,omitempty"`
    // Additional common JSON Schema fields
    MinLength  *int                   `json:"minLength,omitempty"`
    MaxLength  *int                   `json:"maxLength,omitempty"`
    Minimum    *float64               `json:"minimum,omitempty"`
    Maximum    *float64               `json:"maximum,omitempty"`
    Pattern    *string                `json:"pattern,omitempty"`
}

// Common schema constructors
func NewStringEnum(values ...string) JSONSchema {
    enum := make(SchemaEnum, len(values))
    for i, v := range values {
        enum[i] = v
    }
    return JSONSchema{
        Type: TypeString,
        Enum: enum,
    }
}

func NewNumberEnum(values ...float64) JSONSchema {
    enum := make(SchemaEnum, len(values))
    for i, v := range values {
        enum[i] = v
    }
    return JSONSchema{
        Type: TypeNumber,
        Enum: enum,
    }
}

func NewIntegerEnum(values ...int) JSONSchema {
    enum := make(SchemaEnum, len(values))
    for i, v := range values {
        enum[i] = v
    }
    return JSONSchema{
        Type: TypeInteger,
        Enum: enum,
    }
}

// Predefined schemas
var (
    StringSchema = JSONSchema{Type: TypeString}
    NumberSchema = JSONSchema{Type: TypeNumber}
    IntegerSchema = JSONSchema{Type: TypeInteger}
    BooleanSchema = JSONSchema{Type: TypeBoolean}
)

// Schema constructors with constraints
func StringSchemaWithConstraints(opts ...SchemaOption) JSONSchema {
    schema := StringSchema
    for _, opt := range opts {
        opt(&schema)
    }
    return schema
}

// SchemaOption configures a JSONSchema
type SchemaOption func(*JSONSchema)

func WithMinLength(min int) SchemaOption {
    return func(s *JSONSchema) {
        s.MinLength = &min
    }
}

func WithMaxLength(max int) SchemaOption {
    return func(s *JSONSchema) {
        s.MaxLength = &max
    }
}

func WithPattern(pattern string) SchemaOption {
    return func(s *JSONSchema) {
        s.Pattern = &pattern
    }
}

func WithNumberRange(min, max float64) SchemaOption {
    return func(s *JSONSchema) {
        s.Minimum = &min
        s.Maximum = &max
    }
}

// Array and Object schema constructors
func ArraySchema(items JSONSchema) JSONSchema {
    return JSONSchema{
        Type:  TypeArray,
        Items: &items,
    }
}

func ObjectSchema(properties map[string]JSONSchema) JSONSchema {
    return JSONSchema{
        Type:       TypeObject,
        Properties: properties,
    }
}

// Rest of the tool.go implementation remains the same, but now we can use these more type-safe schemas:

/* Usage Example:
func ExampleToolWithSchema() {
    deployTool, err := NewTool("deployService",
        WithToolDescription("Deploy a service to the cloud"),
        WithToolProperty("name", StringSchemaWithConstraints(
            WithMinLength(3),
            WithMaxLength(63),
            WithPattern("^[a-z0-9-]+$"),
        )),
        WithToolProperty("environment", NewStringEnum("dev", "staging", "prod")),
        WithToolProperty("replicas", JSONSchema{
            Type:    TypeInteger,
            Minimum: ptr(1),
            Maximum: ptr(10),
        }),
        WithToolProperty("cpu", NewNumberEnum(0.5, 1.0, 2.0)),
        WithToolRequired("name", "environment"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example of complex nested schema
    serviceSchema := ObjectSchema(map[string]JSONSchema{
        "name": StringSchemaWithConstraints(
            WithMinLength(1),
            WithMaxLength(100),
        ),
        "ports": ArraySchema(ObjectSchema(map[string]JSONSchema{
            "number": IntegerSchema,
            "protocol": NewStringEnum("TCP", "UDP"),
        })),
        "environment": NewStringEnum("dev", "staging", "prod"),
    })

    // Helper functions for creating pointers
    func ptr[T any](v T) *T {
        return &v
    }
}
*/
