package dynamic_parameter

type ParameterName string

func (pn ParameterName) Value() string {
	return string(pn)
}

type ParameterValue interface{}

type DynamicParameter struct {
	Name         ParameterName
	DefaultValue ParameterValue
	DynamicValue ParameterValue
}

func (dp *DynamicParameter) DefaultIfDynamicIsNil() ParameterValue {
	if dp.DynamicValue != nil {
		return dp.DynamicValue
	}

	return dp.DefaultValue
}

func (dp *DynamicParameter) ToPlain() map[string]interface{} {
	return map[string]interface{}{
		"name":          dp.Name.Value(),
		"default":       dp.DefaultValue,
		"dynamic_value": dp.DynamicValue,
	}
}

func FromPlain(values map[string]interface{}) *DynamicParameter {
	return &DynamicParameter{
		Name:         ParameterName(values["name"].(string)),
		DefaultValue: values["default"],
		DynamicValue: values["value"],
	}
}

func (dp *DynamicParameter) WithNewValue(value interface{}) DynamicParameter {
	return DynamicParameter{
		Name:         dp.Name,
		DefaultValue: dp.DefaultValue,
		DynamicValue: value,
	}
}
