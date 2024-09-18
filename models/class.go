package models

import (
	"fmt"
	"fois-generator/internal/enums"
	"fois-generator/internal/transform"
	"fois-generator/internal/utils"
	jsonmodels "fois-generator/models/json_models"
	"strings"
)

type Class struct {
	Name        string
	Modifier    string
	Fields      []Field
	Methods     []Method
	Position    string
	Annotations []string // Entity, Repository, Controller, Exception
}

func (class *Class) createFields(table jsonmodels.Table) {
	var fields []Field
	//{Table1 map[id:int;identity name:string user_id:int;foreign_key{user};]}

	fieldsMap, ok := table.Fields.(map[string]interface{})
	if !ok {
		return
	}

	for key, value := range fieldsMap {
		field := propertyToField(key, value.(string))
		fields = append(fields, field)
	}

	class.Fields = append(class.Fields, fields...)
}

func propertyToField(key, value string) Field {
	field := Field{
		Name:        key,
		Modifier:    "private",
		Annotations: []string{fmt.Sprintf(enums.Column, transform.CamelCaseToSnakeCase(key))},
	}

	values := strings.Split(value, ";")
	for _, val := range values {
		switch val {
		case "int", "string", "float", "char", "boolean":
			field.DataType = transform.TransformDataType(val)
		case "identity":
			field.Annotations = append(
				[]string{enums.Id, fmt.Sprintf(enums.GeneratedValue, "strategy = GenerationType.IDENTITY")},
				field.Annotations...,
			)
		case "foreign_key":
			// Логика для foreign_key
		}
	}

	return field
}

func (class *Class) GenerateEntity(table jsonmodels.Table) string {
	var builder strings.Builder
	class.formationEntity(table)
	class.EntityLombokAnnotations()
	class.Position = "Entity"
	fmt.Fprintf(&builder, "%s class %s implements Serializable {\n", class.Modifier, class.Name)
	for _, field := range class.Fields {
		fmt.Fprintf(&builder, "%s", field.GenerateStringField())
	}
	for i := 0; i < len(class.Methods); i++ {
		fmt.Fprintf(&builder, "%s", class.Methods[i].GenerateStringMethod())
	}
	fmt.Fprintf(&builder, "}\n")

	return utils.AddAnnotations(class.Annotations, builder.String(), enums.Class)
}
