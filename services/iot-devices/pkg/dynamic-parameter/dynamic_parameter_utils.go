package dynamic_parameter

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

func InitDefaultDynamicParametersConfigFromYamlFile(filename string) map[string]interface{} {
	parameters := make(map[string]interface{})

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if yamlErr := yaml.Unmarshal(yamlFile, &parameters); yamlErr != nil {
		panic(yamlErr)
	}

	return parameters
}

func BoolDynamicParameterValue(parameter *DynamicParameter) (bool, error) {
	value := parameter.DefaultIfDynamicIsNil()
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return strconv.ParseBool(fmt.Sprintf("%v", value))
	case reflect.Bool:
		return value.(bool), nil
	default:
		return false, NewInvalidDynamicParameterType(parameter.Name, "bool")
	}
}

func IntDynamicParameterValue(parameter *DynamicParameter) (int, error) {
	value := parameter.DefaultIfDynamicIsNil()
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return strconv.Atoi(fmt.Sprintf("%v", value))
	case reflect.Float64:
		return int(value.(float64)), nil
	case reflect.Int:
		return value.(int), nil
	default:
		return 0, NewInvalidDynamicParameterType(parameter.Name, "int")
	}
}
