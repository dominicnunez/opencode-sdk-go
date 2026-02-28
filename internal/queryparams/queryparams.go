package queryparams

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Marshal converts a struct with query:"name" tags to url.Values.
// This is a simplified version that handles the basic types used in the SDK:
// - string and *string (optional)
// - int, int64 and *int, *int64 (optional)
// - bool and *bool (optional)
// - []string (array of strings)
//
// Empty non-pointer strings are always omitted from the output, regardless
// of whether "omitempty" is set. This means the "required" tag on a non-pointer
// string enforces "required and non-empty". To send an empty string value,
// use *string.
func Marshal(v interface{}) (url.Values, error) {
	if v == nil {
		return url.Values{}, nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("queryparams: expected struct, got %v", val.Kind())
	}

	params := url.Values{}
	if err := marshalFields(params, val); err != nil {
		return nil, err
	}
	return params, nil
}

func marshalFields(params url.Values, val reflect.Value) error {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := typ.Field(i)

		// Recurse into embedded structs
		if typeField.Anonymous && field.Kind() == reflect.Struct {
			if err := marshalFields(params, field); err != nil {
				return err
			}
			continue
		}

		// Skip unexported fields
		if !typeField.IsExported() {
			continue
		}

		// Get query tag
		tag := typeField.Tag.Get("query")
		if tag == "" || tag == "-" {
			continue
		}

		// Parse tag to extract name and options
		name, required, omitempty := parseTag(tag)
		if name == "" {
			continue
		}

		// Add field value to params
		if err := addFieldValue(params, name, field, required, omitempty); err != nil {
			return fmt.Errorf("queryparams: field %s: %w", typeField.Name, err)
		}
	}

	return nil
}

// parseTag parses a struct tag like "name,required", "name,omitempty", or "name,required,omitempty"
// Returns: name, required, omitempty
func parseTag(tag string) (name string, required bool, omitempty bool) {
	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, opt := range parts[1:] {
		switch opt {
		case "required":
			required = true
		case "omitempty":
			omitempty = true
		}
	}
	return
}

func addFieldValue(params url.Values, name string, field reflect.Value, required, omitempty bool) error {
	isPtr := field.Kind() == reflect.Ptr

	// Handle pointer types
	if isPtr {
		if field.IsNil() {
			if required {
				return fmt.Errorf("required query parameter %q is nil", name)
			}
			return nil // omit nil pointers
		}
		field = field.Elem()
	}

	// Handle based on type
	switch field.Kind() {
	case reflect.String:
		s := field.String()
		if s == "" {
			if required {
				return fmt.Errorf("required query parameter %q is empty", name)
			}
			return nil
		}
		params.Add(name, s)
	case reflect.Int, reflect.Int64:
		if omitempty && !isPtr && field.Int() == 0 {
			return nil
		}
		params.Add(name, strconv.FormatInt(field.Int(), 10))
	case reflect.Bool:
		if omitempty && !isPtr && !field.Bool() {
			return nil
		}
		params.Add(name, strconv.FormatBool(field.Bool()))
	case reflect.Slice:
		if required && field.Len() == 0 {
			return fmt.Errorf("required query parameter %q is empty", name)
		}
		// Handle []string
		if field.Type().Elem().Kind() == reflect.String {
			for j := 0; j < field.Len(); j++ {
				params.Add(name, field.Index(j).String())
			}
		} else {
			return fmt.Errorf("unsupported slice element type: %v", field.Type().Elem().Kind())
		}
	default:
		return fmt.Errorf("unsupported query param type: %v", field.Kind())
	}

	return nil
}
