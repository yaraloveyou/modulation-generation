package models

import (
	"fmt"
	"strings"
)

type Field struct {
	Name        string
	DataType    string
	Modifier    string
	Annotations []string //@Id @GeneratedValue() @Column() @
}

func (field *Field) GenerateMethods() []Method {
	getter := Method{
		Name:     strings.Title(field.Name),
		DataType: field.DataType,
		Modifier: field.Modifier,
	}
	setter := Method{
		Name:     strings.Title(field.Name),
		DataType: field.DataType,
		Modifier: field.Modifier,
	}

	return []Method{getter, setter}
}

func (field *Field) GeneratedStringField() string {
	var builder strings.Builder

	for _, ent := range field.Annotations {
		fmt.Fprintf(&builder, "\t@%s\n", ent)
	}

	fmt.Fprintf(&builder, "\t%s %s %s;\n\n", field.Modifier, field.DataType, field.Name)

	return builder.String()
}
