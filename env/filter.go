package env

import (
	"reflect"
	"strings"

	"github.com/traefik/paerser/parser"
)

// FindPrefixedEnvVars finds prefixed environment variables.
func FindPrefixedEnvVars(environ []string, prefix string, element interface{}) []string {
	prefixes := getRootPrefixes(element, prefix)

	var values []string
	for _, px := range prefixes {
		for _, value := range environ {
			if strings.HasPrefix(value, px) {
				values = append(values, value)
			}
		}
	}

	return values
}

func getRootPrefixes(element interface{}, prefix string) []string {
	if element == nil {
		return nil
	}

	rootType := reflect.TypeOf(element)

	return getPrefixes(prefix, rootType)
}

func getPrefixes(prefix string, rootType reflect.Type) []string {
	if rootType.Kind() == reflect.Pointer {
		rootType = rootType.Elem()
	}

	if rootType.Kind() != reflect.Struct {
		return nil
	}

	var names []string
	for i := range rootType.NumField() {
		field := rootType.Field(i)

		if !parser.IsExported(field) {
			continue
		}

		if field.Anonymous &&
			(field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct || field.Type.Kind() == reflect.Struct) {
			names = append(names, getPrefixes(prefix, field.Type)...)
			continue
		}

		names = append(names, prefix+strings.ToUpper(field.Name))
	}

	return names
}
