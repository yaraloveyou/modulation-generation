package models

import (
	"fmt"
	"fois-generator/internal/enums"
	"fois-generator/internal/transform"
	"fois-generator/internal/utils"
	"strings"
)

type Method struct {
	Name              string
	ClassName         string
	Modifier          string
	Annotations       []string
	ExternalVariables []Variable
	Variables         []Variable
	Return            Variable
}

const (
	GETTER_TEMPLATE = `%s %s %s() {
		return %s;
	}`
	SETTER_TEMPLATE = `%s %s %s(%s %s) {
		this.%s = %s;
	}`
	CONSTRUCTOR_TEMPLATE = `%s %s(%s) {
%s
	}`
	TOSTRING_TEMPLATE = `%s %s %s() {
		return String.format(
			"%s[%s]"%s
		);
	}`
)

func (method *Method) determineTemplate() string {
	if strings.Contains(method.Name, "get") {
		return GETTER_TEMPLATE
	} else if strings.Contains(method.Name, "set") {
		return SETTER_TEMPLATE
	} else if strings.Contains(method.Name, method.ClassName) {
		return CONSTRUCTOR_TEMPLATE
	} else if strings.Contains(method.Name, "toString") {
		return TOSTRING_TEMPLATE
	}
	return ""
}

func (method *Method) GenerateStringMethod() string {
	template := method.determineTemplate()
	var mtd string
	switch template {
	case GETTER_TEMPLATE:
		mtd = method.GenerateStringGetter()
	case SETTER_TEMPLATE:
		mtd = method.GenerateStringSetter()
	case CONSTRUCTOR_TEMPLATE:
		mtd = method.GenerateStringConstructors()
	case TOSTRING_TEMPLATE:
		mtd = method.GenerateStringToString()
	}
	mtd = utils.AddAnnotations(method.Annotations, mtd, enums.Method)

	return fmt.Sprintf("\n%s\n\n", mtd)
}

func (method *Method) GenerateStringConstructors() string {
	var builder strings.Builder
	var params string
	var eq string
	for i, param := range method.Variables {
		params += fmt.Sprintf("%s %s", param.DataType, param.Name)
		eq += fmt.Sprintf("\t\tthis.%s = %s;", param.Name, param.Name)
		if i < len(method.Variables)-1 {
			params += ", "
			eq += "\n"
		}
	}
	fmt.Fprintf(&builder, CONSTRUCTOR_TEMPLATE, method.Modifier, method.Name, params, eq)
	return builder.String()
}

func (method *Method) GenerateStringSetter() string {
	return fmt.Sprintf(
		SETTER_TEMPLATE,
		method.Modifier,
		method.Return.DataType,
		method.Name,
		method.Variables[0].DataType,
		method.Variables[0].Name,
		method.ExternalVariables[0].Name,
		method.Variables[0].Name,
	)
}

func (method *Method) GenerateStringGetter() string {
	return fmt.Sprintf(
		GETTER_TEMPLATE,
		method.Modifier,
		method.Return.DataType,
		method.Name,
		method.Return.Name,
	)
}

func (method *Method) GenerateStringToString() string {
	var format string
	var variables string
	for i, v := range method.ExternalVariables {
		format += fmt.Sprintf("%s=%s", v.Name, transform.TransformDataTypeToFormat(v.DataType))
		variables += fmt.Sprintf(", %s", v.Name)
		if i < len(method.ExternalVariables)-1 {
			format += ", "
		}
	}
	return fmt.Sprintf(
		TOSTRING_TEMPLATE,
		method.Modifier,
		method.Return.DataType,
		method.Name,
		method.ClassName,
		format,
		variables,
	)
}
