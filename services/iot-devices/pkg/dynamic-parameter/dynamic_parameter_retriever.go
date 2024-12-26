package dynamic_parameter

import (
	"context"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"
)

type DynamicParameterRetriever struct {
	repository       DynamicParameterRepository
	parametersConfig map[string]interface{}
}

func NewDynamicParameterRetriever(
	repository DynamicParameterRepository,
	parametersConfig map[string]interface{},
) *DynamicParameterRetriever {
	return &DynamicParameterRetriever{
		repository:       repository,
		parametersConfig: parametersConfig,
	}
}

func NewDynamicParameterRetrieverFromConfigFile(
	repository DynamicParameterRepository,
	configFilePath string,
) *DynamicParameterRetriever {
	parameters := InitDefaultDynamicParametersConfigFromYamlFile(configFilePath)
	return NewDynamicParameterRetriever(repository, parameters)
}

func (dr *DynamicParameterRetriever) Get(ctx context.Context, name ParameterName) (*DynamicParameter, error) {
	defaultValue, ok := dr.parametersConfig[name.Value()]
	if !ok {
		return nil, NewDynamicParameterNotExists(name)
	}

	if value, ok := defaultValue.(map[interface{}]interface{}); ok {
		defaultValue = utils.MapInterfaceInterfaceToStringInterface(value)
	}

	parameterValue, err := dr.repository.Search(ctx, name)
	if err != nil {
		return nil, err
	}

	return &DynamicParameter{
		Name:         name,
		DefaultValue: defaultValue,
		DynamicValue: parameterValue,
	}, nil
}
