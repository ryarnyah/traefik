package types

import (
	"fmt"
	"strings"
)

const (
	// AccessLogKeep is the keep string value
	AccessLogKeep = "keep"
	// AccessLogDrop is the drop string value
	AccessLogDrop = "drop"
	// AccessLogRedact is the redact string value
	AccessLogRedact = "redact"
)

// TraefikLog holds the configuration settings for the traefik logger.
type TraefikLog struct {
	FilePath string `json:"file,omitempty" description:"Traefik log file path. Stdout is used when omitted or empty"`
	Format   string `json:"format,omitempty" description:"Traefik log format: json | common"`
}

// AccessLog holds the configuration settings for the access logger (middlewares/accesslog).
type AccessLog struct {
	FilePath string            `json:"file,omitempty" description:"Access log file path. Stdout is used when omitted or empty" export:"true"`
	Format   string            `json:"format,omitempty" description:"Access log format: json | common" export:"true"`
	Filters  *AccessLogFilters `json:"filters,omitempty" description:"Access log filters, used to keep only specific access logs" export:"true"`
	Fields   *AccessLogFields  `json:"fields,omitempty" description:"AccessLogFields" export:"true"`
	Async    *AccessLogAsync   `json:"async,omitempty" description:"Async configuration" export:"true"`
}

// AccessLogAsync holds async log writer configuration for access log
type AccessLogAsync struct {
	AsyncWriterChanSize int64 `json:"asyncWriterChanSize,omitempty" description:"Async Writer chan size" export:"true"`
}

// StatusCodes holds status codes ranges to filter access log
type StatusCodes []string

// AccessLogFilters holds filters configuration
type AccessLogFilters struct {
	StatusCodes   StatusCodes `json:"statusCodes,omitempty" description:"Keep access logs with status codes in the specified range" export:"true"`
	RetryAttempts bool        `json:"retryAttempts,omitempty" description:"Keep access logs when at least one retry happened" export:"true"`
}

// FieldNames holds maps of fields with specific mode
type FieldNames map[string]string

// AccessLogFields holds configuration for access log fields
type AccessLogFields struct {
	DefaultMode string        `json:"defaultMode,omitempty" description:"Default mode for fields: keep | drop" export:"true"`
	Names       FieldNames    `json:"names,omitempty" description:"Override mode for fields" export:"true"`
	Headers     *FieldHeaders `json:"headers,omitempty" description:"Headers to keep, drop or redact" export:"true"`
}

// FieldHeaderNames holds maps of fields with specific mode
type FieldHeaderNames map[string]string

// FieldHeaders holds configuration for access log headers
type FieldHeaders struct {
	DefaultMode string           `json:"defaultMode,omitempty" description:"Default mode for fields: keep | drop | redact" export:"true"`
	Names       FieldHeaderNames `json:"names,omitempty" description:"Override mode for headers" export:"true"`
}

// Set adds strings elem into the the parser
// it splits str on , and ;
func (s *StatusCodes) Set(str string) error {
	fargs := func(c rune) bool {
		return c == ',' || c == ';'
	}
	// get function
	slice := strings.FieldsFunc(str, fargs)
	*s = append(*s, slice...)
	return nil
}

// Get StatusCodes
func (s *StatusCodes) Get() interface{} { return *s }

// String return slice in a string
func (s *StatusCodes) String() string { return fmt.Sprintf("%v", *s) }

// SetValue sets StatusCodes into the parser
func (s *StatusCodes) SetValue(val interface{}) {
	*s = val.(StatusCodes)
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (f *FieldNames) String() string {
	return fmt.Sprintf("%+v", *f)
}

// Get return the FieldNames map
func (f *FieldNames) Get() interface{} {
	return *f
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a space-separated list, so we split it.
func (f *FieldNames) Set(value string) error {
	fields := strings.Fields(value)

	for _, field := range fields {
		n := strings.SplitN(field, "=", 2)
		if len(n) == 2 {
			(*f)[n[0]] = n[1]
		}
	}

	return nil
}

// SetValue sets the FieldNames map with val
func (f *FieldNames) SetValue(val interface{}) {
	*f = val.(FieldNames)
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (f *FieldHeaderNames) String() string {
	return fmt.Sprintf("%+v", *f)
}

// Get return the FieldHeaderNames map
func (f *FieldHeaderNames) Get() interface{} {
	return *f
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a space-separated list, so we split it.
func (f *FieldHeaderNames) Set(value string) error {
	fields := strings.Fields(value)

	for _, field := range fields {
		n := strings.SplitN(field, "=", 2)
		(*f)[n[0]] = n[1]
	}

	return nil
}

// SetValue sets the FieldHeaderNames map with val
func (f *FieldHeaderNames) SetValue(val interface{}) {
	*f = val.(FieldHeaderNames)
}

// Keep check if the field need to be kept or dropped
func (f *AccessLogFields) Keep(field string) bool {
	defaultKeep := true
	if f != nil {
		defaultKeep = checkFieldValue(f.DefaultMode, defaultKeep)

		if v, ok := f.Names[field]; ok {
			return checkFieldValue(v, defaultKeep)
		}
	}
	return defaultKeep
}

func checkFieldValue(value string, defaultKeep bool) bool {
	switch value {
	case AccessLogKeep:
		return true
	case AccessLogDrop:
		return false
	default:
		return defaultKeep
	}
}

// KeepHeader checks if the headers need to be kept, dropped or redacted and returns the status
func (f *AccessLogFields) KeepHeader(header string) string {
	defaultValue := AccessLogKeep
	if f != nil && f.Headers != nil {
		defaultValue = checkFieldHeaderValue(f.Headers.DefaultMode, defaultValue)

		if v, ok := f.Headers.Names[header]; ok {
			return checkFieldHeaderValue(v, defaultValue)
		}
	}
	return defaultValue
}

func checkFieldHeaderValue(value string, defaultValue string) string {
	if value == AccessLogKeep || value == AccessLogDrop || value == AccessLogRedact {
		return value
	}
	return defaultValue
}
