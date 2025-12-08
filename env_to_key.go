package mykonf

import (
	"reflect"
	"strings"
)

func EnvToKey(structNilPtr any, tag string) map[string]string {
	result := make(map[string]string)
	traverseType(reflect.TypeOf(structNilPtr), "", tag, result)
	return result
}

func traverseType(t reflect.Type, jsonPrefix, tagKey string, result map[string]string) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := range t.NumField() {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue
		}

		tagValue := field.Tag.Get(tagKey)
		jsonName := strings.Split(tagValue, ",")[0]
		if jsonName == "-" {
			continue
		}
		if jsonName == "" {
			jsonName = field.Name
		}

		currentJSONPrefix := jsonName
		if jsonPrefix != "" {
			currentJSONPrefix = jsonPrefix + "." + currentJSONPrefix
		}

		fieldType := field.Type

		ft := fieldType
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			switch ft.String() {
			// Time/Duration are leaf nodes
			case "time.Time", "time.Duration":
			default:
				traverseType(fieldType, currentJSONPrefix, tagKey, result)
			}
		}

		// leaf node
		envKey := strings.ToUpper(strings.ReplaceAll(currentJSONPrefix, ".", "_"))
		result[envKey] = currentJSONPrefix
	}
}
