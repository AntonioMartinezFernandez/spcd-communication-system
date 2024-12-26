package dynamic_parameter

import (
	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/domain"
)

const parameterWithInvalidTypeErrorMessage = "Invalid dynamic parameter type"

type InvalidDynamicParameterType struct {
	items map[string]interface{}
	domain.RootDomainError
}

func (ipt InvalidDynamicParameterType) Error() string {
	return parameterWithInvalidTypeErrorMessage
}

func (ipt InvalidDynamicParameterType) ExtraItems() map[string]interface{} {
	return ipt.items
}

func NewInvalidDynamicParameterType(name ParameterName, toType string) *InvalidDynamicParameterType {
	return &InvalidDynamicParameterType{items: map[string]interface{}{
		"name":          name.Value(),
		"expected_type": toType,
	}}
}
