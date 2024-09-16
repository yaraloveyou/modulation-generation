package models

import (
	"fmt"
	"fois-generator/internal/transform"
	jsonmodels "fois-generator/models/json_models"
	"strings"
)

type Class struct {
	Name     string
	Modifier string
	Fields   []Field
	Methods  []Method
	Layer    string // Entity, Repository, Controller, Exception
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

func (class *Class) createMethods() {
	for _, field := range class.Fields {
		methods := field.GenerateMethods()
		class.Methods = append(class.Methods, methods...)
	}
}

func (class *Class) formationEntity(table jsonmodels.Table) {
	class.createFields(table)
	class.createMethods()
	class.Layer = "Entity"
}

func propertyToField(key, value string) Field {
	field := Field{
		Name:        key,
		Modifier:    "protected",
		Annotations: []string{fmt.Sprintf("Column(name=\"%s\")", transform.CamelCaseToSnakeCase(key))},
	}

	values := strings.Split(value, ";")
	for _, val := range values {
		switch val {
		case "int", "string", "float", "char", "boolean":
			field.DataType = transform.TransformDataType(val)
		case "identity":
			field.Annotations = append([]string{"Id", "GeneratedValue(strategy= GenerationType.IDENTITY)"}, field.Annotations...)
		case "foreign_key":
			// Логика для foreign_key
		}
	}

	return field
}

func (class *Class) GeneratedEntity(table jsonmodels.Table) {
	var builder strings.Builder
	class.formationEntity(table)

	fmt.Fprintf(&builder, "@%s\n", class.Layer)
	fmt.Fprintf(&builder, "%s class %s implements Serializable {\n", class.Modifier, class.Name)
	for _, field := range class.Fields {
		fmt.Fprintf(&builder, "%s", field.GeneratedStringField())
	}

	for i := 0; i < len(class.Methods); i += 2 {
		fmt.Fprintf(&builder, "%s", class.Methods[i].GenerateStringMethods())
	}
	fmt.Fprintf(&builder, "}")

	fmt.Println(builder.String())
}
