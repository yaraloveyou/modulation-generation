package models

import (
	"fmt"
	"fois-generator/internal/enums"
	"fois-generator/internal/transform"
	"fois-generator/internal/utils"
	jsonmodels "fois-generator/models/json_models"
	"regexp"
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

// Create fields based on the table structure
func (class *Class) CreateFields(table jsonmodels.Table) {
	fieldsMap, ok := table.Fields.(map[string]interface{})
	if !ok {
		return
	}

	for key, value := range fieldsMap {
		field := propertyToField(table.Name, key, value.(string))
		class.Fields = append(class.Fields, field)
	}
}

// Add relationships to the fields
func (class *Class) AddRelated(classes []*Class, table jsonmodels.Table) {
	fieldsMap, ok := table.Fields.(map[string]interface{})
	if !ok {
		return
	}

	for key, value := range fieldsMap {
		addRelatedField(classes, table.Name, key, value.(string))
	}
}

// Transforms table property into field
func propertyToField(tableName, key, value string) Field {
	field := Field{
		Name:        key,
		Modifier:    "private",
		Annotations: []string{fmt.Sprintf(enums.Column, transform.CamelCaseToSnakeCase(key))},
		Position:    tableName,
	}

	// Split value by semicolon to handle attributes like 'int;identity'
	values := strings.Split(value, ";")
	for _, val := range values {
		if transform.IsValidDataType(val) {
			field.DataType = transform.TransformDataType(val)
		}
		if val == "identity" {
			field.Annotations = append(
				[]string{enums.Id, fmt.Sprintf(enums.GeneratedValue, "strategy = GenerationType.IDENTITY")},
				field.Annotations...,
			)
		}
	}
	return field
}

// Handles relationships between classes
func addRelatedField(classes []*Class, tableName, key, value string) {
	field := Field{
		Name:        key,
		Modifier:    "private",
		Annotations: []string{fmt.Sprintf(enums.Column, transform.CamelCaseToSnakeCase(key))},
		Position:    tableName,
	}

	values := strings.Split(value, ";")
	for _, val := range values {
		if strings.HasPrefix(val, "foreign_key") {
			handleForeignKey(classes, tableName, key, val)
		}
	}
}

func handleForeignKey(classes []*Class, tableName, key, foreignKey string) {
	re := regexp.MustCompile(`foreign_key\{([^\}]+)\}`)
	match := re.FindStringSubmatch(foreignKey)

	relatedClasses := regexp.MustCompile(`\s*,\s*`).Split(match[1], -1)
	for _, class := range classes {
		if utils.Contains(relatedClasses, class.Name) {
			class.Fields = append(class.Fields, Field{
				Name:        strings.ToLower(class.Name),
				DataType:    strings.Title(class.Name),
				Modifier:    "private",
				Annotations: []string{enums.ManyToOne},
				Position:    class.Name,
			})
		}
	}
}
