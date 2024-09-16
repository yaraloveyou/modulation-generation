package models

import (
	"fmt"
	"strings"
)

type Field struct {
	Name        string
	DataType    string
	Modifier    string
	Annotations []string //@Id @GeneratedValue() @Column()
}

func (field *Field) GenerateMethods(class *Class) []Method {
	getter := Method{
		Name:     fmt.Sprintf("get%s", strings.Title(field.Name)),
		Modifier: class.Modifier,
		ExternalVariables: []Variable{
			{
				Name:     field.Name,
				DataType: field.DataType,
			},
		},
		Position: "Entity",
		Return: Variable{
			Name:     field.Name,
			DataType: field.DataType,
		},
		ClassName: class.Name,
	}
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
		Position: "Entity",
		Return: Variable{
			DataType: "void",
		},
		ClassName: class.Name,
	}

	return []Method{getter, setter}
}

func (field *Field) GenerateStringField() string {
	var builder strings.Builder

	for _, ent := range field.Annotations {
		fmt.Fprintf(&builder, "\t@%s\n", ent)
	}

	fmt.Fprintf(&builder, "\t%s %s %s;\n\n", field.Modifier, field.DataType, field.Name)

	return builder.String()
}
