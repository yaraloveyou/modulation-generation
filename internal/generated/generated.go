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

func GeneratedClass() error {
	return generateClasses("../../example/tables.json")
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
		err = GeneratedFile(&class, table)
		if err != nil {
			return err
		}
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

func GeneratedFile(class *models.Class, table jsonmodels.Table) error {
	text := class.GenerateEntity(table)
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
