package models

import (
	"fois-generator/internal/enums"
	"strings"
)

func (c *Class) Getter() {
	c.Layers = append(c.Layers, enums.Getter)
	var updatedMethods []Method
	for _, m := range c.Methods {
		if strings.Contains(m.Name, "get") {
			continue
		}
		updatedMethods = append(updatedMethods, m)
	}

	c.Methods = updatedMethods
}

func (c *Class) Setter() {
	c.Layers = append(c.Layers, enums.Setter)
	var updatedMethods []Method
	for _, m := range c.Methods {
		if strings.Contains(m.Name, "set") {
			continue
		}
		updatedMethods = append(updatedMethods, m)
	}

	c.Methods = updatedMethods
}

func (c *Class) Builder() {
	c.Layers = append(c.Layers, enums.Builder)
}
