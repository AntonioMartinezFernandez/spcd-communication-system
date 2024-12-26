package dynamic_parameter

import (
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/domain"
)

const parameterNotExistsErrorMessage = "Dynamic parameter not exists"

type DynamicParameterNotExists struct {
	items map[string]interface{}
	domain.RootDomainError
}

func (pne DynamicParameterNotExists) Error() string {
	return parameterNotExistsErrorMessage
}

func (pne DynamicParameterNotExists) ExtraItems() map[string]interface{} {
	return pne.items
}

func NewDynamicParameterNotExists(name ParameterName) *DynamicParameterNotExists {
	return &DynamicParameterNotExists{items: map[string]interface{}{"name": name.Value()}}
}
