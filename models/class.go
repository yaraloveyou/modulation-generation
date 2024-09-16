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
	Layers   []string // Entity, Repository, Controller, Exception
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
		methods := field.GenerateMethods(class)
		class.Methods = append(class.Methods, methods...)
	}
}

func (class *Class) formationEntity(table jsonmodels.Table) {
	class.createFields(table)
	class.GenerateConstructor()
	class.GenerateEmptyConstructor()
	class.createMethods()

}

func propertyToField(key, value string) Field {
	field := Field{
		Name:        key,
		Modifier:    "private",
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

func (class *Class) GenerateEntity(table jsonmodels.Table) {
	var builder strings.Builder
	class.formationEntity(table)
	for _, layer := range class.Layers {
		fmt.Fprintf(&builder, "%s\n", layer)
	}
	fmt.Fprintf(&builder, "%s class %s implements Serializable {\n", class.Modifier, class.Name)
	for _, field := range class.Fields {
		fmt.Fprintf(&builder, "%s", field.GenerateStringField())
	}
	for i := 0; i < len(class.Methods); i++ {
		fmt.Fprintf(&builder, "%s", class.Methods[i].GenerateStringMethod())
	}
	fmt.Fprintf(&builder, "}\n")

	fmt.Println(builder.String())
}

func (class *Class) GenerateConstructor() {
	constructor := Method{
		Name:      class.Name,
		Modifier:  class.Modifier,
		ClassName: class.Name,
	}

	for _, field := range class.Fields {
		variable := Variable{
			Name:     field.Name,
			DataType: field.DataType,
		}
		constructor.ExternalVariables = append(constructor.ExternalVariables, variable)
		constructor.Variables = append(constructor.Variables, variable)
	}

	class.Methods = append(class.Methods, constructor)
}

func (class *Class) GenerateEmptyConstructor() {
	constructor := Method{
		Name:      class.Name,
		Modifier:  class.Modifier,
		ClassName: class.Name,
	}

	class.Methods = append(class.Methods, constructor)
}
