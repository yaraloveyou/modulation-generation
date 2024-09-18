package models

import (
	"fois-generator/config"
	"fois-generator/internal/enums"
	"strings"
)

func (class *Class) EntityLombokAnnotations() {
	if !config.GetConfig().IsLombok {
		return
	}

	// class.Getter()
	// class.Setter()
	// class.ToString()
	// class.NoArgsConstructor()
	// class.AllArgsConstructor()
}

func (c *Class) Getter() {
	c.addAnnotationAndFilterMethods(enums.Getter, func(m Method) bool {
		return !strings.Contains(m.Name, "get")
	})
}

func (c *Class) Setter() {
	c.addAnnotationAndFilterMethods(enums.Setter, func(m Method) bool {
		return !strings.Contains(m.Name, "set")
	})
}

func (c *Class) ToString() {
	c.addAnnotationAndFilterMethods(enums.ToString, func(m Method) bool {
		return m.Name != "toString"
	})
}

func (c *Class) NoArgsConstructor() {
	c.addAnnotationAndFilterMethods(enums.NoArgsConstructor, func(m Method) bool {
		return !(len(m.Variables) == 0 && len(m.ExternalVariables) == 0 && c.Name == m.Name)
	})
}

func (c *Class) AllArgsConstructor() {
	c.addAnnotationAndFilterMethods(enums.AllArgsConstructor, func(m Method) bool {
		return !(len(m.Variables) > 0 && len(m.Variables) == len(m.ExternalVariables) && c.Name == m.Name)
	})
}

func (c *Class) addAnnotationAndFilterMethods(annotation string, filterFunc func(Method) bool) {
	c.Annotations = append(c.Annotations, annotation)
	var updatedMethods []Method
	for _, m := range c.Methods {
		if filterFunc(m) {
			updatedMethods = append(updatedMethods, m)
		}
	}
	c.Methods = updatedMethods
}
