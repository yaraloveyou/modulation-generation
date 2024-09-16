package transform

import "strings"

func TransformDataType(value string) string {
	switch value {
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
