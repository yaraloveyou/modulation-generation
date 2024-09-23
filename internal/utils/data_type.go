package utils

import "fmt"

func List(dataType string) string {
	return fmt.Sprintf("List<%s>", dataType)
}
