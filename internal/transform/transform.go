package transform

import (
	"strings"
)

func TransformDataType(s string) string {
	switch s {
	case "int":
		return "Integer"
	case "string":
		return "String"
	case "float":
		return "Double"
	case "bool":
		return "Boolean"
	case "char":
		return "Character"
	default:
		return "Unknown type"
	}
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
	switch s {
	case "Integer":
		return "%d"
	case "String":
		return "%s"
	case "Double":
		return "%.2f"
	case "Boolean":
		return "%b"
	case "Character":
		return "%c"
	default:
		return "%s"
	}
}
