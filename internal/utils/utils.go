package utils

import (
	"fmt"
	"strings"
)

func AddAnnotations(annotations []string, str string, tabs int) string {
	var builder strings.Builder
	var tabulation string
	for i := 0; i < tabs; i++ {
		tabulation += "\t"
	}
	for _, anotation := range annotations {
		fmt.Fprintf(&builder, "%s%s\n", tabulation, anotation)
	}
	fmt.Fprintf(&builder, "%s%s", tabulation, str)
	return builder.String()
}

func Contains(slice []string, val string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, val) {
			return true
		}
	}

	return false
}
