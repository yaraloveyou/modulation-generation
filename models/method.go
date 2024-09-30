package models

import (
	"fmt"
	"fois-generator/internal/enums"
	"fois-generator/internal/templates"
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

var templateMapping = map[string]string{
	"get":      templates.GETTER_TEMPLATE,
	"set":      templates.SETTER_TEMPLATE,
	"toString": templates.TOSTRING_TEMPLATE,
}

var generatorFunctions = map[string]func(*Method) (string, error){
	"get":      (*Method).GenerateStringGetter,
	"set":      (*Method).GenerateStringSetter,
	"toString": (*Method).GenerateStringToString,
}

func init() {
	generatorFunctions["constructor"] = (*Method).GenerateStringConstructors
}

func (method *Method) determineTemplate() (string, string) {
	for key := range templateMapping {
		if strings.Contains(method.Name, key) {
			return key, templateMapping[key]
		}
	}
	if strings.Contains(method.Name, method.ClassName) {
		return "constructor", templates.CONSTRUCTOR_TEMPLATE
	}
	return "", ""
}

func (method *Method) GenerateStringMethod() (string, error) {
	methodType, tmpl := method.determineTemplate()
	if tmpl == "" {
		return "", fmt.Errorf("template not found for method: %s", method.Name)
	}

	generatorFunc, exists := generatorFunctions[methodType]
	if !exists {
		return "", fmt.Errorf("generation function not found for method: %s", method.Name)
	}

	result, err := generatorFunc(method)
	if err != nil {
		return "", err
	}

	result = utils.AddAnnotations(method.Annotations, result, enums.Method)
	return fmt.Sprintf("\n%s\n", result), nil
}

func (method *Method) GenerateStringConstructors() (string, error) {
	return method.generateStringFromTemplate(templates.CONSTRUCTOR_TEMPLATE, funcMapConstructor)
}

func (method *Method) GenerateStringSetter() (string, error) {
	return method.generateStringFromTemplate(templates.SETTER_TEMPLATE, nil)
}

func (method *Method) GenerateStringGetter() (string, error) {
	return method.generateStringFromTemplate(templates.GETTER_TEMPLATE, funcMapGetterSetter)
}

func (method *Method) GenerateStringToString() (string, error) {
	return method.generateStringFromTemplate(templates.TOSTRING_TEMPLATE, funcMapToString)
}

func (method *Method) generateStringFromTemplate(tmpl string, funcMap template.FuncMap) (string, error) {
	t, err := template.New("method").Funcs(funcMap).Parse(tmpl)
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

var funcMapConstructor = template.FuncMap{
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
	"VarsJoin": func(variables []Variable, sep string) string {
		var res []string
		var builder strings.Builder
		for _, v := range variables {
			fmt.Fprintf(&builder, "this.%s = %s;", v.Name, v.Name)
			res = append(res, builder.String())
			builder.Reset()
		}
		return strings.Join(res, sep)
	},
}

var funcMapGetterSetter = template.FuncMap{
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
}

var funcMapToString = template.FuncMap{
	"StringsJoin": func(variables []Variable, sep string) string {
		var res []string
		var builder strings.Builder
		for _, v := range variables {
			fmt.Fprintf(&builder, "%s=%s", v.Name, transform.TransformDataTypeToFormat(v.DataType))
			res = append(res, builder.String())
			builder.Reset()
		}
		return strings.Join(res, sep)
	},
	"VarsJoin": func(variables []Variable, sep string) string {
		var res []string
		for _, v := range variables {
			res = append(res, v.Name)
		}
		return strings.Join(res, sep)
	},
}
