package models

import (
	"fois-generator/internal/enums"
)

func (class *Class) createMethods() {
	for _, field := range class.Fields {
		getter := field.GenerateGetter(class)
		setter := field.GenerateSetter(class)
		class.Methods = append(class.Methods, []Method{getter, setter}...)
	}
}

func (class *Class) formationEntity() {
	// Создание полей
	// Создание конструкторов
	class.GenerateConstructor()
	class.GenerateEmptyConstructor()
	// Создание метода toString()
	class.GenerateToString()
	// Создание геттеров и сеттеров
	class.createMethods()
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

func (class *Class) GenerateToString() {
	var externalVariables []Variable
	for _, f := range class.Fields {
		externalVariables = append(externalVariables, Variable{Name: f.Name, DataType: f.DataType})
	}

	toString := Method{
		Name:              "toString",
		Modifier:          class.Modifier,
		ExternalVariables: externalVariables,
		Return: Variable{
			DataType: "String",
		},
		Annotations: []string{enums.Override},
		ClassName:   class.Name,
	}

	class.Methods = append(class.Methods, toString)
}
