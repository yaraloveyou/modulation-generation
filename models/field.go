package models

import (
	"fmt"
	"fois-generator/internal/enums"
	"fois-generator/internal/utils"
	"strings"
)

type Field struct {
	Name        string
	DataType    string
	Modifier    string
	Position    string
	Annotations []string //@Id @GeneratedValue() @Column()
}

func (field *Field) GenerateGetter(class *Class) Method {
	getter := Method{
		Name:     fmt.Sprintf("get%s", strings.Title(field.Name)),
		Modifier: class.Modifier,
		ExternalVariables: []Variable{
			{
				Name:     field.Name,
				DataType: field.DataType,
			},
		},
		Return: Variable{
			Name:     field.Name,
			DataType: field.DataType,
		},
		ClassName: class.Name,
	}

	return getter
}

func (field *Field) GenerateSetter(class *Class) Method {
	setter := Method{
		Name:     fmt.Sprintf("set%s", strings.Title(field.Name)),
		Modifier: class.Modifier,
		ExternalVariables: []Variable{
			{
				Name:     field.Name,
				DataType: field.DataType,
			},
		},
		Variables: []Variable{
			{
				Name:     field.Name,
				DataType: field.DataType,
			},
		},
		Return: Variable{
			DataType: "void",
		},
		ClassName: class.Name,
	}

	return setter
}

func (field *Field) GenerateStringField() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%s %s %s;\n\n", field.Modifier, field.DataType, field.Name)
	return utils.AddAnnotations(field.Annotations, builder.String(), enums.Field)
}
