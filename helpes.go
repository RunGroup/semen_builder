package semen_builder

import "fmt"

func ToArgs[T any](data []T) []any {
	result := []any{}
	for _, value := range data {
		fmt.Println(value)
		result = append(result, value)
	}

	return result
}
