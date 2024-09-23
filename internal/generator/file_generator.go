package generator

import (
	"encoding/json"
	"fmt"
	"fois-generator/config"
	"fois-generator/models"
	jsonmodels "fois-generator/models/json_models"
	"os"
	"path"
	"strings"
)

func GenerateFiles() error {
	for _, class := range GetClasses() {
		if err := GenerateFile(class); err != nil {
			return err
		}
	}
	return nil
}

func GenerateFile(class *models.Class) error {
	return WriteFile(class)
}

func WriteFile(class *models.Class) error {
	text := class.GenerateEntity()
	dirPath := findFullPath(class.Position)
	fullPath := path.Join(dirPath, fmt.Sprintf("%s%s.java", class.Name, class.Position))

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	pkg := fmt.Sprintf("package %s;\n\n", strings.Trim(strings.ReplaceAll(dirPath, "/", "."), "."))

	data := []byte(pkg + text)
	if _, err := file.Write(data); err != nil {
		return err
	}

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
