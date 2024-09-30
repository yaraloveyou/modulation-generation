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
	Annotations []string
}

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

func (class *Class) AddRelated(allClasses []*Class, table jsonmodels.Table) {
	fieldsMap, ok := table.Fields.(map[string]interface{})
	if !ok {
		return
	}
	for key, value := range fieldsMap {
		class.addManyToOneRelation(allClasses, table.Name, key, value.(string))
		class.addOneToOne(allClasses, table.Name, value.(string))
	}
}

func propertyToField(tableName, key, value string) Field {
	field := Field{
		Name:        key,
		Modifier:    "private",
		Annotations: []string{fmt.Sprintf(enums.Column, transform.CamelCaseToSnakeCase(key))},
		Position:    tableName,
	}

	for _, val := range strings.Split(value, ";") {
		if val != "" && strings.Contains("int|string|float|char|boolean", val) {
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

func (class *Class) addManyToOneRelation(allClasses []*Class, tableName, key, value string) {
	if !strings.Contains(value, "foreign_key") {
		return
	}

	if strings.Contains(value, "unique") {
		return
	}

	re := regexp.MustCompile(`foreign_key\{([^\}]+)\}`)
	match := re.FindStringSubmatch(value)
	if len(match) < 2 {
		return
	}
	relatedClassNames := regexp.MustCompile(`\s*,\s*`).Split(match[1], -1)
	var relClass string

	for _, relClassName := range relatedClassNames {
		for _, c := range allClasses {
			if strings.EqualFold(relClassName, c.Name) {
				newField := Field{
					Name:        strings.ToLower(tableName),
					DataType:    strings.Title(tableName),
					Modifier:    "private",
					Annotations: []string{enums.ManyToOne},
					Position:    c.Name,
				}
				c.Fields = append(c.Fields, newField)
				relClass = c.Name
			}
		}
		if relClass != "" {
			class.addOneToManyRelation(allClasses, tableName, key, relClass)
		}
	}
}

func (class *Class) addOneToManyRelation(allClasses []*Class, tableName, key, relClass string) {
	for _, c := range allClasses {
		if c.Name == tableName {
			for i, field := range c.Fields {
				if field.Name == key {
					c.Fields = append(c.Fields[:i], c.Fields[i+1:]...)
				}
			}
			field := Field{
				Name:        strings.ToLower(relClass),
				DataType:    utils.List(relClass),
				Modifier:    "private",
				Annotations: []string{fmt.Sprintf(enums.OneToMany, strings.ToLower(tableName))},
				Position:    relClass,
			}
			c.Fields = append(c.Fields, field)
		}
	}
}

func (class *Class) addOneToOne(allClasses []*Class, tableName, value string) {
	if !strings.Contains(value, "unique") {
		return
	}

	re := regexp.MustCompile(`foreign_key\{([^\}]+)\}`)
	match := re.FindStringSubmatch(value)
	if len(match) < 2 {
		return
	}
	relatedClassNames := regexp.MustCompile(`\s*,\s*`).Split(match[1], -1)

	for _, c := range allClasses {
		if strings.EqualFold(tableName, c.Name) {
			newField := Field{
				Name:        strings.ToLower(relatedClassNames[0]),
				DataType:    strings.Title(relatedClassNames[0]),
				Modifier:    "private",
				Annotations: []string{enums.OneToOne, enums.MapsId},
				Position:    c.Name,
			}
			fmt.Println(newField)
			c.Fields = append(c.Fields, newField)
		}
	}
}

func (class *Class) GenerateEntity() string {
	var builder strings.Builder
	class.formationEntity()
	class.EntityLombokAnnotations()
	class.Position = "Entity"

	fmt.Fprintf(&builder, "%s class %s implements Serializable {\n", class.Modifier, class.Name)
	for _, field := range class.Fields {
		fmt.Fprintf(&builder, "%s", field.GenerateStringField())
	}
	for _, method := range class.Methods {
		mtd, err := method.GenerateStringMethod()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(&builder, "%s", mtd)
	}
	fmt.Fprintf(&builder, "}")

	return utils.AddAnnotations(class.Annotations, builder.String(), enums.Class)
}
