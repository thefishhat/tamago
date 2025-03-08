package server

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func findField(component reflect.Value, fieldPath string) (reflect.Value, error) {
	if fieldPath == "" {
		return component, nil
	}

	fields := strings.Split(fieldPath, ".")

	for _, field := range fields {
		// Process slice or map indexing while there are brackets []
		for strings.Contains(field, "[") && strings.Contains(field, "]") {
			fieldName := field[:strings.Index(field, "[")]
			keyStr := field[strings.Index(field, "[")+1 : strings.Index(field, "]")]

			// If there's a field name before [index], access it as a struct field
			if fieldName != "" {
				if component.Kind() == reflect.Ptr {
					component = component.Elem()
				}
				if component.Kind() == reflect.Struct {
					component = component.FieldByName(fieldName)
				} else {
					return reflect.Value{}, errors.New("invalid field access")
				}
			}

			// Handle slices
			if component.Kind() == reflect.Slice {
				index, err := strconv.Atoi(keyStr)
				if err != nil || index < 0 || index >= component.Len() {
					return reflect.Value{}, errors.New("invalid slice index")
				}
				component = component.Index(index)
			} else if component.Kind() == reflect.Map {
				// Handle maps with string keys
				key := reflect.ValueOf(keyStr)
				component = component.MapIndex(key)
				if !component.IsValid() {
					return reflect.Value{}, errors.New("invalid map key")
				}
			} else {
				return reflect.Value{}, errors.New("invalid index access (not a slice or map)")
			}

			// Remove processed part from the field and check if there are more indices
			field = field[strings.Index(field, "]")+1:]
			if len(field) > 0 && field[0] == '.' {
				field = field[1:] // Remove leading dot if present
				break
			}
		}

		// Handle struct field access
		if field != "" {
			if component.Kind() == reflect.Ptr && component.Elem().Kind() == reflect.Struct {
				component = component.Elem()
			}
			if component.Kind() == reflect.Struct {
				component = component.FieldByName(field)
			} else if component.Kind() == reflect.Map {
				key := reflect.ValueOf(field)
				component = component.MapIndex(key)
				if !component.IsValid() {
					return reflect.Value{}, errors.New("invalid map key")
				}
			} else {
				return reflect.Value{}, errors.New("invalid field access")
			}
		}

		// Dereference pointers and interfaces
		if component.Kind() == reflect.Interface || component.Kind() == reflect.Ptr {
			component = component.Elem()
			if !component.IsValid() {
				return reflect.Value{}, nil
			}
		}
	}

	if !component.IsValid() {
		return reflect.Value{}, errors.New("invalid field access")
	}

	return component, nil
}

func GetField(component reflect.Value, fieldPath string) (interface{}, error) {
	field, err := findField(component, fieldPath)
	if err != nil {
		return nil, err
	}

	fieldVal := recursivelyConstructValue(field, 1)
	if fieldVal == nil {
		return nil, nil
	}

	// double quote string values
	if str, ok := fieldVal.(string); ok {
		return strconv.Quote(str), nil
	}

	return fieldVal, nil
}

func SetField(component reflect.Value, fieldPath string, value interface{}) error {
	field, err := findField(component, fieldPath)
	if err != nil {
		return err
	}

	if !field.CanSet() {
		return errors.New("field is not settable")
	}

	if field.Kind() == reflect.Slice || field.Kind() == reflect.Map {
		return errors.New("cannot set value of slice or map directly")
	}

	if value == nil {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	val := reflect.ValueOf(value)
	fieldType := field.Type()
	if !val.Type().ConvertibleTo(fieldType) {
		return fmt.Errorf("cannot set field: value type %s is not convertible to %s", val.Type(), fieldType)
	}

	field.Set(val.Convert(fieldType))
	return nil
}

func recursivelyConstructValue(value reflect.Value, depth int) interface{} {
	if depth <= 0 {
		if !value.IsValid() {
			return nil
		}
		if value.Kind() == reflect.Slice {
			return fmt.Sprintf("slice of %s", value.Type().Elem().String())
		} else if value.Kind() == reflect.Map {
			return fmt.Sprintf("map of %s to %s", value.Type().Key().String(), value.Type().Elem().String())
		}
		return value.Kind().String()
	}
	switch value.Kind() {
	case reflect.Ptr:
		return recursivelyConstructValue(value.Elem(), depth)
	case reflect.Struct:
		fields := make(map[string]interface{})
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			fields[value.Type().Field(i).Name] = recursivelyConstructValue(field, depth-1)
		}
		return fields
	case reflect.Slice:
		slice := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			slice[i] = recursivelyConstructValue(value.Index(i), depth-1)
		}
		return slice
	default:
		if !value.IsValid() {
			return nil
		}
		if !value.CanInterface() {
			return "Unexported field"
		}
		return value.Interface()
	}
}
