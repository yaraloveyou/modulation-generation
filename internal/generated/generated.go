package generated

import (
	"encoding/json"
	"fmt"
	"fois-generator/config"
	"fois-generator/models"
	jsonmodels "fois-generator/models/json_models"
	"io"
	"os"
	"path"
	"strings"
)

var (
	classes []*models.Class
)

func GeneratedClass() error {
	return generateClasses("../../example/tables.json")
}

func Classes() []*models.Class {
	return classes
}

func generateClasses(path string) error {
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
		class.SetInfoAnyClass(classes)
		class.AddRelaited(table)
	}
	for _, class := range classes {
		GeneratedFile(class)
	}
	return nil
}

func GeneratedFileSystem() error {
	var folder jsonmodels.FileSystem
	jsonData, err := os.ReadFile(config.GetConfig().FolderStruct)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &folder)
	if err != nil {
		return err
	}

	folders := strings.ReplaceAll(folder.Pkg, ".", "/")
	dirPath := path.Join(config.GetConfig().OutputDir, folders)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			return err
		}
	}

	for _, value := range folder.Folders {
		layerFolder := path.Join(dirPath, value)
		if _, err := os.Stat(layerFolder); os.IsNotExist(err) {
			err = os.Mkdir(layerFolder, 0777)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GenerateFileds(class *models.Class, table jsonmodels.Table) {
	class.CreateFields(table)
}

func GeneratedFile(class *models.Class) error {
	// Передача информации о всех классах
	class.SetInfoAnyClass(classes)
	// Генерация кода Entity
	text := class.GenerateEntity()
	dirPath := findFullPath(class.Position)
	fullPath := path.Join(dirPath, fmt.Sprintf("%s%s.java", class.Name, class.Position))

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// формирование package
	pkg := fmt.Sprintf("package %s;\n\n", strings.Trim(strings.ReplaceAll(dirPath, "/", "."), "."))

	data := []byte(pkg + text)
	file.Write(data)

	fmt.Printf("File %s has been created\n", fullPath)
	return nil
}

func findFullPath(layer string) string {
	var folder jsonmodels.FileSystem
	jsonData, err := os.ReadFile(config.GetConfig().FolderStruct)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jsonData, &folder)
	if err != nil {
		fmt.Println(err)
	}

	folders := strings.ReplaceAll(folder.Pkg, ".", "/")
	dirPath := path.Join(config.GetConfig().OutputDir, folders)

	for key, value := range folder.Folders {
		if strings.EqualFold(layer, key) {
			layer = strings.ReplaceAll(value, ".", "/")
			break
		}
	}

	return path.Join(dirPath, strings.ToLower(layer))
}
