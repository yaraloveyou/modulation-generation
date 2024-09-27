package models

import (
	"fmt"
	"fois-generator/internal/enums"
	"fois-generator/internal/transform"
	"fois-generator/internal/utils"
	"html/template"
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
	GETTER_TEMPLATE = `{{.Modifier}} {{.Return.DataType}} {{.Name}}({{StringsJoin .ExternalVariables ", "}}) {
		return {{.Return.Name}}
	}`

	SETTER_TEMPLATE = `{{.Modifier}} {{.Return.DataType}} {{.Name}}({{with index .Variables 0}}{{.DataType}} {{.Name}}{{end}}) {
		this.{{with index .ExternalVariables 0}}{{.Name}}{{end}} = {{with index .Variables 0}}{{.Name}}{{end}};
	}`

	CONSTRUCTOR_TEMPLATE = `{{.Modifier}} {{.Name}}({{StringsJoin .ExternalVariables ", "}}) {
		{{VarJoin .ExternalVariables "\n\t\t"}}	
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
		mtd, _ = method.GenerateStringGetter()
	case SETTER_TEMPLATE:
		mtd, _ = method.GenerateStringSetter()
	case CONSTRUCTOR_TEMPLATE:
		mtd, _ = method.GenerateStringConstructors()
	case TOSTRING_TEMPLATE:
		mtd = method.GenerateStringToString()
	}
	mtd = utils.AddAnnotations(method.Annotations, mtd, enums.Method)

	return fmt.Sprintf("\n%s\n\n", mtd)
}

func (method *Method) GenerateStringConstructors() (string, error) {
	t, err := template.New("constructor").Funcs(template.FuncMap{
		"StringsJoin": func(variables []Variable, sep string) string {
			var res []string
			var builder strings.Builder
			for _, v := range variables {
				fmt.Fprintf(&builder, "%s %s", v.DataType, v.Name)
				res = append(res, builder.String())
				builder.Reset()
			}
			return strings.Join(res, sep)
		},
		"VarJoin": func(variables []Variable, sep string) string {
			var res []string
			var builder strings.Builder
			for _, v := range variables {
				fmt.Fprintf(&builder, "this.%s = %s;", v.Name, v.Name)
				res = append(res, builder.String())
				builder.Reset()
			}
			return strings.Join(res, sep)
		},
	}).Parse(CONSTRUCTOR_TEMPLATE)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	err = t.Execute(&builder, method)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (method *Method) GenerateStringSetter() (string, error) {
	t, err := template.New("setter").Parse(SETTER_TEMPLATE)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	err = t.Execute(&builder, method)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (method *Method) GenerateStringGetter() (string, error) {
	t, err := template.New("getter").Funcs(template.FuncMap{
		"StringsJoin": func(variables []Variable, sep string) string {
			var res []string
			var builder strings.Builder
			for _, v := range variables {
				fmt.Fprintf(&builder, "%s %s", v.DataType, v.Name)
				res = append(res, builder.String())
				builder.Reset()
			}
			return strings.Join(res, sep)
		},
	}).Parse(GETTER_TEMPLATE)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	err = t.Execute(&builder, method)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
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
