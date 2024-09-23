package app

import (
	"fmt"
	"fois-generator/internal/generator"
)

func Start() {
	err := generator.GeneratedFileSystem()
	if err != nil {
		fmt.Println(err)
	}
	err = generator.BuildProject()
	if err != nil {
		fmt.Println(err)
	}
}
