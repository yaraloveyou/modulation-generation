package transform

import (
	"strings"
)

type DataType struct {
	FullName  string
	Formating string
}

var (
	dataTypes = map[string]DataType{
		"int":    {"Integer", "%d"},
		"string": {"String", "%s"},
		"float":  {"Double", "%.2f"},
		"bool":   {"Boolean", "%b"},
		"char":   {"Character", "%c"},
	}
)

func TransformDataType(s string) string {
	if v, ok := dataTypes[s]; ok {
		return v.FullName
	}
	return "Unknown"
}

func CamelCaseToSnakeCase(s string) string {
	var result string

	for i, v := range s {
		if i > 0 && v >= 'A' && v <= 'Z' {
			result += "_"
		}
		result += string(v)
	}

	return strings.ToLower(result)
}

func TransformDataTypeToFormat(s string) string {
	if v, ok := dataTypes[s]; ok {
		return v.Formating
	}
	return "%s"
}
