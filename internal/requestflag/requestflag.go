package requestflag

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

// Flag [T] is a generic flag base which can be used to implement the most
// common interfaces used by urfave/cli. Additionally, it allows specifying
// where in an HTTP request the flag values should be placed (e.g. query, body, etc.).
type Flag[
	T []any | []DateTimeValue | []DateValue | []TimeValue | []string |
		[]float64 | []int64 | []bool | any | DateTimeValue | DateValue | TimeValue |
		string | float64 | int64 | bool,
] struct {
	Name        string        // name of the flag
	Category    string        // category of the flag, if any
	DefaultText string        // default text of the flag for usage purposes
	HideDefault bool          // whether to hide the default value in output
	Usage       string        // usage string for help output
	Required    bool          // whether the flag is required or not
	Hidden      bool          // whether to hide the flag in help output
	Default     T             // default value for this flag if not set by from any source
	Aliases     []string      // aliases that are allowed for this flag
	Validator   func(T) error // custom function to validate this flag value

	QueryPath  string // location in the request query string to put this flag's value
	HeaderPath string // location in the request header to put this flag's value
	BodyPath   string // location in the request body to put this flag's value
	BodyRoot   bool   // if true, then use this value as the entire request body

	// unexported fields for internal use
	count      int       // number of times the flag has been set
	hasBeenSet bool      // whether the flag has been set from env or file
	applied    bool      // whether the flag has been applied to a flag set already
	value      cli.Value // value representing this flag's value
}

// Type assertions to verify we implement the relevant urfave/cli interfaces
var _ cli.CategorizableFlag = (*Flag[any])(nil)

// InRequest interface for flags that should be included in HTTP requests
type InRequest interface {
	GetQueryPath() string
	GetHeaderPath() string
	GetBodyPath() string
	IsBodyRoot() bool
}

func (f Flag[T]) GetQueryPath() string {
	return f.QueryPath
}

func (f Flag[T]) GetHeaderPath() string {
	return f.HeaderPath
}

func (f Flag[T]) GetBodyPath() string {
	return f.BodyPath
}

func (f Flag[T]) IsBodyRoot() bool {
	return f.BodyRoot
}

// The values that will be sent in different parts of a request.
type RequestContents struct {
	Queries map[string]any
	Headers map[string]any
	Body    any
}

// Extract query parameters, headers, and body values from command flags.
func ExtractRequestContents(cmd *cli.Command) RequestContents {
	bodyMap := make(map[string]any)
	res := RequestContents{
		Queries: make(map[string]any),
		Headers: make(map[string]any),
		Body:    bodyMap,
	}

	for _, flag := range cmd.Flags {
		if !flag.IsSet() {
			continue
		}

		value := flag.Get()
		if toSend, ok := flag.(InRequest); ok {
			if queryPath := toSend.GetQueryPath(); queryPath != "" {
				res.Queries[queryPath] = value
			}
			if headerPath := toSend.GetHeaderPath(); headerPath != "" {
				res.Headers[headerPath] = value
			}
			if toSend.IsBodyRoot() {
				res.Body = value
			} else if bodyPath := toSend.GetBodyPath(); bodyPath != "" {
				bodyMap[bodyPath] = value
			}
		}
	}
	return res
}

// Implementation of the cli.Flag interface
var _ cli.Flag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) PreParse() error {
	newVal := f.Default
	f.value = &cliValue[T]{newVal}

	// Validate the given default or values set from external sources as well
	if f.Validator != nil {
		if err := f.Validator(f.value.Get().(T)); err != nil {
			return err
		}
	}
	f.applied = true
	return nil
}

func (f *Flag[T]) PostParse() error {
	return nil
}

func (f *Flag[T]) Set(_ string, val string) error {
	// Initialize flag if needed
	if !f.applied {
		if err := f.PreParse(); err != nil {
			return err
		}
		f.applied = true
	}

	f.count++

	// If this is the first time setting a slice type, reset it to empty
	// to avoid appending to the default value
	if f.count == 1 && f.value != nil {
		typ := reflect.TypeOf(f.Default)
		if typ != nil && typ.Kind() == reflect.Slice {
			// Create a new empty slice of the same type and set it
			emptySlice := reflect.MakeSlice(typ, 0, 0).Interface()
			f.value = &cliValue[T]{emptySlice.(T)}
		}
	}

	if err := f.value.Set(val); err != nil {
		return err
	}

	f.hasBeenSet = true

	if f.Validator != nil {
		if err := f.Validator(f.value.Get().(T)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Flag[T]) Get() any {
	if f.value != nil {
		return f.value.Get()
	}
	return f.Default
}

func (f *Flag[T]) String() string {
	return cli.FlagStringer(f)
}

func (f *Flag[T]) IsSet() bool {
	return f.hasBeenSet
}

func (f *Flag[T]) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

// Implementation for the cli.VisibleFlag interface
var _ cli.VisibleFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) IsVisible() bool {
	return !f.Hidden
}

func (f *Flag[T]) GetCategory() string {
	return f.Category
}

func (f *Flag[T]) SetCategory(c string) {
	f.Category = c
}

// Implementation for the cli.RequiredFlag interface
var _ cli.RequiredFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) IsRequired() bool {
	return f.Required
}

// Implementation for the cli.DocGenerationFlag interface
var _ cli.DocGenerationFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) TakesValue() bool {
	var t T
	return reflect.TypeOf(t) == nil || reflect.TypeOf(t).Kind() != reflect.Bool
}

func (f *Flag[T]) GetUsage() string {
	return f.Usage
}

func (f *Flag[T]) GetValue() string {
	if f.value == nil {
		return ""
	}
	return f.value.String()
}

func (f *Flag[T]) GetDefaultText() string {
	return f.DefaultText
}

func (f *Flag[T]) GetEnvVars() []string {
	return nil
}

func (f *Flag[T]) IsDefaultVisible() bool {
	return !f.HideDefault
}

func (f *Flag[T]) TypeName() string {
	ty := reflect.TypeOf(f.Default)
	if ty == nil {
		return ""
	}

	// Get base type name with special handling for built-in types
	getTypeName := func(t reflect.Type) string {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "int"
		case reflect.Float32, reflect.Float64:
			return "float"
		case reflect.Bool:
			return "boolean"
		case reflect.String:
			return "string"
		default:
			switch t.Name() {
			case "dateTimeValue":
				return "datetime"
			case "dateValue":
				return "date"
			case "timeValue":
				return "time"
			default:
				if t.Name() == "" {
					return "any"
				}
				return strings.ToLower(t.Name())
			}
		}
	}

	switch ty.Kind() {
	case reflect.Slice:
		elemType := ty.Elem()
		return getTypeName(elemType)
	case reflect.Map:
		keyType := ty.Key()
		valueType := ty.Elem()
		return fmt.Sprintf("%s=%s", getTypeName(keyType), getTypeName(valueType))
	default:
		return getTypeName(ty)
	}
}

// Implementation for the cli.DocGenerationMultiValueFlag interface
var _ cli.DocGenerationMultiValueFlag = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) IsMultiValueFlag() bool {
	if reflect.TypeOf(f.Default) == nil {
		return false
	}
	kind := reflect.TypeOf(f.Default).Kind()
	return kind == reflect.Slice || kind == reflect.Map
}

func (f *Flag[T]) IsBoolFlag() bool {
	_, isBool := any(f.Default).(bool)
	return isBool
}

// Implementation for the cli.Countable interface
var _ cli.Countable = (*Flag[any])(nil) // Type assertion to ensure interface compliance

func (f *Flag[T]) Count() int {
	return f.count
}

// cliValue is a generic implementation of cli.Value for common types
type cliValue[
	T []any | []DateTimeValue | []DateValue | []TimeValue | []string | []float64 |
		[]int64 | []bool | any | DateTimeValue | DateValue | TimeValue | string |
		float64 | int64 | bool,
] struct {
	value T
}

// Set parses the input string and sets the value
func (c *cliValue[T]) Set(value string) error {
	var parsedValue any
	var err error

	switch any(c.value).(type) {
	case []string:
		// Append to existing slice if it exists
		existingSlice, ok := any(c.value).([]string)
		if !ok {
			existingSlice = []string{}
		}
		parsedValue = append(existingSlice, value)

	case []int64:
		var i int64
		i, err = strconv.ParseInt(value, 0, 64)
		if err == nil {
			// Append to existing slice if it exists
			existingSlice, ok := any(c.value).([]int64)
			if !ok {
				existingSlice = []int64{}
			}
			parsedValue = append(existingSlice, i)
		}

	case []float64:
		var f float64
		f, err = strconv.ParseFloat(value, 64)
		if err == nil {
			// Append to existing slice if it exists
			existingSlice, ok := any(c.value).([]float64)
			if !ok {
				existingSlice = []float64{}
			}
			parsedValue = append(existingSlice, f)
		}

	case []bool:
		var b bool
		b, err = strconv.ParseBool(value)
		if err == nil {
			// Append to existing slice if it exists
			existingSlice, ok := any(c.value).([]bool)
			if !ok {
				existingSlice = []bool{}
			}
			parsedValue = append(existingSlice, b)
		}

	case []DateTimeValue:
		var dt DateTimeValue
		err = (&dt).Parse(value)
		if err == nil {
			// Append to existing slice if it exists
			existingSlice, ok := any(c.value).([]DateTimeValue)
			if !ok {
				existingSlice = []DateTimeValue{}
			}
			parsedValue = append(existingSlice, dt)
		}

	case []DateValue:
		var d DateValue
		err = (&d).Parse(value)
		if err == nil {
			// Append to existing slice if it exists
			existingSlice, ok := any(c.value).([]DateValue)
			if !ok {
				existingSlice = []DateValue{}
			}
			parsedValue = append(existingSlice, d)
		}

	case []TimeValue:
		var t TimeValue
		err = (&t).Parse(value)
		if err == nil {
			// Append to existing slice if it exists
			existingSlice, ok := any(c.value).([]TimeValue)
			if !ok {
				existingSlice = []TimeValue{}
			}
			parsedValue = append(existingSlice, t)
		}

	case []any:
		var yamlValue any
		err = yaml.Unmarshal([]byte(value), &yamlValue)
		if err != nil {
			return fmt.Errorf("failed to parse as YAML: %v", err)
		}

		// Append to existing slice if it exists
		existingSlice, ok := any(c.value).([]any)
		if !ok {
			existingSlice = []any{}
		}
		parsedValue = append(existingSlice, yamlValue)

	case string:
		parsedValue = value

	case int64:
		var i int64
		i, err = strconv.ParseInt(value, 0, 64)
		if err == nil {
			parsedValue = i
		}

	case float64:
		var f float64
		f, err = strconv.ParseFloat(value, 64)
		if err == nil {
			parsedValue = f
		}

	case bool:
		var b bool
		b, err = strconv.ParseBool(value)
		if err == nil {
			parsedValue = b
		}

	case DateTimeValue:
		var dt DateTimeValue
		err = (&dt).Parse(value)
		if err == nil {
			parsedValue = dt
		}

	case DateValue:
		var d DateValue
		err = (&d).Parse(value)
		if err == nil {
			parsedValue = d
		}

	case TimeValue:
		var t TimeValue
		err = (&t).Parse(value)
		if err == nil {
			parsedValue = t
		}

	case any:
	default:
		var yamlValue any
		err = yaml.Unmarshal([]byte(value), &yamlValue)
		if err != nil {
			return fmt.Errorf("failed to parse as YAML: %v", err)
		}
		parsedValue = yamlValue
	}

	if err != nil {
		return err
	}

	c.value = parsedValue.(T)
	return nil
}

func (c *cliValue[T]) Get() any {
	return c.value
}

func (c *cliValue[T]) String() string {
	switch v := any(c.value).(type) {
	case string, int, int64, float64, bool, DateTimeValue, DateValue, TimeValue,
		[]string, []int, []int64, []float64, []bool, []DateTimeValue, []DateValue, []TimeValue:
		// For basic types, use standard string representation
		return fmt.Sprintf("%v", v)

	default:
		// For complex types, convert to YAML
		yamlBytes, err := yaml.Marshal(c.value)
		if err != nil {
			// Fall back to standard format if YAML conversion fails
			return fmt.Sprintf("%v", c.value)
		}
		return string(yamlBytes)
	}
}

func (c *cliValue[T]) IsBoolFlag() bool {
	_, ok := any(c.value).(bool)
	return ok
}

// Time-related value types
type DateValue string
type DateTimeValue string
type TimeValue string

// String methods for time-related types
func (d DateValue) String() string {
	return string(d)
}

func (d DateTimeValue) String() string {
	return string(d)
}

func (t TimeValue) String() string {
	return string(t)
}

// parseTimeWithFormats attempts to parse a string using multiple formats
func parseTimeWithFormats(s string, formats []string) (time.Time, error) {
	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

// Parse methods for time-related types
func (d *DateValue) Parse(s string) error {
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
		"2-Jan-2006",
	}

	t, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse date: %v", err)
	}

	*d = DateValue(t.Format("2006-01-02"))
	return nil
}

func (d *DateTimeValue) Parse(s string) error {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.RFC1123,
		time.RFC822,
		time.ANSIC,
	}

	t, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse datetime: %v", err)
	}

	*d = DateTimeValue(t.Format(time.RFC3339))
	return nil
}

func (t *TimeValue) Parse(s string) error {
	formats := []string{
		"15:04:05",
		"3:04:05PM",
		"3:04 PM",
		"15:04",
		time.Kitchen,
	}

	parsedTime, err := parseTimeWithFormats(s, formats)
	if err != nil {
		return fmt.Errorf("unable to parse time: %v", err)
	}

	*t = TimeValue(parsedTime.Format("15:04:05"))
	return nil
}
