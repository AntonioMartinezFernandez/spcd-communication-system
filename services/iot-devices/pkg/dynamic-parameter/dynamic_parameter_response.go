package dynamic_parameter

type DynamicParameterResponse struct {
	ID           string      `jsonapi:"primary,dynamic_parameter"`
	Name         string      `jsonapi:"attr,name"`
	DefaultValue interface{} `jsonapi:"attr,default_value"`
	DynamicValue interface{} `jsonapi:"attr,dynamic_value"`
}

func NewDynamicParameterFromParameter(id string, parameter *DynamicParameter) *DynamicParameterResponse {
	return &DynamicParameterResponse{
		ID:           id,
		Name:         parameter.Name.Value(),
		DefaultValue: parameter.DefaultValue,
		DynamicValue: parameter.DynamicValue,
	}
}
