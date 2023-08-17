package urlvaluer

import (
	"fmt"
	"github.com/gaigals/gotags"
	"github.com/spf13/cast"
	"net/url"
	"reflect"
)

const ignoreTag = "-"

type Stringer interface {
	String() string
}

var (
	tagSettings = gotags.NewSettings("url").
		WithNoKeyExistValidation()
)

func Marshal(target any) (url.Values, error) {
	fields, err := tagSettings.ParseStruct(target)
	if err != nil {
		return nil, err
	}

	return parseFields(fields)
}

func parseFields(fields []gotags.Field) (url.Values, error) {
	urlValues := make(url.Values)

	for _, field := range fields {
		isPtrKind := field.Kind == reflect.Ptr || field.Kind == reflect.Interface
		if isPtrKind && field.Value.IsNil() {
			continue
		}
		if hasIgnoreTag(field) {
			continue
		}

		strSlice, err := castAsString(field)
		if err != nil {
			return nil, fmt.Errorf("urlvaluer: %w", err)
		}
		if len(strSlice) == 0 {
			continue
		}

		urlValues[field.FirstTag().Key] = strSlice
	}

	return urlValues, nil
}

func hasIgnoreTag(field gotags.Field) bool {
	return field.FirstTag().Key == ignoreTag
}

func castAsString(field gotags.Field) ([]string, error) {
	str, hasIStr := strFromInterface(field)
	if hasIStr {
		return []string{str}, nil
	}

	if field.Kind == reflect.Struct || field.Kind == reflect.Map {
		return nil, nil
	}

	if field.Kind == reflect.Ptr {
		return ptrAsStringSlice(field)
	}

	if field.Kind != reflect.Slice && field.Kind != reflect.Array {
		return valueAsStringSlice(field)
	}

	return sliceAsStringSlice(field)
}

func ptrAsStringSlice(field gotags.Field) ([]string, error) {
	valueElem := field.Value.Elem()
	newField := gotags.Field{Value: valueElem, Kind: valueElem.Kind()}
	return castAsString(newField)
}

func valueAsStringSlice(field gotags.Field) ([]string, error) {
	str, err := cast.ToStringE(field.Value.Interface())
	if err != nil {
		return nil, fmt.Errorf("field=%s casting error: %w",
			field.Name, err)
	}

	return []string{str}, nil
}

func strFromInterface(field gotags.Field) (string, bool) {
	//stringer, ok := field.Value.Interface().(fmt.Stringer)
	stringer, ok := field.Value.Interface().(Stringer)
	if !ok {
		return tryStrFromPtrInterface(field)
	}

	return stringer.String(), true
}

func tryStrFromPtrInterface(field gotags.Field) (string, bool) {
	if !field.Value.CanAddr() {
		return "", false
	}

	stringer, ok := field.Value.Addr().Interface().(Stringer)
	if !ok {
		return "", false
	}
	return stringer.String(), true
}

func sliceAsStringSlice(field gotags.Field) ([]string, error) {
	strSlice := make([]string, field.Value.Len())

	for idx := range strSlice {
		str, err := cast.ToStringE(field.Value.Index(idx).Interface())
		if err != nil {
			return nil, fmt.Errorf("field=%s casting error: %w",
				field.Name, err)
		}

		strSlice[idx] = str
	}

	return strSlice, nil
}
