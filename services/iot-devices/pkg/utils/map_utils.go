package utils

import (
	"fmt"
	"maps"
)

func MapInterfaceInterfaceToStringInterface(input map[interface{}]interface{}) map[string]interface{} {
	casted := make(map[string]interface{})
	for s, i := range input {
		casted[fmt.Sprint(s)] = i
	}

	return casted
}

func GetInMapValueOrDefault(keyToFind []string, myMap map[string]interface{}, defaultValue interface{}) interface{} {
	current := myMap
	for i := 0; i < len(keyToFind); i++ {
		if _, found := current[keyToFind[i]]; !found {
			return defaultValue
		}

		if i+1 == len(keyToFind) {
			return current[keyToFind[i]]
		}

		current = current[keyToFind[i]].(map[string]interface{})
	}

	return current
}

func MapStringStructToSlice(input map[string]struct{}) []string {
	keys := make([]string, 0, len(input))
	for key := range maps.Keys(input) {
		keys = append(keys, key)
	}

	return keys
}
