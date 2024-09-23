package generator

import (
	"encoding/json"
	"fois-generator/models"
	jsonmodels "fois-generator/models/json_models"
	"io"
	"os"
)

var (
	classes []*models.Class
)

func GetClasses() []*models.Class {
	return classes
}

func GenerateFileds(class *models.Class, table jsonmodels.Table) {
	class.CreateFields(table)
}

func GenerateClasses(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var project jsonmodels.Project
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &project)

	for _, table := range project.Tables {
		class := models.Class{
			Name:     table.Name,
			Modifier: table.Modifier,
		}
		GenerateFileds(&class, table)
		classes = append(classes, &class)
	}

	for _, table := range project.Tables {
		var class *models.Class
		for _, cl := range classes {
			if cl.Name == table.Name {
				class = cl
				break
			}
		}
		class.AddRelated(classes, table)
	}
	return nil
}
