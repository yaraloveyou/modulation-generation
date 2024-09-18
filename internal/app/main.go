package app

import (
	"fmt"
	"fois-generator/internal/generated"
)

func Start() {
	err := generated.GeneratedFileSystem()
	if err != nil {
		fmt.Println(err)
	}
	err = generated.GeneratedClass()
	if err != nil {
		fmt.Println(err)
	}
}
