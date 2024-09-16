package generated

import (
	"encoding/json"
	"fois-generator/models"
	jsonmodels "fois-generator/models/json_models"
	"io"
	"os"
)

func GeneratedClass() {
	jsonToFields("../../example/tables.json")
}

func jsonToFields(path string) ([]models.Class, error) {
	classes := []models.Class{}

	jsonFile, err := os.Open(path)
	if err != nil {
		return []models.Class{}, err
	}

	defer jsonFile.Close()

	var tables []jsonmodels.Table
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &tables)
	for _, table := range tables {
		class := models.Class{
			Name:     table.Name,
			Modifier: "protected",
		}
		class.GeneratedEntity(table)
		classes = append(classes, class)
	}
	return classes, nil
}
