package dynamic_parameter

import "context"

type DynamicParameterRepository interface {
	Save(ctx context.Context, parameter DynamicParameter) error
	Search(ctx context.Context, name ParameterName) (interface{}, error)
}
