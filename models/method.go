package models

import (
	"fmt"
)

type Method struct {
	Name     string
	Position string
	DataType string
	Modifier string
}

const (
	GETTER = `%s %s get%s() {
		return %s;
	}`
	SETTER = `%s void set%s(%s %s) {
		this.%s = %s
	}`
)

func (method *Method) GenerateStringMethods() string {
	getter := fmt.Sprintf(
		GETTER,
		method.Modifier, //protected, public, private (для начала автоматом выставляется protected в class.go)
		method.DataType, // тип данных
		method.Name,     // Наименование метода (getName, getAge..)
		method.Name,     // Наименование поля
	)

	setter := fmt.Sprintf(
		SETTER,
		method.DataType,
		method.Name,     // Наименование метода (getName, getAge..)
		method.DataType, // Тип данных
		method.Name,     // Наименование поля
		method.Name,     // Наименование поля
		method.Name,     // Наименование поля
	)

	return fmt.Sprintf("\n\t%s\n\n\t%s\n", getter, setter)
}
